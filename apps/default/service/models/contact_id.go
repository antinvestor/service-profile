package models

import "github.com/rs/xid"

// ContactIDFromLookupToken derives a contact's id from its lookup token, so
// the same normalised detail always gets the same id.
//
// The token is the keyed HMAC-SHA256 of the detail, so the id cannot be
// reversed to the phone number or email by hashing guesses. Its first 12
// bytes are encoded the way an xid is: 20 characters of base32hex, the same
// shape as every other id frame generates, so the id passes xid validation.
// The xid "timestamp" in a derived id is meaningless; nothing reads it.
//
// A token shorter than 12 bytes yields "", leaving frame to generate a
// random id.
func ContactIDFromLookupToken(token []byte) string {
	var id xid.ID
	if len(token) < len(id) {
		return ""
	}
	copy(id[:], token)
	return id.String()
}
