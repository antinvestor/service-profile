package config

type DEK struct {
	KeyID     string
	Key       []byte
	OldKeyID  string
	OldKey    []byte
	LookUpKey []byte
}
