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

package business_test

import (
	"testing"
	"time"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	frevents "github.com/pitabwire/frame/v2/events"
	"github.com/pitabwire/frame/v2/frametests/definition"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/default/service/business"
	"github.com/antinvestor/service-profile/apps/default/tests"
)

// TestProfileCreatedFactIsPublishedExactlyOnce is the K5 guarantee for the
// profile aggregate: the fact is staged in the transaction that writes the
// profile, the relay publishes it once, and a further pass publishes nothing.
func (pts *ProfileTestSuite) TestProfileCreatedFactIsPublishedExactlyOnce() {
	pts.WithTestDependancies(pts.T(), func(t *testing.T, dep *definition.DependencyOption) {
		ctx, svc := pts.CreateService(t, dep)
		pb, _ := pts.getProfileBusiness(ctx, svc)

		contact := util.RandomAlphaNumericString(12) + "@testing.com"
		created, err := pb.CreateProfile(ctx, &profilev1.CreateRequest{
			Type:    profilev1.ProfileType_PERSON,
			Contact: contact,
		})
		require.NoError(t, err)
		require.NotEmpty(t, created.GetId())

		staged, err := tests.StagedFacts(ctx, svc, business.FactProfileCreated)
		require.NoError(t, err)
		require.Len(t, staged, 1, "creating a profile must stage exactly one fact")
		require.Nil(t, staged[0].PublishedAt)
		assert.Equal(t, created.GetId(), staged[0].AggregateID)

		published, err := pts.FactRelay.Drain(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, published)

		require.Eventually(t, func() bool {
			return len(pts.Facts.Named(business.FactProfileCreated)) == 1
		}, 10*time.Second, 25*time.Millisecond, "the fact must reach the topic once")

		fact := pts.Facts.Named(business.FactProfileCreated)[0]
		assert.Equal(t, created.GetId(), fact.Payload["profile_id"])
		assert.NotEmpty(t, fact.Payload["created_at"])
		assert.NotEmpty(t, fact.EventID)
		// The frame event header is what frame's own events.Manager routes on,
		// so a consumer registers this fact as a frame event rather than
		// writing a queue worker of its own.
		assert.Equal(t, business.FactProfileCreated, fact.Metadata[frevents.EventHeaderName])
		assert.Equal(t, business.FactProfileCreated, fact.Metadata["event_name"])
		assert.Equal(t, fact.EventID, fact.Metadata["event_id"])
		assert.NotContains(t, fact.Payload, "contact", "a fact carries no personal data")

		republished, err := pts.FactRelay.Drain(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, republished, "a second relay pass must not republish")
		assert.Len(t, pts.Facts.Named(business.FactProfileCreated), 1)
	})
}
