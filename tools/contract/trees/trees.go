package trees

import "github.com/dipdup-net/go-lib/tools/ast"

var (
	Nat, _           = ast.NewTypedAstFromString(`{"prim":"nat"}`)
	Address, _       = ast.NewTypedAstFromString(`{"prim":"address"}`)
	Token, _         = ast.NewTypedAstFromString(`{"prim":"pair","args":[{"prim":"address"},{"prim":"nat"}]}`)
	String, _        = ast.NewTypedAstFromString(`{"prim":"string"}`)
	Bytes, _         = ast.NewTypedAstFromString(`{"prim":"bytes"}`)
	TokenMetadata, _ = ast.NewTypedAstFromString(`{"prim":"pair","args":[{"prim":"nat","annots":["%token_id"]},{"prim":"map","annots":["%token_info"],"args":[{"prim":"string"},{"prim":"bytes"}]}]}`)
)
