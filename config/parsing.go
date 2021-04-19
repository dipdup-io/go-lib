package config

import (
	"os"
	"regexp"
)

var defaultEnv = regexp.MustCompile(`\${(?P<name>[\w\.]{1,}):-(?P<value>[\w\.:/]*)}`)

func expandEnv(data string) string {
	vars := defaultEnv.FindAllStringSubmatch(data, -1)
	data = defaultEnv.ReplaceAllString(data, `${$name}`)

	for i := range vars {
		if _, ok := os.LookupEnv(vars[i][1]); !ok {
			os.Setenv(vars[i][1], vars[i][2])
		}
	}

	data = os.ExpandEnv(data)
	return data
}
