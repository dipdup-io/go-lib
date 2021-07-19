package ast

import (
	"github.com/dipdup-net/go-lib/tools/consts"
)

// Parameter -
type Parameter struct {
	*SectionType
}

// NewParameter -
func NewParameter(depth int) *Parameter {
	return &Parameter{
		SectionType: NewSectionType(consts.PARAMETER, depth),
	}
}
