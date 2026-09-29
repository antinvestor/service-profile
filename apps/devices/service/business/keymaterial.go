// Copyright 2023-2026 Ant Investor Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package business

import (
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	devicev1 "buf.build/gen/go/antinvestor/device/protocolbuffers/go/device/v1"
	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/data"
	"golang.org/x/crypto/sha3"
)

// Cryptographic device key types required by the Stawi Group Financial OS
// (GFOS §8.3, platform change K4). The generated Go enum constants land in
// this service only once buf.build/antinvestor/device is published from
// proto/device/device/v1/device.proto, so the wire numbers are named here and
// stay the contract in the meantime.
const (
	// KeyTypeSecp256k1PublicKey is device.v1.KeyType.SECP256K1_PUBLIC_KEY.
	KeyTypeSecp256k1PublicKey devicev1.KeyType = 6
	// KeyTypeP256WebauthnPublicKey is device.v1.KeyType.P256_WEBAUTHN_PUBLIC_KEY.
	KeyTypeP256WebauthnPublicKey devicev1.KeyType = 7
)

// Wire names of the key types as they appear in published facts. GFOS matches
// on these exact strings (its Credential.KeyType column).
const (
	KeyTypeNameSecp256k1    = "secp256k1"
	KeyTypeNameP256WebAuthn = "p256-webauthn"
)

// Byte lengths of the accepted public key encodings.
const (
	secp256k1UncompressedLen = 65
	p256RawCoordinateLen     = 64
	coordinateLen            = 32
	keyIDLen                 = 20

	// sec1UncompressedPrefix marks a SEC1 uncompressed EC point.
	sec1UncompressedPrefix = 0x04
	// hexBase is the radix of the curve constants below.
	hexBase = 16
	// derivedFieldCount is how many fields this service adds to extra.
	derivedFieldCount = 4

	maxCredentialIDLen = 1023
	minCredentialIDLen = 16
	maxRPIDLen         = 253
	maxLabelLen        = 64
	maxSessionProofLen = 128

	// coseAlgES256 is the only COSE algorithm a P-256 WebAuthn credential may
	// declare; the GFOS verifier (pkg/p256) implements ES256 and nothing else.
	coseAlgES256 = -7
)

// ErrInvalidKeyMaterial is returned when key bytes or the extra document do
// not satisfy the schema for the declared key type.
var ErrInvalidKeyMaterial = errors.New("invalid key material")

// IsPublicKeyType reports whether a key type holds public key material that is
// safe to carry in a published fact. Every other type may hold a secret (an
// FCM token, a pickle key, a Matrix session key) and its bytes never leave the
// database.
func IsPublicKeyType(keyType devicev1.KeyType) bool {
	switch keyType {
	case KeyTypeSecp256k1PublicKey, KeyTypeP256WebauthnPublicKey,
		devicev1.KeyType_ED25519_KEY, devicev1.KeyType_CURVE25519_KEY:
		return true
	case devicev1.KeyType_MATRIX_KEY, devicev1.KeyType_NOTIFICATION_KEY,
		devicev1.KeyType_FCM_TOKEN, devicev1.KeyType_PICKLE_KEY:
		return false
	default:
		return false
	}
}

// KeyTypeName renders a key type for the fact catalogue. Cryptographic types
// use the names GFOS matches on; everything else uses the proto enum name.
func KeyTypeName(keyType devicev1.KeyType) string {
	switch keyType {
	case KeyTypeSecp256k1PublicKey:
		return KeyTypeNameSecp256k1
	case KeyTypeP256WebauthnPublicKey:
		return KeyTypeNameP256WebAuthn
	case devicev1.KeyType_MATRIX_KEY, devicev1.KeyType_NOTIFICATION_KEY,
		devicev1.KeyType_FCM_TOKEN, devicev1.KeyType_CURVE25519_KEY,
		devicev1.KeyType_ED25519_KEY, devicev1.KeyType_PICKLE_KEY:
		return keyType.String()
	default:
		return keyType.String()
	}
}

// ValidateKeyMaterial checks the key bytes and the extra document against the
// schema for keyType and returns the extra document the service will store,
// with the derived fields (key_id, curve, encoding, key_sha256) filled in.
//
// Key types that predate GFOS are left untouched: their material is opaque to
// this service and their extra document is free-form.
//
// # SECP256K1_PUBLIC_KEY
//
// key: 65 bytes, SEC1 uncompressed (0x04 ‖ x ‖ y), a point on secp256k1.
// Compressed keys are rejected — the GFOS verifier only reads the
// uncompressed form.
//
// extra, client supplied:
//
//	key_id            string  optional  "0x" + 40 hex. Must equal the derived id when present.
//	label             string  optional  free-form, ≤64 chars.
//	session_proof_id  string  optional  ≤128 chars; the authentication-service session
//	                                    proof (GFOS K9) this registration happened under.
//
// extra, always written by this service:
//
//	key_id      string  "0x" + hex(keccak256(x ‖ y)[12:]) — the EVM address GFOS
//	                    resolves as device_key_id.
//	curve       string  "secp256k1"
//	encoding    string  "sec1-uncompressed"
//	key_sha256  string  hex sha256 of the stored key bytes.
//
// # P256_WEBAUTHN_PUBLIC_KEY
//
// key: 64 bytes, raw coordinates (x ‖ y, no 0x04 prefix), a point on P-256.
// This is the form a WebAuthn COSE key decodes to and the form the GFOS
// verifier takes.
//
// extra, client supplied:
//
//	credential_id     string   required  base64url (unpadded), 16..1023 decoded bytes.
//	rp_id             string   required  relying party id, ≤253 chars.
//	cose_alg          number   required  must be -7 (ES256).
//	aaguid            string   optional  36-char UUID.
//	sign_count        number   optional  ≥ 0.
//	transports        []string optional  subset of usb, nfc, ble, smart-card,
//	                                     hybrid, internal, cable.
//	user_verification string   optional  required | preferred | discouraged.
//	backup_eligible   bool     optional
//	backup_state      bool     optional  may only be true when backup_eligible is true.
//	origin            string   optional  ≤255 chars.
//	label             string   optional  free-form, ≤64 chars.
//	session_proof_id  string   optional  ≤128 chars.
//
// extra, always written by this service:
//
//	key_id      string  "0x" + hex(keccak256(x ‖ y)[12:]) — the key id the GFOS
//	                    account contract derives.
//	curve       string  "p-256"
//	encoding    string  "raw-xy"
//	key_sha256  string  hex sha256 of the stored key bytes.
//
// No field of either schema ever holds private key material; a caller that
// sends one is rejected by the unknown-field check.
func ValidateKeyMaterial(
	keyType devicev1.KeyType,
	key []byte,
	extra data.JSONMap,
) (data.JSONMap, error) {
	switch keyType {
	case KeyTypeSecp256k1PublicKey:
		return validateSecp256k1(key, extra)
	case KeyTypeP256WebauthnPublicKey:
		return validateP256WebAuthn(key, extra)
	case devicev1.KeyType_MATRIX_KEY, devicev1.KeyType_NOTIFICATION_KEY,
		devicev1.KeyType_FCM_TOKEN, devicev1.KeyType_CURVE25519_KEY,
		devicev1.KeyType_ED25519_KEY, devicev1.KeyType_PICKLE_KEY:
		fallthrough
	default:
		if len(key) == 0 {
			return nil, badKeyMaterial("key material is empty")
		}
		return extra, nil
	}
}

func badKeyMaterial(format string, args ...any) error {
	return connect.NewError(
		connect.CodeInvalidArgument,
		fmt.Errorf("%w: %s", ErrInvalidKeyMaterial, fmt.Sprintf(format, args...)),
	)
}

// secp256k1FieldPrime is the field prime of the secp256k1 curve.
const secp256k1FieldPrime = "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"

// secp256k1CurveB is the curve's b parameter in y² = x³ + b.
const secp256k1CurveB = 7

func validateSecp256k1(key []byte, extra data.JSONMap) (data.JSONMap, error) {
	if len(key) != secp256k1UncompressedLen {
		return nil, badKeyMaterial(
			"a secp256k1 public key is %d bytes (0x04 ‖ x ‖ y), got %d",
			secp256k1UncompressedLen, len(key),
		)
	}
	if key[0] != sec1UncompressedPrefix {
		return nil, badKeyMaterial(
			"a secp256k1 public key must be SEC1 uncompressed (0x04 prefix), got 0x%02x", key[0],
		)
	}
	if err := onSecp256k1Curve(key[1:coordinateLen+1], key[coordinateLen+1:]); err != nil {
		return nil, err
	}

	out, err := copyExtra(extra, []string{"key_id", "label", "session_proof_id"})
	if err != nil {
		return nil, err
	}
	if err = validateCommonExtra(out); err != nil {
		return nil, err
	}

	keyID := deriveKeyID(key[1:])
	if declared, ok := out["key_id"].(string); ok && !strings.EqualFold(declared, keyID) {
		return nil, badKeyMaterial("key_id %q does not match the key material", declared)
	}

	out["key_id"] = keyID
	out["curve"] = "secp256k1"
	out["encoding"] = "sec1-uncompressed"
	out["key_sha256"] = keySHA256(key)
	return out, nil
}

// onSecp256k1Curve checks that (x, y) is an affine point of secp256k1. The
// curve has no standard library implementation, so the check is done directly:
// both coordinates below the field prime and y² ≡ x³ + 7 (mod p).
func onSecp256k1Curve(xb, yb []byte) error {
	secp256k1P, ok := new(big.Int).SetString(secp256k1FieldPrime, hexBase)
	if !ok {
		return errors.New("secp256k1 field prime is malformed")
	}
	secp256k1B := big.NewInt(secp256k1CurveB)

	x := new(big.Int).SetBytes(xb)
	y := new(big.Int).SetBytes(yb)
	if x.Sign() == 0 && y.Sign() == 0 {
		return badKeyMaterial("a secp256k1 public key may not be the point at infinity")
	}
	if x.Cmp(secp256k1P) >= 0 || y.Cmp(secp256k1P) >= 0 {
		return badKeyMaterial("secp256k1 coordinates must be below the field prime")
	}

	lhs := new(big.Int).Mul(y, y)
	lhs.Mod(lhs, secp256k1P)

	rhs := new(big.Int).Mul(x, x)
	rhs.Mod(rhs, secp256k1P)
	rhs.Mul(rhs, x)
	rhs.Add(rhs, secp256k1B)
	rhs.Mod(rhs, secp256k1P)

	if lhs.Cmp(rhs) != 0 {
		return badKeyMaterial("the supplied coordinates are not a point on secp256k1")
	}
	return nil
}

func validateP256WebAuthn(key []byte, extra data.JSONMap) (data.JSONMap, error) {
	if len(key) != p256RawCoordinateLen {
		return nil, badKeyMaterial(
			"a p256-webauthn public key is %d bytes (x ‖ y, no 0x04 prefix), got %d",
			p256RawCoordinateLen, len(key),
		)
	}
	// crypto/ecdh takes the SEC1 uncompressed form and rejects any point that
	// is not on P-256, which is exactly the check needed here.
	uncompressed := make([]byte, 0, p256RawCoordinateLen+1)
	uncompressed = append(uncompressed, sec1UncompressedPrefix)
	uncompressed = append(uncompressed, key...)
	if _, err := ecdh.P256().NewPublicKey(uncompressed); err != nil {
		return nil, badKeyMaterial("the supplied coordinates are not a point on P-256")
	}

	out, err := copyExtra(extra, []string{
		"credential_id", "rp_id", "cose_alg", "aaguid", "sign_count", "transports",
		"user_verification", "backup_eligible", "backup_state", "origin",
		"key_id", "label", "session_proof_id",
	})
	if err != nil {
		return nil, err
	}
	if err = validateCommonExtra(out); err != nil {
		return nil, err
	}
	if err = validateWebAuthnRequired(out); err != nil {
		return nil, err
	}
	if err = validateWebAuthnOptional(out); err != nil {
		return nil, err
	}

	keyID := deriveKeyID(key)
	if declared, ok := out["key_id"].(string); ok && !strings.EqualFold(declared, keyID) {
		return nil, badKeyMaterial("key_id %q does not match the key material", declared)
	}

	out["key_id"] = keyID
	out["curve"] = "p-256"
	out["encoding"] = "raw-xy"
	out["key_sha256"] = keySHA256(key)
	return out, nil
}

func validateWebAuthnRequired(out data.JSONMap) error {
	credentialID, ok := out["credential_id"].(string)
	if !ok || credentialID == "" {
		return badKeyMaterial("credential_id is required for a p256-webauthn key")
	}
	raw, err := decodeBase64URL(credentialID)
	if err != nil {
		return badKeyMaterial("credential_id must be base64url encoded")
	}
	if len(raw) < minCredentialIDLen || len(raw) > maxCredentialIDLen {
		return badKeyMaterial(
			"credential_id must decode to %d..%d bytes, got %d",
			minCredentialIDLen, maxCredentialIDLen, len(raw),
		)
	}

	rpID, ok := out["rp_id"].(string)
	if !ok || rpID == "" {
		return badKeyMaterial("rp_id is required for a p256-webauthn key")
	}
	if len(rpID) > maxRPIDLen {
		return badKeyMaterial("rp_id must be at most %d characters", maxRPIDLen)
	}

	alg, ok := numeric(out["cose_alg"])
	if !ok {
		return badKeyMaterial("cose_alg is required for a p256-webauthn key")
	}
	if alg != coseAlgES256 {
		return badKeyMaterial("cose_alg must be %d (ES256), got %v", coseAlgES256, alg)
	}
	out["cose_alg"] = coseAlgES256
	return nil
}

func validateWebAuthnOptional(out data.JSONMap) error {
	if err := validateAuthenticatorFields(out); err != nil {
		return err
	}
	if err := validateTransportFields(out); err != nil {
		return err
	}

	eligible, err := boolField(out, "backup_eligible")
	if err != nil {
		return err
	}
	state, err := boolField(out, "backup_state")
	if err != nil {
		return err
	}
	if state && !eligible {
		return badKeyMaterial("backup_state may only be true when backup_eligible is true")
	}
	if v, present := out["origin"]; present {
		s, ok := v.(string)
		if !ok || len(s) > 255 {
			return badKeyMaterial("origin must be a string of at most 255 characters")
		}
	}
	return nil
}

// validateAuthenticatorFields checks the fields that describe the
// authenticator itself.
func validateAuthenticatorFields(out data.JSONMap) error {
	if v, present := out["aaguid"]; present {
		s, ok := v.(string)
		if !ok || len(s) != len("00000000-0000-0000-0000-000000000000") {
			return badKeyMaterial("aaguid must be a 36 character UUID")
		}
	}
	if v, present := out["sign_count"]; present {
		n, ok := numeric(v)
		if !ok || n < 0 {
			return badKeyMaterial("sign_count must be a non negative number")
		}
		out["sign_count"] = n
	}
	return nil
}

// validateTransportFields checks how the authenticator may be reached and what
// verification it performs.
func validateTransportFields(out data.JSONMap) error {
	webAuthnTransports := []string{"usb", "nfc", "ble", "smart-card", "hybrid", "internal", "cable"}
	webAuthnUserVerification := []string{"required", "preferred", "discouraged"}

	if v, present := out["transports"]; present {
		list, ok := v.([]any)
		if !ok {
			return badKeyMaterial("transports must be an array of strings")
		}
		for _, item := range list {
			s, isStr := item.(string)
			if !isStr || !slices.Contains(webAuthnTransports, s) {
				return badKeyMaterial("transports must be a subset of %v", webAuthnTransports)
			}
		}
	}
	if v, present := out["user_verification"]; present {
		s, ok := v.(string)
		if !ok || !slices.Contains(webAuthnUserVerification, s) {
			return badKeyMaterial("user_verification must be one of %v", webAuthnUserVerification)
		}
	}
	return nil
}

func boolField(out data.JSONMap, name string) (bool, error) {
	v, present := out[name]
	if !present {
		return false, nil
	}
	b, ok := v.(bool)
	if !ok {
		return false, badKeyMaterial("%s must be a boolean", name)
	}
	return b, nil
}

func validateCommonExtra(out data.JSONMap) error {
	if v, present := out["label"]; present {
		s, ok := v.(string)
		if !ok || len(s) > maxLabelLen {
			return badKeyMaterial("label must be a string of at most %d characters", maxLabelLen)
		}
	}
	if v, present := out["session_proof_id"]; present {
		s, ok := v.(string)
		if !ok || len(s) > maxSessionProofLen {
			return badKeyMaterial(
				"session_proof_id must be a string of at most %d characters", maxSessionProofLen,
			)
		}
	}
	if v, present := out["key_id"]; present {
		if _, ok := v.(string); !ok {
			return badKeyMaterial("key_id must be a string")
		}
	}
	return nil
}

// copyExtra returns a copy of extra and fails when it carries a field the
// schema does not define. A closed schema is what makes the document safe to
// publish: nothing unreviewed, and in particular nothing secret, can ride
// along on a key registration.
func copyExtra(extra data.JSONMap, allowed []string) (data.JSONMap, error) {
	out := make(data.JSONMap, len(extra)+derivedFieldCount)
	for k, v := range extra {
		if !slices.Contains(allowed, k) {
			// Derived fields are recomputed, so accept (and overwrite) them.
			if k == "curve" || k == "encoding" || k == "key_sha256" {
				continue
			}
			return nil, badKeyMaterial("extra field %q is not part of the schema for this key type", k)
		}
		out[k] = v
	}
	return out, nil
}

// deriveKeyID is keccak256(x ‖ y)[12:] rendered as "0x" + 40 hex characters.
// For secp256k1 that is the EVM address of the key; for P-256 it is the key id
// the Stawi account contract derives. GFOS resolves both as device_key_id.
func deriveKeyID(coordinates []byte) string {
	h := sha3.NewLegacyKeccak256()
	h.Write(coordinates)
	sum := h.Sum(nil)
	return "0x" + hex.EncodeToString(sum[len(sum)-keyIDLen:])
}

func keySHA256(key []byte) string {
	sum := sha256.Sum256(key)
	return hex.EncodeToString(sum[:])
}

func decodeBase64URL(s string) ([]byte, error) {
	if raw, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		return raw, nil
	}
	return base64.URLEncoding.DecodeString(s)
}

func numeric(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return 0, false
	}
}
