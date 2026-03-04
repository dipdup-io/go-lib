package interfaces

import "github.com/dipdup-io/go-lib/tools/consts"

// ViewNat -
type ViewNat struct{}

// GetName -
func (f *ViewNat) GetName() string {
	return consts.ViewNatTag
}

// GetContractInterface -
func (f *ViewNat) GetContractInterface() string {
	return `{
		"entrypoints": {
			"default": {
				"prim": "nat"
			}
		},
		"is_root": true
	}`
}
