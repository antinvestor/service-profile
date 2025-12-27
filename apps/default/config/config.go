package config

import (
	"github.com/pitabwire/frame/config"
)

type ProfileConfig struct {
	config.ConfigurationDefault

	NotificationServiceURI string `envDefault:"127.0.0.1:7020" env:"NOTIFICATION_SERVICE_URI"`
	PartitionServiceURI    string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`

	SystemAccessID string `envDefault:"c8cf0ldstmdlinc3eva0" env:"STATIC_SYSTEM_ACCESS_ID"`

	DEKLookupTokenHMACSHA256Key string `envDefault:"yZ9cW4nY7Jq6B7Xr0sN9dFv2mHkP8QY1EJ5VtLxD0uM=" env:"DEK_LOOKUP_TOKEN"`
	DEKActiveKeyID              string `envDefault:"contacts-dek-2026-01"                         env:"DEK_ACTIVE_KEY_ID"`
	DEKActiveAES256GCMKey       string `envDefault:"GZQ8s2m1Kc1yBZzV8YvWJ0l9M5RqK3a9QY7xYb9o7Ww=" env:"DEK_ACTIVE_ENCRYPTION_TOKEN"`
	DEKOldAES256GCMKey          string `envDefault:""                                             env:"DEK_OLD_ENCRYPTION_TOKEN"`

	QueueRelationshipConnectName string `envDefault:"relationships.connect"               env:"QUEUE_RELATIONSHIP_CONNECT_NAME"`
	QueueRelationshipConnectURI  string `envDefault:"mem://default.relationships.connect" env:"QUEUE_RELATIONSHIP_CONNECT_URI"`

	QueueRelationshipDisConnectName string `envDefault:"relationships.disconnect"               env:"QUEUE_RELATIONSHIP_DISCONNECT_NAME"`
	QueueRelationshipDisConnectURI  string `envDefault:"mem://default.relationships.disconnect" env:"QUEUE_RELATIONSHIP_DISCONNECT_URI"`

	LengthOfVerificationCode       int `envDefault:"6"     env:"LENGTH_OF_VERIFICATION_CODE"`
	VerificationPinExpiryTimeInSec int `envDefault:"86400" env:"VERIFICATION_PIN_EXPIRY_TIME_IN_SEC"`

	MessageTemplateContactVerification string `envDefault:"template.profilev1.contact.verification" env:"MESSAGE_TEMPLATE_CONTACT_VERIFICATION"`
}
