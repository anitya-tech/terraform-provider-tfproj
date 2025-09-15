package filequery

import (
	"fmt"
	"path/filepath"
	"sync"
)

type QuerySource interface {
	Query(pattern *string) any
}

type Querier interface {
	Parse(filename string) (QuerySource, error)
}

var querierRegistry = map[string]Querier{
	"json": &JSONQuerier{},
	"yaml": &YAMLQuerier{},
	"yml":  &YAMLQuerier{},
}

var querierCache = sync.Map{}

func Parse(filename string) (QuerySource, error) {
	qsrc, ok := querierCache.Load(filename)
	if ok && qsrc != nil {
		return qsrc.(QuerySource), nil //nolint:forcetypeassert
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		return nil, fmt.Errorf("no extension found")
	}
	ext = ext[1:]

	querier, ok := querierRegistry[ext]
	if !ok {
		return nil, fmt.Errorf("no querier found for extension %s", ext)
	}

	nqsrc, err := querier.Parse(filename)
	if err != nil {
		return nil, err
	}

	querierCache.Store(filename, nqsrc)

	return nqsrc, nil
}
