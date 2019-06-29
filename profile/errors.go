package profile

import (
	"fmt"
)

type ErrorInProfile struct {
	code    int
	message string
}

func (e *ErrorInProfile) Error() string {
	return fmt.Sprintf("%d  : - %s", e.code, e.message)
}

var (
	ErrorUnspecifiedID          = &ErrorInProfile{code: 1, message: "No id was supplied"}
	ErrorEmptyValueSupplied     = &ErrorInProfile{code: 2, message: "Empty value supplied"}
	ErrorItemExist              = &ErrorInProfile{3, "Specified item already exists"}
	ErrorItemDoesNotExist       = &ErrorInProfile{4, "Specified item does not exist"}
	ErrorContactDetailsNotValid = &ErrorInProfile{5, "Contact details are invalid "}

	ErrorProfileDoesNotExist = &ErrorInProfile{20, "Specified profile does not exist"}

	ErrorContactDoesNotExist = &ErrorInProfile{30, "Specified contact does not exist"}

	ErrorAddressDoesNotExist = &ErrorInProfile{40, "Specified address does not exist"}
	ErrorCountryDoesNotExist = &ErrorInProfile{41, "Specified country does not exist"}

)
