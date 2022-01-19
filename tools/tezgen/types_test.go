package tezgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "test 1",
			data: []byte{0x31, 0x32, 0x33, 0x34},
			want: []byte{0x12, 0x34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Bytes
			if err := json.Unmarshal(tt.data, &b); (err != nil) != tt.wantErr {
				t.Errorf("Bytes.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, []byte(b))
		})
	}
}
