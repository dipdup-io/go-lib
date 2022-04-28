package crypto

import (
	"encoding/hex"

	"github.com/dipdup-net/go-lib/tools/encoding"
)

// Signature -
type Signature struct {
	bytes  []byte
	prefix []byte
}

// NewSignature -
func NewSignature(bytes []byte, prefix []byte) Signature {
	return Signature{bytes, prefix}
}

// Bytes -
func (s Signature) Bytes() []byte {
	return s.bytes
}

// Base58 -
func (s Signature) Base58() (string, error) {
	return encoding.EncodeBase58(s.bytes, s.prefix)
}

// Hex -
func (s Signature) Hex() string {
	return hex.EncodeToString(s.bytes)
}

// String -
func (s Signature) String() string {
	return s.Hex()
}
