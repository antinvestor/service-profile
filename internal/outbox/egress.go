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
	"strings"
	"time"

	"github.com/pitabwire/frame/v2"
	"github.com/pitabwire/frame/v2/datastore/pool"
	"github.com/pitabwire/frame/v2/queue"
	"github.com/pitabwire/util"
)

// EgressOptions registers the publisher for the named fact queue and the
// relay that drains the outbox into it.
//
// With no queue URI there is no durable egress, so nothing is registered and
// facts stay in the outbox unpublished until one is configured. Relaying to a
// default in-process queue instead would mark every fact published while
// delivering it nowhere, losing it for good.
func EgressOptions(
	ctx context.Context,
	dbPool pool.Pool,
	qMan queue.Manager,
	name, uri string,
	interval time.Duration,
	batchSize int,
) []frame.Option {
	if strings.TrimSpace(uri) == "" {
		util.Log(ctx).WithField("queue", name).
			Warn("no egress configured for domain facts; they are retained in the outbox until one is")
		return nil
	}

	relay := NewRelay(dbPool, qMan, name).WithInterval(interval).WithBatchSize(batchSize)
	return []frame.Option{
		frame.WithRegisterPublisher(name, uri),
		frame.WithBackgroundConsumer(relay.Run),
	}
}
