package crypto

import (
	"encoding/hex"

	"github.com/dipdup-net/go-lib/tools/encoding"
	"github.com/pkg/errors"
)

// Key -
type Key struct {
	pubKey  PubKey
	address string
	curve   Curve
	bytes   []byte
}

// NewKey -
func NewKey(curveKind ECKind) (Key, error) {
	curve, err := NewCurve(curveKind)
	if err != nil {
		return Key{}, err
	}

	pk, sk, err := curve.GeneratePrivateKey()
	if err != nil {
		return Key{}, err
	}

	pubKey, err := NewPubKey(pk, curveKind)
	if err != nil {
		return Key{}, err
	}

	return Key{
		curve:  curve,
		bytes:  sk,
		pubKey: pubKey,
	}, nil
}

// NewKeyFromBytes -
func NewKeyFromBytes(bytes []byte, curveKind ECKind) (Key, error) {
	curve, err := NewCurve(curveKind)
	if err != nil {
		return Key{}, err
	}

	pk, err := curve.GetPublicKey(bytes)
	if err != nil {
		return Key{}, err
	}

	pubKey, err := NewPubKey(pk, curveKind)
	if err != nil {
		return Key{}, err
	}

	return Key{
		curve:  curve,
		bytes:  bytes,
		pubKey: pubKey,
	}, nil
}

// NewKeyFromBase58 -
func NewKeyFromBase58(data string) (Key, error) {
	if len(data) < 4 {
		return Key{}, errors.Errorf("invalid public key string: %s", data)
	}

	curve, err := NewCurveFromPrefix(data[:4])
	if err != nil {
		return Key{}, err
	}

	sk, err := encoding.DecodeBase58(data)
	if err != nil {
		return Key{}, err
	}

	pk, err := curve.GetPublicKey(sk)
	if err != nil {
		return Key{}, err
	}

	pubKey, err := NewPubKey(pk, curve.Kind())
	if err != nil {
		return Key{}, err
	}

	return Key{
		curve:  curve,
		bytes:  sk,
		pubKey: pubKey,
	}, nil
}

// Bytes -
func (key Key) Bytes() []byte {
	return key.bytes
}

// Hex -
func (key Key) Hex() string {
	return hex.EncodeToString(key.bytes)
}

// String -
func (key Key) String() string {
	return key.Hex()
}

// Sign -
func (key Key) Sign(data []byte) (Signature, error) {
	return key.curve.Sign(data, append(key.bytes, key.pubKey.bytes...))
}

// Verify -
func (key Key) Verify(data, signature []byte) bool {
	return key.pubKey.Verify(data, signature)
}
