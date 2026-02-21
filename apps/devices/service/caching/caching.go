package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pitabwire/frame/cache"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/devices/config"
)

// TTL constants governing cache expiration across the devices service.
const (
	TTLDevice         = 5 * time.Minute
	TTLSession        = 5 * time.Minute
	TTLLatestSession  = 2 * time.Minute
	TTLPresence       = 1 * time.Hour
	TTLDeviceKeys     = 5 * time.Minute
	TTLGeoIP          = 24 * time.Hour
	TTLGeoIPNegative  = 1 * time.Hour
	TTLRateWindow     = 1 * time.Minute
	TTLSessionBatch   = 2 * time.Minute
	TTLLastSeenBuffer = 30 * time.Second

	// rateWindowSeconds is the sliding window size for rate limiting (1 minute).
	rateWindowSeconds = 60

	// rateCounterBytes is the byte size of a big-endian int64 counter.
	rateCounterBytes = 8

	// rateTTLMultiplier ensures rate-limit keys outlive the window boundary.
	rateTTLMultiplier = 2
)

// Key prefix constants for cache key namespacing.
const (
	prefixDevice        = "device:"
	prefixSession       = "session:"
	prefixLatestSession = "session:latest:"
	prefixPresence      = "presence:"
	prefixDeviceKeys    = "device:keys:"
	prefixGeoIP         = "geoip:"
	prefixGeoIPNeg      = "geoip:neg:"
	prefixRateLog       = "rate:log:"
	prefixRatePresence  = "rate:presence:"
	prefixLastSeen      = "lastseen:"
)

// DeviceCacheService provides typed cache operations for the devices service.
// It wraps frame's cache.Manager and provides domain-specific cache methods
// with consistent key formatting, TTLs, and serialization.
type DeviceCacheService struct {
	devices  cache.RawCache
	presence cache.RawCache
	geoip    cache.RawCache
	rate     cache.RawCache
}

// NewDeviceCacheService creates a DeviceCacheService from the frame cache manager.
// It retrieves the named caches configured during service initialization.
// Returns nil if cache manager is nil (graceful degradation for environments without cache).
func NewDeviceCacheService(cacheMan cache.Manager) *DeviceCacheService {
	if cacheMan == nil {
		return nil
	}

	devices, _ := cacheMan.GetRawCache(config.CacheNameDevices)
	presence, _ := cacheMan.GetRawCache(config.CacheNamePresence)
	geoip, _ := cacheMan.GetRawCache(config.CacheNameGeoIP)
	rate, _ := cacheMan.GetRawCache(config.CacheNameRate)

	// All caches must be available for the service to function correctly.
	if devices == nil || presence == nil || geoip == nil || rate == nil {
		return nil
	}

	return &DeviceCacheService{
		devices:  devices,
		presence: presence,
		geoip:    geoip,
		rate:     rate,
	}
}

// --- Device Cache Operations ---

// GetDevice retrieves a cached device by ID.
func (c *DeviceCacheService) GetDevice(ctx context.Context, id string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	val, found, err := c.devices.Get(ctx, prefixDevice+id)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache get device failed")
		return nil, false
	}
	return val, found
}

// SetDevice stores a serialized device in cache.
func (c *DeviceCacheService) SetDevice(ctx context.Context, id string, data []byte) {
	if c == nil {
		return
	}
	if err := c.devices.Set(ctx, prefixDevice+id, data, TTLDevice); err != nil {
		util.Log(ctx).WithError(err).Debug("cache set device failed")
	}
}

// InvalidateDevice removes a device from cache.
func (c *DeviceCacheService) InvalidateDevice(ctx context.Context, id string) {
	if c == nil {
		return
	}
	if err := c.devices.Delete(ctx, prefixDevice+id); err != nil {
		util.Log(ctx).WithError(err).Debug("cache invalidate device failed")
	}
}

// --- Session Cache Operations ---

// GetSession retrieves a cached session by session ID.
func (c *DeviceCacheService) GetSession(ctx context.Context, sessionID string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	val, found, err := c.devices.Get(ctx, prefixSession+sessionID)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache get session failed")
		return nil, false
	}
	return val, found
}

// SetSession stores a serialized session in cache.
func (c *DeviceCacheService) SetSession(ctx context.Context, sessionID string, data []byte) {
	if c == nil {
		return
	}
	if err := c.devices.Set(ctx, prefixSession+sessionID, data, TTLSession); err != nil {
		util.Log(ctx).WithError(err).Debug("cache set session failed")
	}
}

// GetLatestSession retrieves the latest session for a device.
func (c *DeviceCacheService) GetLatestSession(ctx context.Context, deviceID string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	val, found, err := c.devices.Get(ctx, prefixLatestSession+deviceID)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache get latest session failed")
		return nil, false
	}
	return val, found
}

// SetLatestSession stores the latest session data for a device.
func (c *DeviceCacheService) SetLatestSession(ctx context.Context, deviceID string, data []byte) {
	if c == nil {
		return
	}
	if err := c.devices.Set(ctx, prefixLatestSession+deviceID, data, TTLLatestSession); err != nil {
		util.Log(ctx).WithError(err).Debug("cache set latest session failed")
	}
}

// InvalidateLatestSession removes the latest session for a device from cache.
func (c *DeviceCacheService) InvalidateLatestSession(ctx context.Context, deviceID string) {
	if c == nil {
		return
	}
	if err := c.devices.Delete(ctx, prefixLatestSession+deviceID); err != nil {
		util.Log(ctx).WithError(err).Debug("cache invalidate latest session failed")
	}
}

// --- Presence Cache Operations ---

// PresenceEntry represents cached presence state.
type PresenceEntry struct {
	DeviceID      string `json:"device_id"`
	ProfileID     string `json:"profile_id"`
	Status        int32  `json:"status"`
	StatusMessage string `json:"status_message"`
	ExpiryTime    string `json:"expiry_time,omitempty"`
	LastActive    string `json:"last_active"`
}

// GetPresence retrieves the current presence state for a device.
func (c *DeviceCacheService) GetPresence(ctx context.Context, deviceID string) (*PresenceEntry, bool) {
	if c == nil {
		return nil, false
	}
	val, found, err := c.presence.Get(ctx, prefixPresence+deviceID)
	if err != nil || !found {
		return nil, false
	}
	var entry PresenceEntry
	if unmarshalErr := json.Unmarshal(val, &entry); unmarshalErr != nil {
		util.Log(ctx).WithError(unmarshalErr).Debug("cache unmarshal presence failed")
		return nil, false
	}
	return &entry, true
}

// SetPresence stores the current presence state for a device with TTL.
func (c *DeviceCacheService) SetPresence(
	ctx context.Context,
	deviceID string,
	entry *PresenceEntry,
	ttl time.Duration,
) {
	if c == nil || entry == nil {
		return
	}
	encoded, err := json.Marshal(entry)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache marshal presence failed")
		return
	}
	if setErr := c.presence.Set(ctx, prefixPresence+deviceID, encoded, ttl); setErr != nil {
		util.Log(ctx).WithError(setErr).Debug("cache set presence failed")
	}
}

// InvalidatePresence removes presence state from cache.
func (c *DeviceCacheService) InvalidatePresence(ctx context.Context, deviceID string) {
	if c == nil {
		return
	}
	if err := c.presence.Delete(ctx, prefixPresence+deviceID); err != nil {
		util.Log(ctx).WithError(err).Debug("cache invalidate presence failed")
	}
}

// --- Device Keys Cache Operations ---

// GetDeviceKeys retrieves cached keys for a device.
func (c *DeviceCacheService) GetDeviceKeys(ctx context.Context, deviceID string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	val, found, err := c.devices.Get(ctx, prefixDeviceKeys+deviceID)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache get device keys failed")
		return nil, false
	}
	return val, found
}

// SetDeviceKeys stores serialized device keys in cache.
func (c *DeviceCacheService) SetDeviceKeys(ctx context.Context, deviceID string, data []byte) {
	if c == nil {
		return
	}
	if err := c.devices.Set(ctx, prefixDeviceKeys+deviceID, data, TTLDeviceKeys); err != nil {
		util.Log(ctx).WithError(err).Debug("cache set device keys failed")
	}
}

// InvalidateDeviceKeys removes device keys from cache.
func (c *DeviceCacheService) InvalidateDeviceKeys(ctx context.Context, deviceID string) {
	if c == nil {
		return
	}
	if err := c.devices.Delete(ctx, prefixDeviceKeys+deviceID); err != nil {
		util.Log(ctx).WithError(err).Debug("cache invalidate device keys failed")
	}
}

// --- GeoIP Cache Operations ---

// GetGeoIP retrieves cached GeoIP data for an IP address.
func (c *DeviceCacheService) GetGeoIP(ctx context.Context, ip string) ([]byte, bool, bool) {
	if c == nil {
		return nil, false, false
	}

	// Check negative cache first (failed lookups).
	negExists, _ := c.geoip.Exists(ctx, prefixGeoIPNeg+ip)
	if negExists {
		return nil, true, true // found but negative
	}

	val, found, err := c.geoip.Get(ctx, prefixGeoIP+ip)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache get geoip failed")
		return nil, false, false
	}
	return val, found, false
}

// SetGeoIP stores GeoIP data in cache.
func (c *DeviceCacheService) SetGeoIP(ctx context.Context, ip string, data []byte) {
	if c == nil {
		return
	}
	if err := c.geoip.Set(ctx, prefixGeoIP+ip, data, TTLGeoIP); err != nil {
		util.Log(ctx).WithError(err).Debug("cache set geoip failed")
	}
}

// SetGeoIPNegative marks an IP as having a failed GeoIP lookup.
func (c *DeviceCacheService) SetGeoIPNegative(ctx context.Context, ip string) {
	if c == nil {
		return
	}
	if err := c.geoip.Set(ctx, prefixGeoIPNeg+ip, []byte("1"), TTLGeoIPNegative); err != nil {
		util.Log(ctx).WithError(err).Debug("cache set geoip negative failed")
	}
}

// --- Rate Limiting Operations ---

// CheckRateLimit checks whether the given key has exceeded the limit within the current window.
// Returns (allowed bool, currentCount int64).
//
// Uses a fixed-window approach keyed by (prefix, id, window_epoch).
// To ensure counter keys don't accumulate without TTL, we initialize the key
// via Set (with TTL) before the first Increment. The Increment implementation
// preserves the expiration set by Set, so subsequent increments keep the TTL.
//
// There is a narrow race where two goroutines both see Exists=false and both
// call Set, but the worst case is resetting a counter that just started —
// allowing at most one extra request through, which is acceptable for rate limiting.
func (c *DeviceCacheService) CheckRateLimit(ctx context.Context, prefix, deviceID string, limit int64) (bool, int64) {
	if c == nil {
		return true, 0
	}

	windowKey := fmt.Sprintf("%s%s:%d", prefix, deviceID, time.Now().Unix()/rateWindowSeconds)

	// Ensure the key exists with a TTL before incrementing.
	// Increment preserves the expiration, so the TTL set here persists.
	exists, _ := c.rate.Exists(ctx, windowKey)
	if !exists {
		// Initialize with zero-value counter and TTL of 2x window to survive boundary.
		initBytes := make([]byte, rateCounterBytes)
		_ = c.rate.Set(ctx, windowKey, initBytes, rateTTLMultiplier*TTLRateWindow)
	}

	count, err := c.rate.Increment(ctx, windowKey, 1)
	if err != nil {
		util.Log(ctx).WithError(err).Debug("cache rate limit increment failed")
		return true, 0 // allow on error
	}

	return count <= limit, count
}

// CheckLogRateLimit checks the device log rate limit.
func (c *DeviceCacheService) CheckLogRateLimit(ctx context.Context, deviceID string, limit int64) (bool, int64) {
	return c.CheckRateLimit(ctx, prefixRateLog, deviceID, limit)
}

// CheckPresenceRateLimit checks the presence update rate limit.
func (c *DeviceCacheService) CheckPresenceRateLimit(ctx context.Context, deviceID string, limit int64) (bool, int64) {
	return c.CheckRateLimit(ctx, prefixRatePresence, deviceID, limit)
}

// --- LastSeen Buffer Operations ---

// BufferLastSeen stores the latest LastSeen timestamp for a session in cache.
// This is used to coalesce frequent LastSeen DB writes.
func (c *DeviceCacheService) BufferLastSeen(ctx context.Context, sessionID string, lastSeen time.Time) {
	if c == nil {
		return
	}
	data := []byte(lastSeen.Format(time.RFC3339Nano))
	if err := c.devices.Set(ctx, prefixLastSeen+sessionID, data, TTLLastSeenBuffer); err != nil {
		util.Log(ctx).WithError(err).Debug("cache buffer lastseen failed")
	}
}

// GetBufferedLastSeen retrieves the buffered LastSeen for a session.
// Returns zero time and false if not found.
func (c *DeviceCacheService) GetBufferedLastSeen(ctx context.Context, sessionID string) (time.Time, bool) {
	if c == nil {
		return time.Time{}, false
	}
	val, found, err := c.devices.Get(ctx, prefixLastSeen+sessionID)
	if err != nil || !found {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339Nano, string(val))
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
