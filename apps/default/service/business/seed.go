package business

import (
	"context"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"
)

// SystemBotProfileID is the fixed ID for the default system bot profile.
// All internal service accounts should reference this profile.
const SystemBotProfileID = "system_bot_profile_01"

// SystemBotContact is the contact detail for the system bot profile.
const SystemBotContact = "system.bot@stawi.org"

// SeedSystemBotContact ensures the system bot profile has its contact record.
// The profile row is created by SQL migration; this function adds the encrypted
// contact since contacts require application-level encryption.
// Safe to call multiple times — skips if the contact already exists.
func SeedSystemBotContact(ctx context.Context, pb ProfileBusiness, cb ContactBusiness) error {
	log := util.Log(ctx)

	// Check if profile exists
	profile, err := pb.GetByID(ctx, SystemBotProfileID)
	if err != nil {
		log.WithError(err).Warn("system bot profile lookup failed — run migration first")
		return nil
	}
	if profile == nil {
		log.Warn("system bot profile not found — run migration first")
		return nil
	}

	// Check if contact already exists for this profile
	contacts, err := cb.GetByProfile(ctx, SystemBotProfileID)
	if err == nil && len(contacts) > 0 {
		log.WithField("profile_id", SystemBotProfileID).Debug("system bot contact already exists")
		return nil
	}

	// Create the contact
	contact, err := cb.CreateContact(ctx, SystemBotContact, data.JSONMap{
		"source": "seed",
	})
	if err != nil {
		return err
	}

	// Link contact to the profile
	contact.ProfileID = SystemBotProfileID
	if _, updateErr := cb.UpdateContact(ctx, contact.GetID(), SystemBotProfileID, nil); updateErr != nil {
		log.WithError(updateErr).Warn("failed to link system bot contact to profile")
	}

	log.WithField("profile_id", SystemBotProfileID).
		WithField("contact_id", contact.GetID()).
		Info("seeded system bot contact")
	return nil
}
