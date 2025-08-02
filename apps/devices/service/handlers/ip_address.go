package handlers

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// GetClientIP extracts the client's IP address from the context of a gRPC call.
func GetClientIP(ctx context.Context) string {
	// 1. Check for proxy headers in gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// gRPC metadata keys are automatically converted to lowercase
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			// X-Forwarded-For can be a comma-separated list, take the first one
			ips := strings.Split(xff[0], ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		if xrip := md.Get("x-real-ip"); len(xrip) > 0 {
			return strings.TrimSpace(xrip[0])
		}
	}

	// 2. Fallback to peer information
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}

	// 3. Check peer information from gRPC context
	if tcpAddr, tcpOk := p.Addr.(*net.TCPAddr); tcpOk {
		return tcpAddr.IP.String()
	}

	return ""
}
