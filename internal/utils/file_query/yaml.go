package filequery

import (
	"os"

	"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
)

type YAMLQuerier struct{}

func (q *YAMLQuerier) Parse(filename string) (QuerySource, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	json, err := yaml.YAMLToJSON(content)
	if err != nil {
		return nil, err
	}

	return &JSONQuerier{data: gjson.Parse(string(json))}, nil
}
