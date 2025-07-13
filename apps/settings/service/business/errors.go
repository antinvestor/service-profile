package business

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnspecifiedID      = status.Error(codes.InvalidArgument, "No id was supplied")
	ErrEmptyValueSupplied = status.Error(codes.InvalidArgument, "Empty value supplied")
	ErrItemExist          = status.Error(codes.AlreadyExists, "Specified item already exists")
	ErrItemDoesNotExist   = status.Error(codes.NotFound, "Specified item does not exist")

	ErrInitializationFail = status.Error(codes.Internal, "Internal configuration is invalid")
)
