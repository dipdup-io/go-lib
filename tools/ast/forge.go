package ast

import "github.com/dipdup-io/go-lib/tools/forge"

// Forge -
func Forge(node Base, optimized bool) (string, error) {
	baseAST, err := node.ToBaseNode(optimized)
	if err != nil {
		return "", err
	}
	return forge.ToString(baseAST)
}
