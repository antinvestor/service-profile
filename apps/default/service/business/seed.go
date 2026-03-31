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
	AdminProfileID     = "d75qclkpf2t1uum8ij3g"

	SvcAuthenticationProfileID             = "d75qclkpf2t1uum8ij40"
	SvcProfileProfileID                    = "d75qclkpf2t1uum8ij4g"
	SvcTenancyProfileID                    = "d75qclkpf2t1uum8ij50"
	SvcNotificationProfileID               = "d75qclkpf2t1uum8ij5g"
	SvcDevicesProfileID                    = "d75qclkpf2t1uum8ij60"
	SvcSettingProfileID                    = "d75qclkpf2t1uum8ij6g"
	SvcPaymentProfileID                    = "d75qclkpf2t1uum8ij70"
	SvcPaymentJengaProfileID               = "d75qclkpf2t1uum8ij7g"
	SvcLedgerProfileID                     = "d75qclkpf2t1uum8ij80"
	SvcBillingProfileID                    = "d75qclkpf2t1uum8ij8g"
	SvcFileProfileID                       = "d75qclkpf2t1uum8ij90"
	SvcChatDroneProfileID                  = "d75qclkpf2t1uum8ij9g"
	SvcChatGatewayProfileID                = "d75qclkpf2t1uum8ija0"
	SvcFoundryProfileID                    = "d75qclkpf2t1uum8ijag"
	SvcGitvaultProfileID                   = "d75qclkpf2t1uum8ijb0"
	SvcTrustageProfileID                   = "d75qclkpf2t1uum8ijbg"
	SvcNotificationAfricastalkingProfileID = "d75qclkpf2t1uum8ijc0"
	SvcNotificationEmailSMTPProfileID      = "d75qclkpf2t1uum8ijcg"
	SvcLenderProfileID                     = "d75qclkpf2t1uum8ijd0"
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
	{SvcAuthenticationProfileID, "authentication.bot@stawi.org"},
	{SvcProfileProfileID, "profile.bot@stawi.org"},
	{SvcTenancyProfileID, "tenancy.bot@stawi.org"},
	{SvcNotificationProfileID, "notification.bot@stawi.org"},
	{SvcDevicesProfileID, "devices.bot@stawi.org"},
	{SvcSettingProfileID, "setting.bot@stawi.org"},
	{SvcPaymentProfileID, "payment.bot@stawi.org"},
	{SvcPaymentJengaProfileID, "payment-jenga.bot@stawi.org"},
	{SvcLedgerProfileID, "ledger.bot@stawi.org"},
	{SvcBillingProfileID, "billing.bot@stawi.org"},
	{SvcFileProfileID, "file.bot@stawi.org"},
	{SvcChatDroneProfileID, "chat-drone.bot@stawi.org"},
	{SvcChatGatewayProfileID, "chat-gateway.bot@stawi.org"},
	{SvcFoundryProfileID, "foundry.bot@stawi.org"},
	{SvcGitvaultProfileID, "gitvault.bot@stawi.org"},
	{SvcTrustageProfileID, "trustage.bot@stawi.org"},
	{SvcNotificationAfricastalkingProfileID, "notification-africastalking.bot@stawi.org"},
	{SvcNotificationEmailSMTPProfileID, "notification-emailsmtp.bot@stawi.org"},
	{SvcLenderProfileID, "lender.bot@stawi.org"},
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
