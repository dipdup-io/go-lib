package crypto

import (
	"encoding/hex"

	"github.com/dipdup-net/go-lib/tools/encoding"
	"github.com/pkg/errors"
)

// PubKey -
type PubKey struct {
	address string
	curve   Curve
	bytes   []byte
}

// NewPubKey -
func NewPubKey(bytes []byte, curveKind ECKind) (PubKey, error) {
	curve, err := NewCurve(curveKind)
	if err != nil {
		return PubKey{}, err
	}
	return PubKey{
		bytes: bytes,
		curve: curve,
	}, nil
}

// NewPubKeyFromBase58 -
func NewPubKeyFromBase58(data string) (PubKey, error) {
	if len(data) < 4 {
		return PubKey{}, errors.Errorf("invalid public key string: %s", data)
	}
	curve, err := NewCurveFromPrefix(data[:4])
	if err != nil {
		return PubKey{}, err
	}

	bytes, err := encoding.DecodeBase58(data)
	if err != nil {
		return PubKey{}, err
	}
	return PubKey{
		curve: curve,
		bytes: bytes,
	}, nil
}

// Address -
func (pk PubKey) Address() (string, error) {
	if pk.address != "" {
		return pk.address, nil
	}

	data, err := Blake2b160(pk.bytes)
	if err != nil {
		return "", err
	}

	result, err := encoding.EncodeBase58(data, pk.curve.AddressPrefix())
	if err != nil {
		return "", err
	}
	pk.address = result
	return pk.address, nil
}

// Bytes -
func (pk PubKey) Bytes() []byte {
	return pk.bytes
}

// Base58 -
func (pk PubKey) Base58() (string, error) {
	return encoding.EncodeBase58(pk.bytes, pk.curve.PublicKeyPrefix())
}

// Hex -
func (pk PubKey) Hex() string {
	return hex.EncodeToString(pk.bytes)
}

// String -
func (pk PubKey) String() string {
	return pk.Hex()
}

// Verify -
func (pk PubKey) Verify(data, signature []byte) bool {
	return pk.curve.Verify(data, signature, pk.bytes)
}
