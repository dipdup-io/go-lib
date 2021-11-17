package forge

import (
	"bytes"
	"math/big"

	"github.com/dipdup-net/go-lib/tools/base"
	"github.com/dipdup-net/go-lib/tools/types"
)

// Int -
type Int base.Node

// NewInt -
func NewInt() *Int {
	return &Int{
		IntValue: types.NewBigInt(0),
	}
}

// Unforge -
func (val *Int) Unforge(data []byte) (int, error) {
	buffer := new(bytes.Buffer)
	for i := range data {
		buffer.WriteByte(data[i])
		if data[i] < 128 {
			break
		}
	}

	parts := buffer.Bytes()
	for i := len(parts) - 1; i > 0; i-- {
		num := int64(parts[i] & 0x7f)
		val.IntValue.Int = val.IntValue.Lsh(val.IntValue.Int, 7)
		val.IntValue.Int = val.IntValue.Or(val.IntValue.Int, big.NewInt(num))
	}

	if len(parts) > 0 {
		num := int64(parts[0] & 0x3f)
		val.IntValue.Int = val.IntValue.Lsh(val.IntValue.Int, 6)
		val.IntValue.Int = val.IntValue.Or(val.IntValue.Int, big.NewInt(num))

		if parts[0]&0x40 > 0 {
			val.IntValue.Int = val.IntValue.Neg(val.IntValue.Int)
		}
	}

	return buffer.Len(), nil
}

// Forge -
func (val *Int) Forge() ([]byte, error) {
	data, err := ForgeInt(val.IntValue)
	if err != nil {
		return nil, err
	}
	return append([]byte{ByteInt}, data...), nil
}
