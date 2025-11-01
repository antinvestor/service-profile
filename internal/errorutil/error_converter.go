package errorutil

import (
	"github.com/pitabwire/frame/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrToAPI(err error) error {
	grpcError, ok := status.FromError(err)

	if ok {
		return grpcError.Err()
	}

	if data.ErrorIsNoRows(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	return grpcError.Err()
}
