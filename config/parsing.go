package config

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/text/transform"
)

func expandVariables(data []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewBuffer(data), newExpandTransformer()))
}

// expandTransformer implements transform.Transformer
type expandTransformer struct {
	transform.NopResetter
}

func newExpandTransformer() *expandTransformer {
	return &expandTransformer{}
}

// Transform -
func (t *expandTransformer) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	var buf bytes.Buffer
	var index int

	startIndex := bytes.Index(src, []byte{'$', '{'})
	for startIndex != -1 {
		if _, err := buf.Write(src[index : startIndex+index]); err != nil {
			return 0, 0, err
		}
		var name, def string

		endIndex := bytes.Index(src[startIndex+index:], []byte{'}'})
		separatorIndex := bytes.Index(src[startIndex+index:startIndex+index+endIndex], []byte{':', '-'})
		if separatorIndex == -1 {
			name = string(src[startIndex+index+2 : startIndex+index+endIndex])
		} else {
			name = string(src[startIndex+index+2 : startIndex+index+separatorIndex])
			def = string(src[startIndex+index+separatorIndex+2 : startIndex+index+endIndex])
		}

		if envVal, ok := os.LookupEnv(name); ok {
			def = envVal
		}

		if def == "" {
			def = `""`
		}

		if _, err := buf.Write([]byte(def)); err != nil {
			return 0, 0, err
		}

		index += startIndex + endIndex + 1
		startIndex = bytes.Index(src[index:], []byte{'$', '{'})
	}

	if _, err := buf.Write(src[index:]); err != nil {
		return 0, 0, err
	}

	return copy(dst, buf.Bytes()), len(src), nil
}
