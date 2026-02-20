package config

import "github.com/pitabwire/frame/config"

// Cache name constants for the devices service.
const (
	CacheNameDevices  = "devices"
	CacheNamePresence = "presence"
	CacheNameGeoIP    = "geoip"
	CacheNameRate     = "rate"
)

type DevicesConfig struct {
	config.ConfigurationDefault

	PartitionServiceURI string `envDefault:"127.0.0.1:7003" env:"PARTITION_SERVICE_URI"`

	QueueDeviceAnalysis     string `envDefault:"mem://device_analysis_queue" env:"QUEUE_DEVICE_ANALYSIS_URI"`
	QueueDeviceAnalysisName string `envDefault:"device_analysis_queue"       env:"QUEUE_DEVICE_ANALYSIS_NAME"`

	FCMMaxBatchSize int `envDefault:"500" env:"FCM_MAX_BATCH_SIZE"`

	// RateLimitLogPerMinute is the max device log events per device per minute.
	RateLimitLogPerMinute int64 `envDefault:"120" env:"RATE_LIMIT_LOG_PER_MINUTE"`
	// RateLimitPresencePerMinute is the max presence updates per device per minute.
	RateLimitPresencePerMinute int64 `envDefault:"30" env:"RATE_LIMIT_PRESENCE_PER_MINUTE"`
}
