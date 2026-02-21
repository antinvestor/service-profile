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
	// RateLimitTURNPerMinute is the max TURN credential requests per caller per minute.
	RateLimitTURNPerMinute int64 `envDefault:"10" env:"RATE_LIMIT_TURN_PER_MINUTE"`

	// TURNProvider selects the TURN credential provider: "cloudflare" or "static".
	// "cloudflare" uses Cloudflare's TURN API to generate credentials.
	// "static" generates time-limited credentials locally using a shared secret (for coturn/pion).
	TURNProvider string `envDefault:"static" env:"TURN_PROVIDER"`
	// TURNTTL is the time-to-live in seconds for generated TURN credentials. Defaults to 3600 (1 hour).
	// Valid range: 60–86400 seconds.
	TURNTTL int32 `envDefault:"3600" env:"TURN_TTL"`
	// CloudflareTURNTokenID is the TURN Token ID from Cloudflare's Real-Time Communications settings.
	CloudflareTURNTokenID string `env:"CLOUDFLARE_TURN_TOKEN_ID"`
	// CloudflareTURNAPIToken is the API token used to authenticate with Cloudflare's TURN credential generation endpoint.
	CloudflareTURNAPIToken string `env:"CLOUDFLARE_TURN_API_TOKEN"`

	// TURNServerURLs is a comma-separated list of TURN/STUN server URLs for the static provider.
	// Each URL must start with "turn:", "turns:", or "stun:".
	// Example: "turn:turn.example.com:3478,turns:turn.example.com:5349,stun:stun.example.com:3478"
	TURNServerURLs string `env:"TURN_SERVER_URLS"`
	// TURNSharedSecret is the shared secret used to generate time-limited credentials (coturn/pion static-auth-secret).
	TURNSharedSecret string `env:"TURN_SHARED_SECRET"`
}
