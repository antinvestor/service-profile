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

// Package outboxtest gives test suites a subscriber that records the domain
// facts an outbox relay published, so a test can assert on what actually left
// the service and how many times.
package outboxtest

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pitabwire/frame/v2/queue"
)

// Fact is one message seen on a facts topic.
type Fact struct {
	Name     string
	EventID  string
	Metadata map[string]string
	Payload  map[string]any
}

// Recorder is a queue subscriber that keeps every fact it receives.
type Recorder struct {
	mu    sync.Mutex
	facts []Fact
}

var _ queue.SubscribeWorker = (*Recorder)(nil)

// Handle records one published fact.
func (r *Recorder) Handle(_ context.Context, metadata map[string]string, message []byte) error {
	payload := map[string]any{}
	if err := json.Unmarshal(message, &payload); err != nil {
		return err
	}

	name, _ := payload["name"].(string)
	eventID, _ := payload["event_id"].(string)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.facts = append(r.facts, Fact{
		Name:     name,
		EventID:  eventID,
		Metadata: metadata,
		Payload:  payload,
	})
	return nil
}

// All returns every fact recorded so far.
func (r *Recorder) All() []Fact {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Fact(nil), r.facts...)
}

// Named returns the recorded facts with the given name.
func (r *Recorder) Named(name string) []Fact {
	var out []Fact
	for _, fact := range r.All() {
		if fact.Name == name {
			out = append(out, fact)
		}
	}
	return out
}

// Reset drops everything recorded so far.
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.facts = nil
}
