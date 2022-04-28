package translator

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverter_FromFile(t *testing.T) {
	files, err := ioutil.ReadDir("./tests/")
	if err != nil {
		t.Errorf("ioutil.ReadDir(./tests/) error = %v", err)
		return
	}

	c, err := NewConverter()
	if err != nil {
		t.Errorf("Converter.NewConverter() error = %v", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			resultFilename := fmt.Sprintf("tests/%s/code.json", file.Name())
			resultBytes, err := ioutil.ReadFile(resultFilename)
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
				return
			}

			filename := fmt.Sprintf("tests/%s/code.tz", file.Name())
			got, err := c.FromFile(filename)
			if err != nil {
				t.Errorf("Converter.FromFile() error = %v", err)
				return
			}

			assert.JSONEq(t, string(resultBytes), got, "JSON comparing")
		})
	}
}
