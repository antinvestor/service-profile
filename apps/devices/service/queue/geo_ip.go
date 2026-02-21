package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pitabwire/frame/client"
	"github.com/pitabwire/util"

	"github.com/antinvestor/service-profile/apps/devices/service/caching"
)

const geoIPRequestTimeout = 5 * time.Second

type GeoIP struct {
	IP                 string  `json:"ip"`
	Network            string  `json:"network"`
	Version            string  `json:"version"`
	City               string  `json:"city"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	Country            string  `json:"country"`
	CountryName        string  `json:"country_name"`
	CountryCode        string  `json:"country_code"`
	CountryCodeIso3    string  `json:"country_code_iso3"`
	CountryCapital     string  `json:"country_capital"`
	CountryTld         string  `json:"country_tld"`
	ContinentCode      string  `json:"continent_code"`
	InEu               bool    `json:"in_eu"`
	Postal             string  `json:"postal"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
	CountryCallingCode string  `json:"country_calling_code"`
	Currency           string  `json:"currency"`
	CurrencyName       string  `json:"currency_name"`
	Languages          string  `json:"languages"`
	CountryArea        float64 `json:"country_area"`
	CountryPopulation  int64   `json:"country_population"`
	Asn                string  `json:"asn"`
	Org                string  `json:"org"`
}

// QueryIPGeo resolves an IP address to geographic location data.
// It checks the cache first and falls back to the external ipapi.co API.
// Failed lookups are negative-cached to avoid hammering the API for bad IPs.
func QueryIPGeo(
	ctx context.Context,
	cli client.Manager,
	ip string,
	cacheSvc *caching.DeviceCacheService,
) (*GeoIP, error) {
	// Try cache first.
	if cacheSvc != nil {
		cached, found, isNegative := cacheSvc.GetGeoIP(ctx, ip)
		if isNegative {
			return nil, fmt.Errorf("geoip lookup previously failed for ip %s", ip)
		}
		if found {
			var geoData GeoIP
			if err := json.Unmarshal(cached, &geoData); err == nil {
				return &geoData, nil
			}
		}
	}

	// Use a bounded timeout for external GeoIP lookups to prevent slow responses
	// from blocking the queue handler.
	geoCtx, cancel := context.WithTimeout(ctx, geoIPRequestTimeout)
	defer cancel()

	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	resp, err := cli.Invoke(geoCtx, http.MethodGet, url, nil, nil)
	if err != nil {
		// Negative-cache on network errors to avoid retrying rapidly.
		if cacheSvc != nil {
			cacheSvc.SetGeoIPNegative(ctx, ip)
		}
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		result, _ := resp.ToContent(geoCtx)
		// Negative-cache non-OK responses (rate limits, bad IPs, etc.).
		if cacheSvc != nil {
			cacheSvc.SetGeoIPNegative(ctx, ip)
		}
		return nil, fmt.Errorf("unexpected status code: %d : %s", resp.StatusCode, string(result))
	}

	var geoData GeoIP
	err = resp.Decode(geoCtx, &geoData)
	if err != nil {
		return nil, err
	}

	// Cache the successful result.
	if cacheSvc != nil {
		encoded, encErr := json.Marshal(&geoData)
		if encErr == nil {
			cacheSvc.SetGeoIP(ctx, ip, encoded)
		} else {
			util.Log(ctx).WithError(encErr).Debug("failed to marshal geoip for cache")
		}
	}

	return &geoData, nil
}
