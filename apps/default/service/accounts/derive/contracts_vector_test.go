// Vendored from github.com/stawilabs/stawi/pkg/protocol/derive at tag
// pkg/protocol/derive/v0.1.0 (commit a9d937d). Do not edit here: change the
// stawi source, retag, and copy it again. Its correctness is pinned by the
// contract-generated vector in ../testdata/account_derivation.json
// (contracts_vector_test.go) and the upstream fixed-input tests.

package derive

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// accountDerivationVector is packages/contracts/test/vectors/account_derivation.json, written by
// the Foundry test AccountDerivationVectorTest from the StawiAccount creation code the factory
// actually deploys (linked against its libraries, read back from the factory's CodeStore).
type accountDerivationVector struct {
	Factory          string `json:"factory"`
	AccountVersion   uint32 `json:"account_version"`
	IdentitySaltHash string `json:"identity_salt_hash"`
	AccountSalt      string `json:"account_salt"`
	CreationCode     string `json:"creation_code"`
	CreationCodeHash string `json:"creation_code_hash"`
	InitCodeHash     string `json:"init_code_hash"`
	Account          string `json:"account"`
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.TrimPrefix(strings.ToLower(s), "0x"))
	if err != nil {
		t.Fatalf("hex %q: %v", s, err)
	}
	return b
}

// TestContractsAccountDerivationVector asserts that Go derives, from the published creation
// code, exactly the address the Solidity factory deployed (P18).
func TestContractsAccountDerivationVector(t *testing.T) {
	raw, err := os.ReadFile("../testdata/account_derivation.json")
	if err != nil {
		t.Fatalf("contracts vector missing: %v", err)
	}
	var v accountDerivationVector
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	var ish [32]byte
	copy(ish[:], mustHex(t, v.IdentitySaltHash))
	var factory [20]byte
	copy(factory[:], mustHex(t, v.Factory))
	code := mustHex(t, v.CreationCode)

	if got := keccak(code); hex.EncodeToString(got[:]) != hex.EncodeToString(mustHex(t, v.CreationCodeHash)) {
		t.Fatalf("creation code hash %x", got)
	}
	salt, err := AccountSalt(FamilyEVM, v.AccountVersion, ish)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(salt[:]) != hex.EncodeToString(mustHex(t, v.AccountSalt)) {
		t.Fatalf("account salt %x", salt)
	}
	ich := InitCodeHash(code, ish)
	if hex.EncodeToString(ich[:]) != hex.EncodeToString(mustHex(t, v.InitCodeHash)) {
		t.Fatalf("init code hash %x", ich)
	}
	addr, err := AccountAddress(FamilyEVM, v.AccountVersion, ish, factory, ich)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(addr[:]) != hex.EncodeToString(mustHex(t, v.Account)) {
		t.Fatalf("address %x, want %s", addr, v.Account)
	}
}
