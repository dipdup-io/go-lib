package crypto

import "golang.org/x/crypto/blake2b"

// Blake2b160 -
func Blake2b160(data []byte) ([]byte, error) {
	hash, err := blake2b.New(20, nil)
	if err != nil {
		return nil, err
	}

	if _, err := hash.Write(data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
