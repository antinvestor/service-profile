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

package business

import (
	"context"
	"time"

	"github.com/pitabwire/frame/v2/data"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// FactProfileCreated states that a profile now exists (GFOS §10.1, K5). It is
// staged in the transaction that writes the profile row.
const FactProfileCreated = "profile.created"

const aggregateProfile = "profile"

// ProfileCreatedFact builds the profile.created fact.
//
// Payload: profile_id, profile_type (the profile.v1.ProfileType name),
// profile_type_uid (its enum number), profile_type_id and created_at, plus the
// envelope fields event_id, name and occurred_at.
//
// Contact details and profile properties are deliberately absent: a consumer
// that needs them asks for them under its own authority, and the fact itself
// stays free of personal data.
func ProfileCreatedFact(ctx context.Context, profile *models.Profile) *outbox.Event {
	return outbox.NewEvent(ctx, FactProfileCreated, aggregateProfile, profile.GetID(), data.JSONMap{
		"profile_id":       profile.GetID(),
		"profile_type":     profile.ProfileType.Name,
		"profile_type_uid": profile.ProfileType.UID,
		"profile_type_id":  profile.ProfileTypeID,
		"created_at":       profile.CreatedAt.UTC().Format(time.RFC3339Nano),
	})
}
