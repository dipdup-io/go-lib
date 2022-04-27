package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyFromBase58(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		address string
		pub     string
		wantErr bool
	}{
		{
			name:    "test 1",
			data:    "edsk4GGs6oeqasc61QtgmWuQb6Yhkpx5MSva8Euq7bvzVEF3VpHdZR",
			pub:     "edpkv4B2gbqAFfcqcBTaa7DAt9w1senrDECNx9itDyBiXo9pQYDCtg",
			address: "tz1MZZApycqLQdCiCUqkFT498vqUCX3zQ4fX",
		}, {
			name:    "test 2",
			data:    "edsk3LNo1oq8bo2SQtsfinv1fcR6828pvianYKDsNr31QNta9KFH9q",
			pub:     "edpkvZNKsgFb7D7HLxnJ68cUgqEsZ47Qw81WGMndQLDvvziqcn9nVQ",
			address: "tz1TBFXHmJGZh7XpdMDCrPkxFKYo3ffNrYaq",
		}, {
			name:    "test 3",
			data:    "edsk2zpXnyz3yoFQpVekcZgbgnXbHfSrheLRkxMLNkfVjaCDQaViRa",
			pub:     "edpkvNbNZBn9PgDP6FYXzDe2fECAQLqaRPkT3SfjyJgHgDZ5ExPw4K",
			address: "tz1gR5pEVRysV4j7391xCXdbTGtQwxdGhWUY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKeyFromBase58(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKeyFromBase58() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			address, err := got.pubKey.Address()
			if err != nil {
				t.Errorf("got.pubKey.Address() error = %v", err)
				return
			}
			assert.Equal(t, tt.address, address)

			pub, err := got.pubKey.Base58()
			if err != nil {
				t.Errorf("got.pubKey.Base58() error = %v", err)
				return
			}
			assert.Equal(t, tt.pub, pub)
		})
	}
}
