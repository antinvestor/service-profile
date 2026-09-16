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

// Package outbox stages domain facts in the same database transaction as the
// state change they describe, and relays them to the service's existing frame
// queue publisher.
//
// A fact is only true because a row changed; writing the fact through a
// separate call would let the two disagree whenever the process dies between
// them. The outbox removes that window: the fact row and the state change
// commit together or not at all, and a relay publishes whatever is committed
// and unpublished. The relay keeps a single egress — the frame queue publisher
// registered by the service — so there is exactly one path a fact can leave on.
//
// Delivery is at-least-once: a fact whose publish succeeds but whose
// bookkeeping fails is retried, so every fact carries an event_id and
// consumers deduplicate on it.
package outbox

import (
	"context"
	"time"

	"github.com/pitabwire/frame/v2/data"
	"gorm.io/gorm"
)

// Event is one staged domain fact.
type Event struct {
	data.BaseModel

	// Name is the fact name, e.g. "device.key.added".
	Name string `gorm:"index;size:120;not null" json:"name"`
	// AggregateType and AggregateID identify what the fact is about.
	AggregateType string `gorm:"size:60;not null" json:"aggregate_type"`
	AggregateID   string `gorm:"index;size:50"    json:"aggregate_id"`
	// Payload is the fact body as published.
	Payload data.JSONMap `gorm:"type:jsonb" json:"payload"`
	// OccurredAt is when the state change happened.
	OccurredAt time.Time `gorm:"not null" json:"occurred_at"`
	// PublishedAt is set once the relay has handed the fact to the queue.
	PublishedAt *time.Time `gorm:"index" json:"published_at,omitempty"`
	// Attempts counts relay passes that have touched this row.
	Attempts int `gorm:"not null;default:0" json:"attempts"`
}

// TableName pins the table name so both applications that embed this package
// keep the same shape.
func (Event) TableName() string {
	return "outbox_events"
}

// NewEvent stages a fact. The returned row is not written anywhere yet; pass
// it to Enqueue inside the transaction that performs the state change.
func NewEvent(ctx context.Context, name, aggregateType, aggregateID string, payload data.JSONMap) *Event {
	evt := &Event{
		Name:          name,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       payload,
		OccurredAt:    time.Now().UTC(),
	}
	evt.GenID(ctx)
	if evt.Payload == nil {
		evt.Payload = data.JSONMap{}
	}
	// The consumer deduplicates on event_id, so it must be in the body as
	// well as in the message header.
	evt.Payload["event_id"] = evt.GetID()
	evt.Payload["name"] = name
	evt.Payload["occurred_at"] = evt.OccurredAt.Format(time.RFC3339Nano)
	return evt
}

// Enqueue writes a staged fact on the supplied transaction handle. It must be
// the same handle the state change was written on — that is the whole point.
func Enqueue(tx *gorm.DB, evt *Event) error {
	return tx.Create(evt).Error
}
