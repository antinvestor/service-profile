package handlers

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/pitabwire/frame/data"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Cloudflare header keys (lowercase, as gRPC metadata normalises to lowercase).
const (
	headerCFConnectingIP = "cf-connecting-ip"
	headerCFIPCountry    = "cf-ipcountry"
	headerCFIPCity       = "cf-ipcity"
	headerCFIPContinent  = "cf-ipcontinent"
	headerCFIPLatitude   = "cf-iplatitude"
	headerCFIPLongitude  = "cf-iplongitude"
	headerCFRegion       = "cf-region"
	headerCFRegionCode   = "cf-region-code"
	headerCFPostalCode   = "cf-postal-code"
	headerCFTimezone     = "cf-timezone"
)

// GetClientIP extracts the client's IP address from the context of a gRPC call.
// Prefers Cloudflare's CF-Connecting-IP (set by the edge and not spoofable by clients)
// over X-Forwarded-For.
func GetClientIP(ctx context.Context) string {
	if ip := getIPFromMetadata(ctx); ip != "" {
		return ip
	}
	return getIPFromPeer(ctx)
}

// getIPFromMetadata extracts the client IP from gRPC metadata headers.
// Priority: CF-Connecting-IP > X-Forwarded-For > X-Real-IP.
func getIPFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// CF-Connecting-IP is injected by the Cloudflare edge — most trustworthy.
	if ip := firstMDValue(md, headerCFConnectingIP); ip != "" {
		return ip
	}

	// X-Forwarded-For — take the leftmost (original client) IP.
	if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
		if parts := strings.SplitN(xff[0], ",", 2); len(parts) > 0 { //nolint:mnd // Split into at most 2 parts.
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}

	// X-Real-IP — set by some reverse proxies.
	return firstMDValue(md, "x-real-ip")
}

// getIPFromPeer extracts the IP from the gRPC peer transport address.
func getIPFromPeer(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	if tcpAddr, tcpOk := p.Addr.(*net.TCPAddr); tcpOk {
		return tcpAddr.IP.String()
	}
	return ""
}

// firstMDValue returns the first non-empty trimmed value for a metadata key.
func firstMDValue(md metadata.MD, key string) string {
	vals := md.Get(key)
	if len(vals) > 0 {
		if v := strings.TrimSpace(vals[0]); v != "" {
			return v
		}
	}
	return ""
}

// ExtractCloudflareGeo reads Cloudflare geolocation headers from either
// Connect HTTP headers or gRPC metadata and returns them as a JSONMap.
//
// When the service runs behind Cloudflare, these headers provide geo data
// at the edge with zero latency, avoiding external GeoIP API calls.
//
// Cloudflare geo headers reference:
// https://developers.cloudflare.com/fundamentals/reference/http-request-headers/
func ExtractCloudflareGeo(ctx context.Context, httpHeader http.Header) data.JSONMap {
	geo := data.JSONMap{}

	// Try Connect/HTTP headers first (preferred — they come directly from the request).
	if httpHeader != nil {
		setIfPresent(geo, "cf_country", httpHeader.Get("Cf-Ipcountry"))
		setIfPresent(geo, "cf_city", httpHeader.Get("Cf-Ipcity"))
		setIfPresent(geo, "cf_continent", httpHeader.Get("Cf-Ipcontinent"))
		setIfPresent(geo, "cf_region", httpHeader.Get("Cf-Region"))
		setIfPresent(geo, "cf_region_code", httpHeader.Get("Cf-Region-Code"))
		setIfPresent(geo, "cf_postal_code", httpHeader.Get("Cf-Postal-Code"))
		setIfPresent(geo, "cf_timezone", httpHeader.Get("Cf-Timezone"))
		setFloatIfPresent(geo, "cf_latitude", httpHeader.Get("Cf-Iplatitude"))
		setFloatIfPresent(geo, "cf_longitude", httpHeader.Get("Cf-Iplongitude"))

		if len(geo) > 0 {
			return geo
		}
	}

	// Fallback to gRPC metadata (keys are lowercase).
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return geo
	}

	setMDIfPresent(geo, "cf_country", md, headerCFIPCountry)
	setMDIfPresent(geo, "cf_city", md, headerCFIPCity)
	setMDIfPresent(geo, "cf_continent", md, headerCFIPContinent)
	setMDIfPresent(geo, "cf_region", md, headerCFRegion)
	setMDIfPresent(geo, "cf_region_code", md, headerCFRegionCode)
	setMDIfPresent(geo, "cf_postal_code", md, headerCFPostalCode)
	setMDIfPresent(geo, "cf_timezone", md, headerCFTimezone)
	setMDFloatIfPresent(geo, "cf_latitude", md, headerCFIPLatitude)
	setMDFloatIfPresent(geo, "cf_longitude", md, headerCFIPLongitude)

	return geo
}

func setIfPresent(m data.JSONMap, key, val string) {
	v := strings.TrimSpace(val)
	if v != "" {
		m[key] = v
	}
}

func setFloatIfPresent(m data.JSONMap, key, val string) {
	v := strings.TrimSpace(val)
	if v == "" {
		return
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		m[key] = f
	}
}

func setMDIfPresent(m data.JSONMap, key string, md metadata.MD, headerKey string) {
	vals := md.Get(headerKey)
	if len(vals) > 0 {
		setIfPresent(m, key, vals[0])
	}
}

func setMDFloatIfPresent(m data.JSONMap, key string, md metadata.MD, headerKey string) {
	vals := md.Get(headerKey)
	if len(vals) > 0 {
		setFloatIfPresent(m, key, vals[0])
	}
}
