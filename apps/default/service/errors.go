package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnspecifiedID          = status.Error(codes.InvalidArgument, "No id was supplied")
	ErrEmptyValueSupplied     = status.Error(codes.InvalidArgument, "Empty value supplied")
	ErrItemExist              = status.Error(codes.AlreadyExists, "Specified item already exists")
	ErrItemDoesNotExist       = status.Error(codes.NotFound, "Specified item does not exist")
	ErrContactTypeNotValid    = status.Error(codes.InvalidArgument, "Contact type is not known ")
	ErrContactDetailsNotValid = status.Error(codes.InvalidArgument, "Contact details are invalid ")
	ErrContactProfileNotValid = status.Error(codes.InvalidArgument, "Contact profile is invalid ")

	ErrProfileDoesNotExist = status.Error(codes.NotFound, "Specified profile does not exist")
	ErrContactDoesNotExist = status.Error(codes.NotFound, "Specified contact does not exist")
	ErrAddressDoesNotExist = status.Error(codes.NotFound, "Specified address does not exist")
	ErrCountryDoesNotExist = status.Error(codes.NotFound, "Specified country does not exist")
)
