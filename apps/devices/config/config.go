package config

import "github.com/pitabwire/frame"

type DevicesConfig struct {
	frame.ConfigurationDefault

	PartitionServiceURI string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`

	QueueDeviceAnalysis     string `envDefault:"mem://device_analysis_queue" env:"QUEUE_DEVICE_ANALYSIS_URI"`
	QueueDeviceAnalysisName string `envDefault:"device_analysis_queue"       env:"QUEUE_DEVICE_ANALYSIS_NAME"`
}
