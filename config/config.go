package config

import "github.com/pitabwire/frame"

type ProfileConfig struct {
	frame.ConfigurationDefault

	NotificationServiceURI string `default:"127.0.0.1:7020" envconfig:"NOTIFICATION_SERVICE_URI"`
	PartitionServiceURI    string `default:"127.0.0.1:7003" envconfig:"PARTITION_SERVICE_URI"`

	SystemAccessID        string `default:"c8cf0ldstmdlinc3eva0" envconfig:"STATIC_SYSTEM_ACCESS_ID"`
	QueueVerification     string `default:"mem://contact_verification_queue" envconfig:"QUEUE_VERIFICATION_URI"`
	QueueVerificationName string `default:"contact_verification_queue" envconfig:"QUEUE_VERIFICATION_NAME"`

	QueueRelationshipConnectName string `default:"relationships.connect" envconfig:"QUEUE_RELATIONSHIP_CONNECT_NAME"`
	QueueRelationshipConnectURI  string `default:"mem://default.relationships.connect" envconfig:"QUEUE_RELATIONSHIP_CONNECT_URI"`

	QueueRelationshipDisConnectName string `default:"relationships.disconnect" envconfig:"QUEUE_RELATIONSHIP_DISCONNECT_NAME"`
	QueueRelationshipDisConnectURI  string `default:"mem://default.relationships.disconnect" envconfig:"QUEUE_RELATIONSHIP_DISCONNECT_URI"`

	QueueDeviceAnalysis     string `default:"mem://device_analysis_queue" envconfig:"QUEUE_DEVICE_ANALYSIS_URI"`
	QueueDeviceAnalysisName string `default:"device_analysis_queue" envconfig:"QUEUE_DEVICE_ANALYSIS_NAME"`

	LengthOfVerificationPin        int `default:"5" envconfig:"LENGTH_OF_VERIFICATION_PIN"`
	LengthOfVerificationLinkHash   int `default:"70" envconfig:"LENGTH_OF_VERIFICATION_LINK_HASH"`
	VerificationPinExpiryTimeInSec int `default:"86400" envconfig:"VERIFICATION_PIN_EXPIRY_TIME_IN_SEC"`

	MessageTemplateContactVerification string `default:"template.papi.contact.verification" envconfig:"MESSAGE_TEMPLATE_CONTACT_VERIFICATION"`
}
