package config

import (
	"time"

	"github.com/pitabwire/frame/v2/config"
)

type ProfileConfig struct {
	config.ConfigurationDefault

	//nolint:golines // Struct tags must remain a single valid reflect.StructTag literal.
	NotificationSvcURI                       string `envDefault:"127.0.0.1:7020" env:"NOTIFICATION_SERVICE_URI"`
	TenancyServiceURI                        string `envDefault:"127.0.0.1:7003" env:"TENANCY_SERVICE_URI"`
	NotificationServiceWorkloadAPITargetPath string `envDefault:"/ns/notifications/sa/service-notification" env:"NOTIFICATION_SERVICE_WORKLOAD_API_TARGET_PATH"`

	SystemAccessID string `envDefault:"c8cf0ldstmdlinc3eva0" env:"STATIC_SYSTEM_ACCESS_ID"`

	DEKLookupTokenHMACSHA256Key string `envDefault:"yZ9cW4nY7Jq6B7Xr0sN9dFv2mHkP8QY1EJ5VtLxD0uM=" env:"DEK_LOOKUP_TOKEN"`
	DEKActiveKeyID              string `envDefault:"contacts-dek-2026-01"                         env:"DEK_ACTIVE_KEY_ID"`
	DEKActiveAES256GCMKey       string `envDefault:"GZQ8s2m1Kc1yBZzV8YvWJ0l9M5RqK3a9QY7xYb9o7Ww=" env:"DEK_ACTIVE_ENCRYPTION_TOKEN"`
	DEKOldAES256GCMKey          string `envDefault:""                                             env:"DEK_OLD_ENCRYPTION_TOKEN"`

	QueueRelationshipConnectName string `envDefault:"relationships.connect"               env:"QUEUE_RELATIONSHIP_CONNECT_NAME"`
	QueueRelationshipConnectURI  string `envDefault:"mem://default.relationships.connect" env:"QUEUE_RELATIONSHIP_CONNECT_URI"`

	QueueRelationshipDisConnectName string `envDefault:"relationships.disconnect"               env:"QUEUE_RELATIONSHIP_DISCONNECT_NAME"`
	QueueRelationshipDisConnectURI  string `envDefault:"mem://default.relationships.disconnect" env:"QUEUE_RELATIONSHIP_DISCONNECT_URI"`

	// QueueProfileEvents is the single egress for the durable profile domain
	// facts staged in the outbox (profile.created). Consumers route on the
	// event_name header. Empty URI = no egress: facts stay in the outbox.
	QueueProfileEventsName string `envDefault:"profile.events" env:"QUEUE_PROFILE_EVENTS_NAME"`
	QueueProfileEventsURI  string `envDefault:""               env:"QUEUE_PROFILE_EVENTS_URI"`

	// OutboxRelayInterval is how often the relay looks for staged facts.
	OutboxRelayInterval time.Duration `envDefault:"2s" env:"OUTBOX_RELAY_INTERVAL"`
	// OutboxRelayBatchSize is how many staged facts one relay pass claims.
	OutboxRelayBatchSize int `envDefault:"100" env:"OUTBOX_RELAY_BATCH_SIZE"`

	LengthOfVerificationCode       int `envDefault:"6"     env:"LENGTH_OF_VERIFICATION_CODE"`
	VerificationPinExpiryTimeInSec int `envDefault:"86400" env:"VERIFICATION_PIN_EXPIRY_TIME_IN_SEC"`

	MessageTemplateContactVerification string `envDefault:"template.profilev1.contact.verification" env:"MESSAGE_TEMPLATE_CONTACT_VERIFICATION"`

	AuditServiceURI string `envDefault:"" env:"AUDIT_SERVICE_URI"`
}
