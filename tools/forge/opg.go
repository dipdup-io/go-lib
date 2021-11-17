package forge

import (
	"bytes"

	"github.com/dipdup-net/go-lib/node"
	"github.com/dipdup-net/go-lib/tools/encoding"
	"github.com/pkg/errors"
)

// OPG -
func OPG(branch string, operations ...node.Operation) ([]byte, error) {
	result := new(bytes.Buffer)
	if branch != "" {
		decoded, err := encoding.DecodeBase58(branch)
		if err != nil {
			return nil, errors.Wrap(err, "failed to forge operation")
		}
		result.Write(decoded)
	}

	for i := range operations {
		switch typ := operations[i].Body.(type) {

		// TODO: realize others operation kinds
		case node.Transaction:
			forged, err := Transaction(typ)
			if err != nil {
				return nil, errors.Wrap(err, "failed to forge transaction")
			}
			result.Write(forged)
		default:
		}

	}

	return result.Bytes(), nil
}
