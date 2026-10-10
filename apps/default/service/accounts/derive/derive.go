// Vendored from github.com/stawilabs/stawi/pkg/protocol/derive at tag
// pkg/protocol/derive/v0.1.0 (commit a9d937d). Do not edit here: change the
// stawi source, retag, and copy it again. Its correctness is pinned by the
// contract-generated vector in ../testdata/account_derivation.json
// (contracts_vector_test.go) and the upstream fixed-input tests.

// Package derive implements account address derivation (architecture §8.1,
// P18) and the CREATE2 computation shared by finance, settlement, the SDK and
// deploy scripts.
//
//	identity_salt_hash = keccak256(identity_salt)
//	ACCOUNT_SALT       = keccak256("STAWI_ACCOUNT_V1" ‖ bytes32 account_family ‖ uint32 account_version ‖ bytes32 identity_salt_hash)
//	address            = CREATE2(canonical_factory, ACCOUNT_SALT, keccak256(init_code(identity_salt_hash)))
//
// Only the hash of the identity salt ever reaches the chain or the factory
// calldata; the salt itself stays in finance and the user's recovery kit. The
// address depends only on family, version and the salt hash. Nothing mutable
// (policy, keys, chain) enters the derivation; chain equality is a deployment
// property (identical factory address and init code), recorded per chain in
// the protocol manifest.
package derive

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/sha3"
)

const saltTag = "STAWI_ACCOUNT_V1"

var (
	ErrBadFamily  = errors.New("derive: account_family must be 32 bytes")
	ErrBadSalt    = errors.New("derive: identity_salt must be 32 non-zero bytes")
	ErrBadFactory = errors.New("derive: factory address must be 20 bytes")
)

// Family is a 32-byte account family identifier; FamilyEVM is the only one in v1.
type Family [32]byte

// FamilyEVM is keccak256("EVM"), so that families are content-addressed and
// cannot collide with hand-written constants.
var FamilyEVM = Family(keccak([]byte("EVM")))

// IdentitySaltHash is keccak256 of the raw identity salt; it is the only form
// of the salt that leaves the finance service.
func IdentitySaltHash(identitySalt [32]byte) ([32]byte, error) {
	if identitySalt == ([32]byte{}) {
		return [32]byte{}, ErrBadSalt
	}
	return keccak(identitySalt[:]), nil
}

// AccountSalt computes ACCOUNT_SALT from the identity salt hash.
func AccountSalt(family Family, accountVersion uint32, identitySaltHash [32]byte) ([32]byte, error) {
	if identitySaltHash == ([32]byte{}) {
		return [32]byte{}, ErrBadSalt
	}
	// account_version is a packed uint32 (4 bytes, big-endian), matching
	// abi.encodePacked(uint32) in AccountFactory.accountSalt.
	pre := make([]byte, 0, len(saltTag)+32+4+32)
	pre = append(pre, saltTag...)
	pre = append(pre, family[:]...)
	pre = binary.BigEndian.AppendUint32(pre, accountVersion)
	pre = append(pre, identitySaltHash[:]...)
	return keccak(pre), nil
}

// Create2 computes the EVM CREATE2 address:
// keccak256(0xff ‖ deployer ‖ salt ‖ initCodeHash)[12:].
func Create2(deployer [20]byte, salt [32]byte, initCodeHash [32]byte) [20]byte {
	pre := make([]byte, 0, 1+20+32+32)
	pre = append(pre, 0xff)
	pre = append(pre, deployer[:]...)
	pre = append(pre, salt[:]...)
	pre = append(pre, initCodeHash[:]...)
	h := keccak(pre)
	var out [20]byte
	copy(out[:], h[12:])
	return out
}

// AccountAddress derives the account address for an identity on a family and
// version, given the canonical factory and the keccak of the canonical init
// code for that version (which is a function of the identity salt hash only).
func AccountAddress(family Family, accountVersion uint32, identitySaltHash [32]byte, factory [20]byte, initCodeHash [32]byte) ([20]byte, error) {
	salt, err := AccountSalt(family, accountVersion, identitySaltHash)
	if err != nil {
		return [20]byte{}, err
	}
	return Create2(factory, salt, initCodeHash), nil
}

// InitCodeHash returns keccak256(creationCode ‖ abi.encode(identitySaltHash)).
// The StawiAccount constructor takes exactly one argument, the identity salt
// hash, so that the salt itself never appears in calldata. The creation code
// is version-fixed and published in the manifest.
func InitCodeHash(creationCode []byte, identitySaltHash [32]byte) [32]byte {
	pre := make([]byte, 0, len(creationCode)+32)
	pre = append(pre, creationCode...)
	pre = append(pre, identitySaltHash[:]...) // abi.encode(bytes32) is the 32 bytes as-is
	return keccak(pre)
}

func keccak(b []byte) [32]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(b)
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}
