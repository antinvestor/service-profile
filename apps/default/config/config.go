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

	// DeploymentEnvironment refuses the static identity key when it is
	// "prod" or "production" (SERVICE_ENVIRONMENT is checked as well).
	DeploymentEnvironment string `envDefault:"" env:"ENVIRONMENT"`

	// Identity salt (HMAC-SHA256 over "stawi/identity/v1" ‖ profile_id).
	// Vault Transit when VAULT_ADDR is set; otherwise the static key; with
	// neither, accounts are not derived.
	VaultAddress          string `envDefault:""               env:"VAULT_ADDR"`
	VaultK8sAuthRole      string `envDefault:""               env:"VAULT_K8S_AUTH_ROLE"`
	VaultK8sAuthMount     string `envDefault:"kubernetes"     env:"VAULT_K8S_AUTH_MOUNT"`
	VaultK8sTokenPath     string `envDefault:""               env:"VAULT_K8S_TOKEN_PATH"`
	VaultTransitMount     string `envDefault:"transit"        env:"VAULT_TRANSIT_MOUNT"`
	IdentityTransitKey    string `envDefault:"stawi-identity" env:"STAWI_IDENTITY_TRANSIT_KEY"`
	IdentityTransitKeyVer int    `envDefault:"0"              env:"STAWI_IDENTITY_KEY_VERSION"`
	IdentityStaticKeyHex  string `envDefault:""               env:"STAWI_IDENTITY_STATIC_KEY"`

	// Account derivation inputs from the protocol manifest. The init code is
	// keccak256(creation_code ‖ identity_salt_hash), so it differs per
	// profile and the creation code itself is required.
	AccountFactory          string `envDefault:"" env:"STAWI_ACCOUNT_FACTORY"`
	AccountCreationCode     string `envDefault:"" env:"STAWI_ACCOUNT_CREATION_CODE"`
	AccountCreationCodeFile string `envDefault:"" env:"STAWI_ACCOUNT_CREATION_CODE_FILE"`
	AccountCreationCodeHash string `envDefault:"" env:"STAWI_ACCOUNT_CREATION_CODE_HASH"`
	AccountVersion          uint32 `envDefault:"1" env:"STAWI_ACCOUNT_VERSION"`

	// AccountBackfillBatchSize is how many profiles one backfill pass derives.
	AccountBackfillBatchSize int `envDefault:"200" env:"ACCOUNT_BACKFILL_BATCH_SIZE"`
}
