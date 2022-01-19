package tezgen

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Unit -
type Unit struct{}

// Bytes -
type Bytes []byte

// UnmarshalJSON -
func (b *Bytes) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if len(data)%2 == 1 {
		return errors.Errorf("invalid bytes value with length %d: %v", data, len(data))
	}
	if len(data) > 1 {
		if data[0] == '"' && data[len(data)-1] == '"' {
			data = data[1 : len(data)-1]
		}
	}
	byt := make([]byte, hex.DecodedLen(len(data)))
	if _, err := hex.Decode(byt, data); err != nil {
		return err
	}

	*b = make([]byte, 0)
	*b = append(*b, byt...)
	return nil
}

// MarshalJSON -
func (b Bytes) MarshalJSON() ([]byte, error) {
	str := hex.EncodeToString(b)
	return []byte(strconv.Quote(str)), nil
}

// Timestamp -
type Timestamp struct {
	val time.Time
}

// NewTimestamp -
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{t}
}

// Value -
func (t Timestamp) Value() time.Time {
	return t.val
}

// UnmarshalJSON -
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	ts, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		t.val = time.Unix(ts, 0)
		return nil
	}

	t.val, err = time.Parse(time.RFC3339, str)
	return err
}

// MarshalJSON -
func (t Timestamp) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(t.val.Unix(), 10)
	return []byte(str), nil
}

// Address -
type Address string

// Contract -
type Contract string

// SaplingTransaction -
type SaplingTransaction string

// SaplingState -
type SaplingState struct {
	State int64
	Array []interface{}
}

// UnmarshalJSON -
func (ss *SaplingState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &ss.State); err != nil {
		return err
	}

	return json.Unmarshal(data, &ss.Array)
}

// MarshalJSON -
func (ss SaplingState) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(ss.State, 10)
	return []byte(str), nil
}

// Int -
type Int struct{ *big.Int }

// NewInt -
func NewInt(val int64) Int {
	return Int{
		big.NewInt(val),
	}
}

// UnmarshalJSON -
func (i *Int) UnmarshalJSON(data []byte) error {
	if i.Int == nil {
		i.Int = big.NewInt(0)
	}
	if len(data) > 2 {
		if data[0] == '"' && data[len(data)-1] == '"' {
			data = data[1 : len(data)-1]
		}
	}
	return i.Int.UnmarshalJSON(data)
}

// MarshalJSON -
func (i Int) MarshalJSON() ([]byte, error) {
	if i.Int == nil {
		i.Int = big.NewInt(0)
	}
	return i.Int.MarshalJSON()
}
