package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strings"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	envToken := strings.TrimSpace(os.Getenv("PROFILE_AUTH_TOKEN"))
	if envToken != "" {

		requestToken := ctx.Value("token")
		if requestToken == nil {
			return nil, status.Errorf(codes.Unauthenticated, "incorrect access token")
		}
		token := requestToken.(string)
		// Validate token
		if token != envToken {
			return nil, status.Errorf(codes.Unauthenticated, "incorrect access token")
		}

	}

	return handler(ctx, req)
}
