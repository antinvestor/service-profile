package config

import "github.com/pitabwire/frame"

type ProfileConfig struct {
	frame.ConfigurationDefault

	NotificationServiceURI string `envDefault:"127.0.0.1:7020" env:"NOTIFICATION_SERVICE_URI"`
	PartitionServiceURI    string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`

	SystemAccessID        string `envDefault:"c8cf0ldstmdlinc3eva0" env:"STATIC_SYSTEM_ACCESS_ID"`
	QueueVerification     string `envDefault:"mem://contact_verification_queue" env:"QUEUE_VERIFICATION_URI"`
	QueueVerificationName string `envDefault:"contact_verification_queue" env:"QUEUE_VERIFICATION_NAME"`

	QueueRelationshipConnectName string `envDefault:"relationships.connect" env:"QUEUE_RELATIONSHIP_CONNECT_NAME"`
	QueueRelationshipConnectURI  string `envDefault:"mem://default.relationships.connect" env:"QUEUE_RELATIONSHIP_CONNECT_URI"`

	QueueRelationshipDisConnectName string `envDefault:"relationships.disconnect" env:"QUEUE_RELATIONSHIP_DISCONNECT_NAME"`
	QueueRelationshipDisConnectURI  string `envDefault:"mem://default.relationships.disconnect" env:"QUEUE_RELATIONSHIP_DISCONNECT_URI"`

	QueueDeviceAnalysis     string `envDefault:"mem://device_analysis_queue" env:"QUEUE_DEVICE_ANALYSIS_URI"`
	QueueDeviceAnalysisName string `envDefault:"device_analysis_queue" env:"QUEUE_DEVICE_ANALYSIS_NAME"`

	LengthOfVerificationPin        int `envDefault:"5" env:"LENGTH_OF_VERIFICATION_PIN"`
	LengthOfVerificationLinkHash   int `envDefault:"70" env:"LENGTH_OF_VERIFICATION_LINK_HASH"`
	VerificationPinExpiryTimeInSec int `envDefault:"86400" env:"VERIFICATION_PIN_EXPIRY_TIME_IN_SEC"`

	MessageTemplateContactVerification string `envDefault:"template.papi.contact.verification" env:"MESSAGE_TEMPLATE_CONTACT_VERIFICATION"`
}
