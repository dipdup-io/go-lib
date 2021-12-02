package forge

import (
	"encoding/hex"
	"testing"

	"github.com/dipdup-net/go-lib/node"
	"github.com/stretchr/testify/assert"
)

func TestEndorsement(t *testing.T) {
	tests := []struct {
		name        string
		endorsement node.Endorsement
		branch      string
		want        string
		wantErr     bool
	}{
		{
			name: "test 1",
			endorsement: node.Endorsement{
				Level:    751292,
				Metadata: &node.EndorsementMetadata{},
			},
			branch: "BMbpxQAU7Jat7g9ZnKrP3brgqFX6r2VX8PPXCxNbFZeURA6DbEF",
			want:   "f8bc58c3ceaa7aaaa09d2892d0ee234231ffe46b484e5f7e7b32b5bfd618b67200000b76bc",
		}, {
			name: "test 2",
			endorsement: node.Endorsement{
				Level:    751179,
				Metadata: &node.EndorsementMetadata{},
			},
			branch: "BM2JkusQmT885mqjKiJXfMJrgQXZTwoEsCM8tkuvGyJLPrSw2ih",
			want:   "ac9f7b86a813cf29a18a16bda49a434c92a5a583f1630f0ad2b8224b7b26f05a00000b764b",
		}, {
			name: "test 3",
			endorsement: node.Endorsement{
				Level:    751447,
				Metadata: &node.EndorsementMetadata{},
			},
			branch: "BLp1dxsyPLc58x4cSMKGVevdQfgo9VBHy46kqnhsJSrNgteDPex",
			want:   "90b48e1e6ff05a6a4bd2527dbb9853a5e152906847bc4c3eda57bd1e5742a39900000b7757",
		}, {
			name: "test 4",
			endorsement: node.Endorsement{
				Level:    1479809,
				Metadata: &node.EndorsementMetadata{},
			},
			branch: "BL38RNz32eAVhgvV5bUMWUxGMq2v2wkst9UaN7CbW7hLC6THCQ6",
			want:   "2acaf73c5f06c812083f5989ea19d1b2a4d71335f988233c5145814d695b3b4c0000169481",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Endorsement(tt.endorsement, tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endorsement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, hex.EncodeToString(got))
		})
	}
}
