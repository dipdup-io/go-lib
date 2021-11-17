package forge

import (
	"math/big"
	"reflect"
	"testing"
)

func TestForgeInt(t *testing.T) {
	tests := []struct {
		name    string
		data    *big.Int
		want    []byte
		wantErr bool
	}{
		{
			name: "Small int",
			data: big.NewInt(6),
			want: []byte{0x06},
		},
		{
			name: "Negative small int",
			data: big.NewInt(-6),
			want: []byte{0x46},
		},
		{
			name: "Medium int",
			data: big.NewInt(900),
			want: []byte{0x84, 0x0e},
		},
		{
			name: "Negative medium int",
			data: big.NewInt(-900),
			want: []byte{0xc4, 0x0e},
		},
		{
			name: "Large int",
			data: big.NewInt(917431994),
			want: []byte{0xba, 0x9a, 0xf7, 0xea, 0x06},
		},
		{
			name: "Negative large int",
			data: big.NewInt(-610913435200),
			want: []byte{0xc0, 0xf9, 0xb9, 0xd4, 0xc7, 0x23},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := NewInt()
			val.IntValue.Set(tt.data)
			got, err := ForgeInt(val.IntValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForgeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForgeInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
