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

package outbox_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/internal/outbox"
)

// Without a configured egress nothing may drain the outbox: a relay into an
// in-process default queue would mark facts published and deliver them nowhere.
func TestEgressOptionsRetainFactsWithoutAQueueURI(t *testing.T) {
	t.Parallel()

	for _, uri := range []string{"", "   "} {
		opts := outbox.EgressOptions(t.Context(), nil, nil, "device.events", uri, 0, 0)
		require.Empty(t, opts, "uri %q must not register a publisher or relay", uri)
	}
}

func TestEgressOptionsRegisterPublisherAndRelayForAQueueURI(t *testing.T) {
	t.Parallel()

	opts := outbox.EgressOptions(t.Context(), nil, nil, "device.events", "mem://device.events", 0, 0)
	require.Len(t, opts, 2, "publisher and relay")
}
