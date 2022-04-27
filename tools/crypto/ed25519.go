package crypto

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/dipdup-net/go-lib/tools/encoding"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

// Ed25519 -
type Ed25519 EmptyCurve

// NewEd25519 -
func NewEd25519() Ed25519 {
	return Ed25519{
		seedKey:          []byte{101, 100, 50, 53, 53, 49, 57, 32, 115, 101, 101, 100},
		addressPrefix:    []byte(encoding.PrefixPublicKeyTZ1),
		publicKeyPrefix:  []byte(encoding.PrefixED25519PublicKey),
		privateKeyPrefix: []byte(encoding.PrefixED25519SecretKey),
		signaturePrefix:  []byte(encoding.PrefixED25519Signature),
	}
}

// GeneratePrivateKey -
func (curve Ed25519) GeneratePrivateKey() ([]byte, []byte, error) {
	return ed25519.GenerateKey(rand.Reader)
}

// GetPublicKey -
func (curve Ed25519) GetPublicKey(privateKey []byte) ([]byte, error) {
	var sk ed25519.PrivateKey
	switch len(privateKey) {
	case ed25519.PrivateKeySize:
		sk = ed25519.PrivateKey(privateKey)
	case ed25519.PublicKeySize:
		sk = ed25519.NewKeyFromSeed(privateKey)
	default:
		return nil, errors.Errorf("invalid private key length: %d", len(privateKey))
	}

	pk := sk.Public()
	return []byte(pk.(ed25519.PublicKey)), nil
}

// Sign -
func (curve Ed25519) Sign(data []byte, privateKey []byte) (Signature, error) {
	if ed25519.PrivateKeySize != len(privateKey) {
		return Signature{}, errors.Errorf("invalid private key length: %d != %d", len(privateKey), ed25519.PrivateKeySize)
	}

	digest := blake2b.Sum256(data)
	sign := ed25519.Sign(privateKey, digest[:])
	return NewSignature(sign, curve.signaturePrefix), nil
}

// Verify -
func (curve Ed25519) Verify(data []byte, signature []byte, pubKey []byte) bool {
	if ed25519.PublicKeySize != len(pubKey) || ed25519.SignatureSize != len(signature) {
		return false
	}
	digest := blake2b.Sum256(data)
	return ed25519.Verify(pubKey, digest[:], signature)
}

// AddressPrefix -
func (curve Ed25519) AddressPrefix() []byte {
	return curve.addressPrefix
}

// PublicKeyPrefix -
func (curve Ed25519) PublicKeyPrefix() []byte {
	return curve.publicKeyPrefix
}

// Kind -
func (curve Ed25519) Kind() ECKind {
	return KindEd25519
}
