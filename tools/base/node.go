package base

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/dipdup-net/go-lib/tools/consts"
	"github.com/dipdup-net/go-lib/tools/types"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Node - struct for parsing micheline
type Node struct {
	Prim        string        `json:"prim,omitempty"`
	Args        []*Node       `json:"args,omitempty"`
	Annots      []string      `json:"annots,omitempty"`
	StringValue *string       `json:"string,omitempty"`
	BytesValue  *string       `json:"bytes,omitempty"`
	IntValue    *types.BigInt `json:"int,omitempty"`
}

// UnmarshalJSON -
func (node *Node) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return consts.ErrInvalidJSON
	}
	if data[0] == '[' {
		node.Prim = consts.PrimArray
		node.Args = make([]*Node, 0)
		return json.Unmarshal(data, &node.Args)
	} else if data[0] == '{' {
		type buf Node
		return json.Unmarshal(data, (*buf)(node))
	}
	return consts.ErrInvalidJSON
}

// MarshalJSON -
func (node *Node) MarshalJSON() ([]byte, error) {
	if node.Prim == consts.PrimArray {
		return json.Marshal(node.Args)
	}

	type buf Node
	return json.Marshal((*buf)(node))
}

// GetAnnotations - returns all node`s annotations recursively
func (node *Node) GetAnnotations() map[string]struct{} {
	annots := make(map[string]struct{})
	for i := range node.Annots {
		if len(node.Annots[i]) == 0 {
			continue
		}
		if node.Annots[i][0] == consts.AnnotPrefixFieldName || node.Annots[i][0] == consts.AnnotPrefixrefixTypeName {
			annots[node.Annots[i][1:]] = struct{}{}
		}
	}
	for i := range node.Args {
		for k := range node.Args[i].GetAnnotations() {
			annots[k] = struct{}{}
		}
	}
	return annots
}

// Compare -
func (node *Node) Compare(second *Node) bool {
	if node.Prim != second.Prim {
		return false
	}
	if len(node.Args) != len(second.Args) {
		return false
	}
	for i := range node.Args {
		if !node.Args[i].Compare(second.Args[i]) {
			return false
		}
	}
	return true
}

// String - converts node info to string
func (node *Node) String() string {
	return node.print(0) + "\n"
}

func (node *Node) print(depth int) string {
	var s strings.Builder
	s.WriteByte('\n')
	s.WriteString(strings.Repeat(consts.DefaultIndent, depth))
	switch {
	case node.Prim != "":
		s.WriteString(node.Prim)
		for i := range node.Args {
			s.WriteString(node.Args[i].print(depth + 1))
		}
	case node.IntValue != nil:
		s.WriteString(fmt.Sprintf("Int=%d", *node.IntValue))
	case node.BytesValue != nil:
		s.WriteString(fmt.Sprintf("Bytes=%s", *node.BytesValue))
	case node.StringValue != nil:
		s.WriteString(fmt.Sprintf("String=%s", *node.StringValue))
	}
	return s.String()
}

var lambdaPackedReg = regexp.MustCompile("^0502[0-9a-f]{8}0[3-9]")

// IsLambda -
func (node *Node) IsLambda() bool {
	if node.BytesValue == nil {
		return false
	}
	input := *node.BytesValue
	if len(input) < 24 {
		return false
	}
	if !lambdaPackedReg.MatchString(input) {
		return false
	}
	b, err := hex.DecodeString(input[22:24])
	if err != nil {
		return false
	}
	if len(b) != 1 {
		return false
	}
	if 0x0c > b[0] || 0x75 < b[0] {
		return false
	}
	return 0x58 >= b[0] || 0x6f <= b[0]

}
