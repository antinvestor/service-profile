package business

import (
	"context"

	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"
)

// Static profile IDs — these match the SQL migration and can be referenced
// by other services to look up or authenticate against known accounts.
const (
	SystemBotProfileID = "system_bot_profile_01"
	AdminProfileID     = "admin_profile_01"

	SvcAuthenticationProfileID = "svc_authentication_01"
	SvcNotificationProfileID   = "svc_notification_01"
	SvcTenancyProfileID        = "svc_tenancy_01"
)

// bootstrapProfile pairs a migration-created profile with the contact
// detail that must be added at the application level (requires encryption).
type bootstrapProfile struct {
	ProfileID string
	Contact   string
}

// bootstrapProfiles lists every profile created by migration that needs an
// encrypted contact record. Add new service accounts here and in the
// corresponding SQL migration.
var bootstrapProfiles = []bootstrapProfile{
	{SystemBotProfileID, "system.bot@stawi.org"},
	{AdminProfileID, "bwire517@gmail.com"},
	{SvcAuthenticationProfileID, "service-authentication@internal.antinvestor.com"},
	{SvcNotificationProfileID, "service-notification@internal.antinvestor.com"},
	{SvcTenancyProfileID, "service-tenancy@internal.antinvestor.com"},
}

// SeedBootstrapContacts ensures every bootstrap profile has its encrypted
// contact record. Profile rows are created by SQL migration; this function
// adds the contacts since they require application-level encryption.
// Safe to call multiple times — skips profiles that already have contacts.
func SeedBootstrapContacts(ctx context.Context, pb ProfileBusiness, cb ContactBusiness) error {
	log := util.Log(ctx)

	for _, bp := range bootstrapProfiles {
		if err := seedProfileContact(ctx, pb, cb, bp); err != nil {
			log.WithError(err).WithField("profile_id", bp.ProfileID).
				Warn("failed to seed bootstrap contact — will retry on next startup")
		}
	}

	return nil
}

func seedProfileContact(
	ctx context.Context,
	pb ProfileBusiness,
	cb ContactBusiness,
	bp bootstrapProfile,
) error {
	log := util.Log(ctx)

	profile, err := pb.GetByID(ctx, bp.ProfileID)
	if err != nil {
		log.WithError(err).WithField("profile_id", bp.ProfileID).
			Warn("bootstrap profile not found — run migration first")
		return nil
	}
	if profile == nil {
		return nil
	}

	// Skip if contact already linked.
	contacts, err := cb.GetByProfile(ctx, bp.ProfileID)
	if err == nil && len(contacts) > 0 {
		return nil
	}

	contact, err := cb.CreateContact(ctx, bp.Contact, data.JSONMap{"source": "seed"})
	if err != nil {
		return err
	}

	if _, updateErr := cb.UpdateContact(ctx, contact.GetID(), bp.ProfileID, nil); updateErr != nil {
		log.WithError(updateErr).WithField("profile_id", bp.ProfileID).
			Warn("failed to link bootstrap contact to profile")
	}

	log.WithField("profile_id", bp.ProfileID).
		WithField("contact", bp.Contact).
		Info("seeded bootstrap contact")
	return nil
}
