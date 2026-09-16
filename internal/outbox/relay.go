// Copyright 2023-2026 Ant Investor Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pitabwire/frame/v2/datastore/pool"
	frevents "github.com/pitabwire/frame/v2/events"
	"github.com/pitabwire/frame/v2/queue"
	"github.com/pitabwire/frame/v2/tenancy"
	"github.com/pitabwire/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// DefaultBatchSize is how many staged facts one relay pass claims.
	DefaultBatchSize = 100
	// DefaultInterval is how often the background relay looks for work.
	DefaultInterval = 2 * time.Second
	// HeaderEventName carries the fact name on the queue message.
	HeaderEventName = "event_name"
	// HeaderFrameEventName is the header frame's own events.Manager routes on
	// (events.EventHeaderName). Carrying it lets a consumer handle these facts
	// with a registered frame event instead of a bespoke queue worker.
	HeaderFrameEventName = frevents.EventHeaderName
	// HeaderEventID carries the deduplication key on the queue message.
	HeaderEventID = "event_id"
	// HeaderAggregateID carries the subject of the fact.
	HeaderAggregateID = "aggregate_id"
)

// ErrFactNotClaimed is returned when marking a fact as published changes no
// row, which would otherwise let the next pass publish it again.
var ErrFactNotClaimed = errors.New("outbox fact could not be claimed")

// Relay publishes staged facts through the service's registered frame queue
// publisher. It is the only writer of PublishedAt.
type Relay struct {
	dbPool    pool.Pool
	qMan      queue.Manager
	topic     string
	batchSize int
	interval  time.Duration
}

// NewRelay builds a relay that publishes to the named frame publisher
// reference. The reference must have been registered with
// frame.WithRegisterPublisher.
func NewRelay(dbPool pool.Pool, qMan queue.Manager, topic string) *Relay {
	return &Relay{
		dbPool:    dbPool,
		qMan:      qMan,
		topic:     topic,
		batchSize: DefaultBatchSize,
		interval:  DefaultInterval,
	}
}

// WithInterval overrides how often Run looks for work.
func (r *Relay) WithInterval(interval time.Duration) *Relay {
	if interval > 0 {
		r.interval = interval
	}
	return r
}

// WithBatchSize overrides how many facts one pass claims.
func (r *Relay) WithBatchSize(size int) *Relay {
	if size > 0 {
		r.batchSize = size
	}
	return r
}

// Run drains the outbox until the context is cancelled. Wire it with
// frame.WithBackgroundConsumer.
func (r *Relay) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if _, err := r.Drain(ctx); err != nil {
				util.Log(ctx).WithError(err).Warn("outbox relay pass failed")
			}
		}
	}
}

// Drain publishes every fact that is committed and not yet published, and
// returns how many it published. Facts are claimed FOR UPDATE SKIP LOCKED so
// several replicas can relay the same outbox without publishing a fact twice,
// and PublishedAt is set in the claiming transaction so a later pass — or a
// redelivered message that causes a consumer to ask for a replay — never
// republishes a fact that has already gone out.
func (r *Relay) Drain(ctx context.Context) (int, error) {
	total := 0
	for {
		published, err := r.drainBatch(ctx)
		total += published
		if err != nil {
			return total, err
		}
		if published < r.batchSize {
			return total, nil
		}
	}
}

func (r *Relay) drainBatch(ctx context.Context) (int, error) {
	// The relay is process-local system work over every tenant's facts, so it
	// runs with tenancy enforcement explicitly bypassed rather than on
	// whatever ambient claims happen to be around.
	ctx = tenancy.WithSkipEnforcement(ctx)

	db := r.dbPool.DB(ctx, false)
	if db == nil {
		return 0, gorm.ErrInvalidDB
	}

	published := 0
	err := db.Transaction(func(tx *gorm.DB) error {
		var staged []*Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("published_at IS NULL").
			Order("created_at asc").
			Limit(r.batchSize).
			Find(&staged).Error; err != nil {
			return err
		}

		now := time.Now().UTC()
		for _, evt := range staged {
			headers := map[string]string{
				HeaderFrameEventName: evt.Name,
				HeaderEventName:      evt.Name,
				HeaderEventID:        evt.GetID(),
				HeaderAggregateID:    evt.AggregateID,
			}
			if err := r.qMan.Publish(ctx, r.topic, evt.Payload, headers); err != nil {
				return err
			}
			// A single statement claims the row: publishing once depends on
			// this update, so it is written as the exact SQL that runs.
			marked := tx.Exec(
				"UPDATE outbox_events SET published_at = ?, attempts = attempts + 1 WHERE id = ?",
				now, evt.GetID(),
			)
			if marked.Error != nil {
				return marked.Error
			}
			if marked.RowsAffected != 1 {
				// The claim is what makes publishing once possible; if it did
				// not land, the whole batch is rolled back and retried rather
				// than published again on the next pass.
				return fmt.Errorf("%w: %s", ErrFactNotClaimed, evt.GetID())
			}
			published++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return published, nil
}
