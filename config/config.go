package config

import "github.com/pitabwire/frame"

type ProfileConfig struct {
	frame.ConfigurationDefault
	NotificationServiceURI string `default:"127.0.0.1:7020" envconfig:"NOTIFICATION_SERVICE_URI"`
	PartitionServiceURI    string `default:"127.0.0.1:7003" envconfig:"PARTITION_SERVICE_URI"`

	ContactEncryptionKey  string `required:"true" envconfig:"CONTACT_ENCRYPTION_KEY"`
	ContactEncryptionSalt string `required:"true" envconfig:"CONTACT_ENCRYPTION_SALT"`

	SystemAccessID        string `default:"c8cf0ldstmdlinc3eva0" envconfig:"STATIC_SYSTEM_ACCESS_ID"`
	QueueVerification     string `default:"mem://contact_verification_queue" envconfig:"QUEUE_VERIFICATION"`
	QueueVerificationName string `default:"contact_verification_queue" envconfig:"QUEUE_VERIFICATION_NAME"`

	LengthOfVerificationPin        int `default:"5" envconfig:"LENGTH_OF_VERIFICATION_PIN"`
	LengthOfVerificationLinkHash   int `default:"70" envconfig:"LENGTH_OF_VERIFICATION_LINK_HASH"`
	VerificationPinExpiryTimeInSec int `default:"86400" envconfig:"VERIFICATION_PIN_EXPIRY_TIME_IN_SEC"`

	MessageTemplateContactVerification string `default:"template.papi.contact.verification" envconfig:"MESSAGE_TEMPLATE_CONTACT_VERIFICATION"`
}
