// Vendored from github.com/stawilabs/stawi/pkg/protocol/derive at tag
// pkg/protocol/derive/v0.1.0 (commit a9d937d). Do not edit here: change the
// stawi source, retag, and copy it again. Its correctness is pinned by the
// contract-generated vector in ../testdata/account_derivation.json
// (contracts_vector_test.go) and the upstream fixed-input tests.

package derive

import (
	"encoding/hex"
	"testing"
)

func salt(b byte) [32]byte {
	var s [32]byte
	for i := range s {
		s[i] = b
	}
	return s
}

// Known-answer test for CREATE2 from EIP-1014 example 1:
// address 0x0000000000000000000000000000000000000000, salt 0x00..00, init_code 0x00 → 0x4D1A2e2bB4F88F0250f26Ffff098B0b30B26BF38
func TestCreate2KnownAnswer(t *testing.T) {
	initHash := keccak([]byte{0x00})
	got := Create2([20]byte{}, [32]byte{}, initHash)
	want, _ := hex.DecodeString("4D1A2e2bB4F88F0250f26Ffff098B0b30B26BF38")
	if hex.EncodeToString(got[:]) != hex.EncodeToString(want) {
		t.Fatalf("got %x want %x", got, want)
	}
}

// P18: the address depends only on family, version and identity salt.
func TestAddressDeterminism(t *testing.T) {
	factory := [20]byte{0xfa}
	code := []byte{0x60, 0x80, 0x60, 0x40}
	s, _ := IdentitySaltHash(salt(0x11))
	a1, err := AccountAddress(FamilyEVM, 1, s, factory, InitCodeHash(code, s))
	if err != nil {
		t.Fatal(err)
	}
	a2, _ := AccountAddress(FamilyEVM, 1, s, factory, InitCodeHash(code, s))
	if a1 != a2 {
		t.Fatal("not deterministic")
	}
	// different salt -> different address
	s2, _ := IdentitySaltHash(salt(0x12))
	a3, _ := AccountAddress(FamilyEVM, 1, s2, factory, InitCodeHash(code, s2))
	if a3 == a1 {
		t.Fatal("salt not bound")
	}
	// different version -> different address
	a4, _ := AccountAddress(FamilyEVM, 2, s, factory, InitCodeHash(code, s))
	if a4 == a1 {
		t.Fatal("version not bound")
	}
	// different factory -> different address (cross-chain equality requires the same factory)
	a5, _ := AccountAddress(FamilyEVM, 1, s, [20]byte{0xfb}, InitCodeHash(code, s))
	if a5 == a1 {
		t.Fatal("factory not bound")
	}
}

func TestSaltRejectsZero(t *testing.T) {
	if _, err := AccountSalt(FamilyEVM, 1, [32]byte{}); err != ErrBadSalt {
		t.Fatalf("expected ErrBadSalt, got %v", err)
	}
	if _, err := IdentitySaltHash([32]byte{}); err != ErrBadSalt {
		t.Fatalf("expected ErrBadSalt for zero salt, got %v", err)
	}
}

func TestFamilyIsContentAddressed(t *testing.T) {
	if FamilyEVM != Family(keccak([]byte("EVM"))) {
		t.Fatal("family constant drifted")
	}
}

// TestCrossLanguageVector pins the salt that packages/contracts
// test/account/AccountFactory.t.sol::test_derivation_vector asserts, so a
// change to either side's encoding is caught in CI.
func TestCrossLanguageVector(t *testing.T) {
	h := keccak([]byte("vector"))
	if got := hex.EncodeToString(h[:]); got != "21c55cde1a8b741b30d8e78ab6d05799cd8b24f366a420a4049982f13704a49c" {
		t.Fatalf("identity salt hash: %s", got)
	}
	salt, err := AccountSalt(FamilyEVM, 1, h)
	if err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(salt[:]); got != "b4c8b739036c30d37fb53af39052941f72a4f0826b96ce98c88dfea82941d47f" {
		t.Fatalf("account salt: %s", got)
	}
}
