package filequery

import (
	"errors"
	"os"

	"github.com/tidwall/gjson"
)

type JSONQuerier struct {
	data gjson.Result
}

func (q *JSONQuerier) Parse(filename string) (QuerySource, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	json := string(content)

	if !gjson.Valid(json) {
		return nil, errors.New("invalid json")
	}

	return &JSONQuerier{data: gjson.Parse(json)}, nil
}

func (q *JSONQuerier) Query(pattern *string) any {
	if pattern == nil {
		return q.data.Value()
	}
	return q.data.Get(*pattern).Value()
}
