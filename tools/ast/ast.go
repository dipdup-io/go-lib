package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/dipdup-net/go-lib/tools/base"
	"github.com/dipdup-net/go-lib/tools/consts"
	"github.com/dipdup-net/go-lib/tools/types"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// TypedAst -
type TypedAst struct {
	Nodes   []Node
	settled bool
}

// UnmarshalJSON -
func (ast *TypedAst) UnmarshalJSON(data []byte) error {
	parsed, err := NewTypedAstFromBytes(data)
	if err != nil {
		return err
	}
	ast.Nodes = parsed.Nodes
	return nil
}

// NewTypedAST -
func NewTypedAST() *TypedAst {
	return &TypedAst{
		Nodes: make([]Node, 0),
	}
}

// NewTypedAstFromBytes -
func NewTypedAstFromBytes(data []byte) (*TypedAst, error) {
	var tree UntypedAST
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, err
	}
	return tree.ToTypedAST()
}

// NewTypedAstFromString -
func NewTypedAstFromString(data string) (*TypedAst, error) {
	var tree UntypedAST
	if err := json.UnmarshalFromString(data, &tree); err != nil {
		return nil, err
	}
	return tree.ToTypedAST()
}

// NewSettledTypedAst -
func NewSettledTypedAst(tree, data string) (*TypedAst, error) {
	typ, err := NewTypedAstFromString(tree)
	if err != nil {
		return nil, err
	}

	var treeData UntypedAST
	if err := json.UnmarshalFromString(data, &treeData); err != nil {
		return nil, err
	}

	err = typ.Settle(treeData)
	return typ, err
}

// IsSettled -
func (a *TypedAst) IsSettled() bool {
	return a.settled
}

// String -
func (a *TypedAst) String() string {
	var s strings.Builder
	for i := range a.Nodes {
		s.WriteString(a.Nodes[i].String())
	}
	return s.String()
}

// Settle -
func (a *TypedAst) Settle(untyped UntypedAST) error {
	if len(untyped) == len(a.Nodes) {
		for i := range untyped {
			if err := a.Nodes[i].ParseValue(untyped[i]); err != nil {
				return err
			}
		}
		a.settled = true
		return nil
	} else if len(a.Nodes) == 1 {
		if _, ok := a.Nodes[0].(*Pair); ok {
			newUntyped := &base.Node{
				Prim: consts.Pair,
				Args: untyped,
			}
			if err := a.Nodes[0].ParseValue(newUntyped); err != nil {
				return err
			}
			a.settled = true
			return nil
		}
	}
	return errors.Wrap(consts.ErrTreesAreDifferent, "TypedAst.Settle")
}

// SettleFromBytes -
func (a *TypedAst) SettleFromBytes(data []byte) error {
	var tree UntypedAST
	if err := json.Unmarshal(data, &tree); err != nil {
		return err
	}
	return a.Settle(tree)
}

// ToMiguel -
func (a *TypedAst) ToMiguel() ([]*MiguelNode, error) {
	nodes := make([]*MiguelNode, 0)
	for i := range a.Nodes {
		m, err := a.Nodes[i].ToMiguel()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, m)
	}
	return nodes, nil
}

// GetEntrypoints -
func (a *TypedAst) GetEntrypoints() []string {
	entrypoints := make([]string, 0)
	for i := range a.Nodes {
		entrypoints = append(entrypoints, a.Nodes[i].GetEntrypoints()...)
	}
	if len(entrypoints) == 1 && entrypoints[0] == "" {
		return []string{consts.DefaultEntrypoint}
	}

	for i := range entrypoints {
		if entrypoints[i] == "" {
			entrypoints[i] = fmt.Sprintf("entrypoint_%d", i)
		}
	}

	return entrypoints
}

// ToBaseNode -
func (a *TypedAst) ToBaseNode(optimized bool) (*base.Node, error) {
	if len(a.Nodes) == 1 {
		return a.Nodes[0].ToBaseNode(optimized)
	}
	return arrayToBaseNode(a.Nodes, optimized)
}

// ToJSONSchema -
func (a *TypedAst) ToJSONSchema() (*JSONSchema, error) {
	if len(a.Nodes) == 1 {
		if a.Nodes[0].GetPrim() == consts.UNIT {
			return nil, nil
		}
		return a.Nodes[0].ToJSONSchema()
	}

	s := &JSONSchema{
		Type:       JSONSchemaTypeObject,
		Properties: make(map[string]*JSONSchema),
	}

	for i := range a.Nodes {
		child, err := a.Nodes[i].ToJSONSchema()
		if err != nil {
			return nil, err
		}
		s.Properties[a.Nodes[i].GetName()] = child
	}

	return s, nil
}

// FromJSONSchema -
func (a *TypedAst) FromJSONSchema(data map[string]interface{}) error {
	for i := range a.Nodes {
		if err := a.Nodes[i].FromJSONSchema(data); err != nil {
			return err
		}
	}
	return nil
}

// FindByName -
func (a *TypedAst) FindByName(name string, isEntrypoint bool) Node {
	for i := range a.Nodes {
		if node := a.Nodes[i].FindByName(name, isEntrypoint); node != nil {
			return node
		}
	}

	if isEntrypoint && name == consts.DefaultEntrypoint {
		return a.Nodes[0]
	}

	if isEntrypoint && strings.HasPrefix(name, "entrypoint_") {
		num, err := strconv.ParseInt(strings.TrimPrefix(name, "entrypoint_"), 10, 64)
		if err != nil {
			return nil
		}
		var depth int64
		return findCustomEntrypoint(a.Nodes[0], num, &depth)
	}

	return nil
}

// ToParameters -
func (a *TypedAst) ToParameters(entrypoint string) ([]byte, error) {
	if entrypoint == "" {
		if len(a.Nodes) == 1 {
			return a.Nodes[0].ToParameters()
		}

		return buildListParameters(a.Nodes)
	}

	node := a.FindByName(entrypoint, true)
	if node != nil {
		return node.ToParameters()
	}
	return nil, nil
}

// Docs -
func (a *TypedAst) Docs(entrypoint string) ([]Typedef, error) {
	if entrypoint == DocsFull {
		if len(a.Nodes) == 1 {
			docs, typ, err := a.Nodes[0].Docs(DocsFull)
			if err != nil {
				return nil, err
			}
			if docs != nil {
				return docs, nil
			}
			return []Typedef{
				{
					Name: consts.DefaultEntrypoint,
					Type: typ,
				},
			}, nil
		}
		return buildArrayDocs(a.Nodes)
	}

	node := a.FindByName(entrypoint, true)
	if node != nil {
		docs, typName, err := node.Docs(DocsFull)
		if docs != nil {
			return docs, err
		}
		return []Typedef{
			{
				Name: entrypoint,
				Type: typName,
			},
		}, nil
	}
	return nil, nil
}

// GetEntrypointsDocs -
func (a *TypedAst) GetEntrypointsDocs() ([]EntrypointType, error) {
	docs, err := a.Docs(DocsFull)
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, nil
	}

	if docs[0].Type == consts.OR {
		response := make([]EntrypointType, 0)
		for i := range docs[0].Args {
			name := docs[0].Args[i].Key
			if strings.HasPrefix(name, "@") {
				name = fmt.Sprintf("entrypoint_%d", len(response))
			}
			entrypoint := EntrypointType{
				Name: name,
			}
			eDocs, err := a.Docs(docs[0].Args[i].Key)
			if err != nil {
				return nil, err
			}
			entrypoint.Type = eDocs
			response = append(response, entrypoint)
		}

		return response, nil
	}
	entrypoint := EntrypointType{
		Name: consts.DefaultEntrypoint,
		Type: docs,
	}
	return []EntrypointType{entrypoint}, nil
}

// Compare -
func (a *TypedAst) Compare(b *TypedAst) (int, error) {
	if len(a.Nodes) != len(b.Nodes) {
		return 0, consts.ErrTypeIsNotComparable
	}
	for i := range a.Nodes {
		res, err := a.Nodes[i].Compare(b.Nodes[i])
		if err != nil {
			return res, err
		}
		if res != 0 {
			return res, nil
		}
	}
	return 0, nil
}

// Diff -
func (a *TypedAst) Diff(b *TypedAst) (*MiguelNode, error) {
	if b == nil {
		tree, err := a.ToMiguel()
		if err != nil {
			return nil, err
		}
		if len(tree) == 0 {
			return nil, nil
		}
		for i := range tree {
			tree[i].setDiffType(MiguelKindCreate)
		}
		return tree[0], nil
	}
	if len(b.Nodes) == 1 && len(a.Nodes) == 1 {
		return a.Nodes[0].Distinguish(b.Nodes[0])
	}
	return nil, nil
}

// FromParameters - fill(settle) subtree in `a` with `data`. Returned settled subtree and error if occurred. If `entrypoint` is empty string it will try to settle main tree.
func (a *TypedAst) FromParameters(data *types.Parameters) (*TypedAst, error) {
	var tree UntypedAST
	if err := json.Unmarshal(data.Value, &tree); err != nil {
		return nil, err
	}

	if data.Entrypoint != "" {
		subTree := a.FindByName(data.Entrypoint, true)
		if subTree != nil {
			err := subTree.ParseValue(tree[0])
			return &TypedAst{
				Nodes:   []Node{subTree},
				settled: true,
			}, err
		}
	}
	err := a.Settle(tree)
	return a, err
}

// EnrichBigMap -
func (a *TypedAst) EnrichBigMap(bmd []*types.BigMapDiff) error {
	for i := range a.Nodes {
		if err := a.Nodes[i].EnrichBigMap(bmd); err != nil {
			return err
		}
	}
	return nil
}

// EqualType -
func (a *TypedAst) EqualType(b *TypedAst) bool {
	if len(a.Nodes) != len(b.Nodes) {
		return false
	}

	for i := range a.Nodes {
		if !a.Nodes[i].EqualType(b.Nodes[i]) {
			return false
		}
	}

	return true
}

// FindBigMapByPtr -
func (a *TypedAst) FindBigMapByPtr() map[int64]*BigMap {
	res := make(map[int64]*BigMap)
	for i := range a.Nodes {
		if m := a.Nodes[i].FindPointers(); m != nil {
			for p, bm := range m {
				res[p] = bm
			}
		}
	}
	return res
}

// GetJSONModel -
func (a *TypedAst) GetJSONModel(model JSONModel) {
	if model == nil {
		model = make(JSONModel)
	}

	for i := range a.Nodes {
		a.Nodes[i].GetJSONModel(model)
		if a.Nodes[i].IsPrim(consts.PAIR) {
			name := a.Nodes[i].GetName()
			if val, ok := model[name]; ok {
				if data, ok := val.(JSONModel); ok {
					delete(model, name)
					for key, value := range data {
						model[key] = value
					}
				}
			}
		}
	}
}

// Unwrap - clear parameters from Left/Right. Tree must be settled.
func (a *TypedAst) UnwrapAndGetEntrypointName() (Node, string) {
	if !a.IsSettled() || len(a.Nodes) != 1 {
		return nil, ""
	}

	return unwrap(a.Nodes[0], "0")
}

// ParametersForExecution -
func (a *TypedAst) ParametersForExecution(entrypoint string, data map[string]interface{}) (*types.Parameters, error) {
	if len(a.Nodes) == 0 {
		return nil, consts.ErrEmptyTree
	}

	if strings.HasPrefix(entrypoint, "entrypoint_") {
		if params := wrapForExecution(a.Nodes[0], data, entrypoint); params != nil {
			return params, nil
		}
	}

	if node := a.FindByName(entrypoint, true); node != nil {
		return settleForExecution(node, data, entrypoint)
	}
	return nil, consts.ErrEmptyTree
}

func settleForExecution(node Node, data map[string]interface{}, entrypoint string) (*types.Parameters, error) {
	if node.IsNamed() && node.IsPrim(consts.PAIR) {
		data = map[string]interface{}{
			node.GetName(): data,
		}
	}
	if err := node.FromJSONSchema(data); err != nil {
		return nil, err
	}
	params, err := node.ToParameters()
	if err != nil {
		return nil, err
	}
	return &types.Parameters{
		Entrypoint: entrypoint,
		Value:      params,
	}, nil
}

func wrapForExecution(node Node, data map[string]interface{}, entrypoint string) *types.Parameters {
	num, err := strconv.ParseInt(strings.TrimPrefix(entrypoint, "entrypoint_"), 10, 64)
	if err != nil {
		return nil
	}

	var depth int64
	current := findCustomEntrypoint(node, num, &depth)
	if current == nil {
		return nil
	}
	if err := current.FromJSONSchema(data); err != nil {
		return nil
	}
	params, err := node.ToParameters()
	if err != nil {
		return nil
	}
	return &types.Parameters{
		Entrypoint: consts.DefaultEntrypoint,
		Value:      params,
	}
}

func findCustomEntrypoint(node Node, num int64, depth *int64) Node {
	typ, ok := node.(*Or)
	if ok {
		if node := findCustomEntrypoint(typ.LeftType, num, depth); node != nil {
			typ.key = leftKey
			return node
		}
		if node := findCustomEntrypoint(typ.RightType, num, depth); node != nil {
			typ.key = rightKey
			return node
		}
		return nil
	}

	if num == *depth {
		return node
	}
	*depth += 1
	return nil
}

func unwrap(node Node, path string) (Node, string) {
	or, ok := node.(*Or)
	if !ok {
		name := node.GetName()
		if strings.HasPrefix(name, "@") {
			i, err := strconv.ParseInt(path, 2, 64)
			if err == nil {
				name = fmt.Sprintf("@entrypoint_%d", i)
			}
		}
		return node, name
	}

	switch or.key {
	case leftKey:
		return unwrap(or.LeftType, path+"0")
	case rightKey:
		return unwrap(or.RightType, path+"1")
	}
	return node, node.GetName()
}

func marshalJSON(prim string, annots []string, args ...Node) ([]byte, error) {
	var builder bytes.Buffer
	builder.WriteByte('{')
	builder.WriteString(fmt.Sprintf(`"prim": "%s"`, prim))
	if len(args) > 0 {
		builder.WriteString(`, "args": [`)
		for i := range args {
			typ, err := json.Marshal(args[i])
			if err != nil {
				return nil, err
			}
			if _, err := builder.Write(typ); err != nil {
				return nil, err
			}
			if i < len(args)-1 {
				builder.WriteByte(',')
			}
		}
		builder.WriteByte(']')
	}
	if len(annots) > 0 {
		builder.WriteString(fmt.Sprintf(`, "annots": ["%s"]`, strings.Join(annots, `","`)))

	}
	builder.WriteByte('}')
	return builder.Bytes(), nil
}
