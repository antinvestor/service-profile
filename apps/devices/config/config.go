package config

import "github.com/pitabwire/frame/config"

type DevicesConfig struct {
	config.ConfigurationDefault

	PartitionServiceURI string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`

	QueueDeviceAnalysis     string `envDefault:"mem://device_analysis_queue" env:"QUEUE_DEVICE_ANALYSIS_URI"`
	QueueDeviceAnalysisName string `envDefault:"device_analysis_queue"       env:"QUEUE_DEVICE_ANALYSIS_NAME"`

	FCMEndpoint          string `envDefault:"https://fcm.googleapis.com/fcm/send" env:"FCM_ENDPOINT"`
	FCMServerKey         string `env:"FCM_SERVER_KEY"`
	NotificationBatchSize int    `envDefault:"500" env:"NOTIFICATION_BATCH_SIZE"`
}
