package forge

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/dipdup-net/go-lib/tools/types"
	"github.com/pkg/errors"
)

// ForgeNat -
func ForgeNat(value *types.BigInt) ([]byte, error) {
	if value == nil || value.Sign() == -1 {
		return nil, errors.Errorf("invalid nat value: %v", value)
	}

	buf := new(bytes.Buffer)
	val := value.Int64()

	var end bool
	for !end {
		b := byte(val & 0x7f)
		val >>= 7
		end = val <= 0
		if !end {
			b |= 0x80
		}

		buf.WriteByte(b)
	}

	return buf.Bytes(), nil
}

// ForgeInt -
func ForgeInt(value *types.BigInt) ([]byte, error) {
	if value == nil {
		return nil, errors.New("Invalid int value")
	}

	isNegative := value.Sign() == -1
	bits := value.Text(2)
	if isNegative {
		bits = bits[1:]
	}
	bitsCount := len(bits)

	var pad int
	switch {
	case (bitsCount-6)%7 == 0:
		pad = bitsCount
	case bitsCount > 6:
		pad = bitsCount + 7 - (bitsCount-6)%7
	default:
		pad = 6
	}
	bits = fmt.Sprintf("%0*s", pad, bits)

	segments := make([]string, 0)
	for i := 0; i <= pad/7; i++ {
		idx := 7 * i
		length := int(math.Min(7, float64(pad-7*i)))
		segments = append(segments, bits[idx:(idx+length)])
	}

	segments = reverse(segments)
	if isNegative {
		segments[0] = fmt.Sprintf("1%s", segments[0])
	} else {
		segments[0] = fmt.Sprintf("0%s", segments[0])
	}

	data := make([]byte, 0)

	for i := 0; i < len(segments); i++ {
		prefix := "1"
		if i == len(segments)-1 {
			prefix = "0"
		}
		val, err := strconv.ParseUint(prefix+segments[i], 2, 8)
		if err != nil {
			return nil, err
		}
		data = append(data, byte(val))
	}

	return data, nil
}

// ForgeBool -
func ForgeBool(value bool) []byte {
	if value {
		return []byte{255}
	}
	return []byte{0}
}

// ForgeString -
func ForgeString(value string) ([]byte, error) {
	return nil, nil
}

func reverse(arr []string) []string {
	for i := len(arr)/2 - 1; i >= 0; i-- {
		opp := len(arr) - 1 - i
		arr[i], arr[opp] = arr[opp], arr[i]
	}
	return arr
}

func reverseBytes(arr []byte) []byte {
	for i := len(arr)/2 - 1; i >= 0; i-- {
		opp := len(arr) - 1 - i
		arr[i], arr[opp] = arr[opp], arr[i]
	}
	return arr
}
