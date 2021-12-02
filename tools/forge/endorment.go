package forge

import (
	"bytes"
	"encoding/binary"

	"github.com/dipdup-net/go-lib/node"
	"github.com/dipdup-net/go-lib/tools/types"
)

// Endorsement -
func Endorsement(endorsement node.Endorsement, branch string) ([]byte, error) {
	var buf bytes.Buffer

	branchForged, err := ForgeString(branch)
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(branchForged); err != nil {
		return nil, err
	}

	tag, err := ForgeNat(types.NewBigInt(0))
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(tag); err != nil {
		return nil, err
	}

	level := make([]byte, 4)
	binary.BigEndian.PutUint32(level, uint32(endorsement.Level))
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(level); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
