package crypto

import "github.com/pkg/errors"

// ECKind -
type ECKind int

const (
	KindEd25519 ECKind = iota + 1
	KindSecp256k1
	KindNistP256
)

// Curve -
type Curve interface {
	GeneratePrivateKey() ([]byte, []byte, error)
	GetPublicKey(privateKey []byte) ([]byte, error)
	Sign(data []byte, privateKey []byte) (Signature, error)
	Verify(data []byte, signature []byte, pubKey []byte) bool
	AddressPrefix() []byte
	PublicKeyPrefix() []byte
	Kind() ECKind
}

// EmptyCurve -
type EmptyCurve struct {
	addressPrefix    []byte
	publicKeyPrefix  []byte
	privateKeyPrefix []byte
	signaturePrefix  []byte
	seedKey          []byte
}

// NewCurve -
func NewCurve(kind ECKind) (Curve, error) {
	switch kind {
	case KindEd25519:
		return NewEd25519(), nil
	// case KindSecp256k1:
	// case KindNistP256:
	default:
		return nil, errors.Errorf("unknown curve kind: %d", kind)
	}
}

// NewCurveFromPrefix -
func NewCurveFromPrefix(prefix string) (Curve, error) {
	switch prefix {
	case "edsig", "edsk", "edpk", "tz1", "edesk":
		return NewEd25519(), nil
	// case "spsig", "sppk", "spsk", "tz2", "spesk":
	// 	// Secp256k1
	// case "p2sig", "p2pk", "p2sk", "tz3", "p2esk":
	// 	// NistP256
	default:
		return nil, errors.Errorf("unknown curve prefix: %s", prefix)
	}
}
