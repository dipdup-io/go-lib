package forge

import (
	"bytes"
	"encoding/binary"

	"github.com/dipdup-net/go-lib/node"
	"github.com/dipdup-net/go-lib/tools/base"
	"github.com/dipdup-net/go-lib/tools/types"
	"github.com/pkg/errors"
)

// Transaction -
func Transaction(transaction node.Transaction) ([]byte, error) {
	buf := new(bytes.Buffer)

	tag, ok := operationTags[node.KindTransaction]
	if !ok {
		return nil, errors.Errorf("unknown operation tag: %s", node.KindTransaction)
	}

	kind, err := ForgeNat(types.NewBigIntFromString(tag))
	if err != nil {
		return nil, errors.Wrap(err, "kind forging")
	}
	buf.Write(kind)

	source, err := Address(transaction.Source, true)
	if err != nil {
		return nil, errors.Wrap(err, "source forging")
	}
	buf.Write(source)

	fee, err := ForgeNat(types.NewBigIntFromString(transaction.Fee))
	if err != nil {
		return nil, errors.Wrap(err, "fee forging")
	}
	buf.Write(fee)

	counter, err := ForgeNat(types.NewBigIntFromString(transaction.Counter))
	if err != nil {
		return nil, errors.Wrap(err, "counter forging")
	}
	buf.Write(counter)

	gasLimit, err := ForgeNat(types.NewBigIntFromString(transaction.GasLimit))
	if err != nil {
		return nil, errors.Wrap(err, "gas limit forging")
	}
	buf.Write(gasLimit)

	storageLimit, err := ForgeNat(types.NewBigIntFromString(transaction.StorageLimit))
	if err != nil {
		return nil, errors.Wrap(err, "storage limit forging")
	}
	buf.Write(storageLimit)

	amount, err := ForgeNat(types.NewBigIntFromString(transaction.Amount))
	if err != nil {
		return nil, errors.Wrap(err, "amount forging")
	}
	buf.Write(amount)

	destination, err := Address(transaction.Destination, false)
	if err != nil {
		return nil, errors.Wrap(err, "destination forging")
	}
	buf.Write(destination)

	hasParams := transaction.Parameters != nil
	buf.Write(ForgeBool(hasParams))
	if hasParams {
		buf.Write(forgeEntrypoint(transaction.Parameters.Entrypoint))

		var node base.Node
		if err := json.Unmarshal(*transaction.Parameters.Value, &node); err != nil {
			return nil, errors.Wrap(err, "parameter's value unmarshaling")
		}
		forger := NewMichelson()
		forger.Nodes = []*base.Node{&node}
		value, err := forger.Forge()
		if err != nil {
			return nil, errors.Wrap(err, "parameter's value forging")
		}
		buf.Write(ForgeArray(value, 4))
	}

	return buf.Bytes(), nil
}

// ForgeArray -
func ForgeArray(value []byte, l uint64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, uint64(len(value)))

	bytes := reverse(buf.Bytes()[0:l])
	return append(bytes, value...)
}

func forgeEntrypoint(value string) []byte {
	buf := new(bytes.Buffer)

	if val, ok := entrypointTags[value]; ok {
		buf.WriteByte(val)
	} else {
		buf.WriteByte(byte(255))
		buf.Write(ForgeArray(bytes.NewBufferString(value).Bytes(), 1))
	}

	return buf.Bytes()
}
