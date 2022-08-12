package translator

import (
	"io/ioutil"
	"regexp"

	"github.com/yhirose/go-peg"
)

// Converter -
type Converter struct {
	parser  *peg.Parser
	grammar string
	err     error
}

// NewConverter -
func NewConverter(opts ...ConverterOption) (Converter, error) {
	c := Converter{}
	for i := range opts {
		opts[i](&c)
	}
	if c.err != nil {
		return c, c.err
	}

	if c.grammar == "" {
		c.grammar = defaultGrammar
	}

	parser, err := peg.NewParser(c.grammar)
	if err != nil {
		return c, err
	}

	if err := parser.EnableAst(); err != nil {
		return c, err
	}
	c.parser = parser
	return c, nil
}

// FromFile -
func (c Converter) FromFile(filename string) (string, error) {
	michelson, err := readFileToString(filename)
	if err != nil {
		return "", err
	}

	return c.FromString(michelson)
}

// FromString -
func (c Converter) FromString(input string) (string, error) {
	input = removeComments(input)

	ast, err := c.parser.ParseAndGetAst(input, nil)
	if err != nil {
		return "", err
	}

	return NewJSONTranslator().Translate(ast)
}

func readFileToString(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var oneLineComment = regexp.MustCompile("#[^\n]*")
var multiLineComment = regexp.MustCompile(`\/\*[^\*]*\*/`)

func removeComments(data string) string {
	data = oneLineComment.ReplaceAllString(data, "")
	return multiLineComment.ReplaceAllString(data, "")
}
