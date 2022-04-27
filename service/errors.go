package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorUnspecifiedID          = status.Error(codes.InvalidArgument, "No id was supplied")
	ErrorEmptyValueSupplied     = status.Error(codes.InvalidArgument, "Empty value supplied")
	ErrorItemExist              = status.Error(codes.AlreadyExists, "Specified item already exists")
	ErrorItemDoesNotExist       = status.Error(codes.NotFound, "Specified item does not exist")
	ErrorContactDetailsNotValid = status.Error(codes.InvalidArgument, "Contact details are invalid ")

	ErrorProfileDoesNotExist = status.Error(codes.NotFound, "Specified profile does not exist")

	ErrorContactDoesNotExist = status.Error(codes.NotFound, "Specified contact does not exist")

	ErrorAddressDoesNotExist = status.Error(codes.NotFound, "Specified address does not exist")
	ErrorCountryDoesNotExist = status.Error(codes.NotFound, "Specified country does not exist")
)
