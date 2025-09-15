package tf_project

import (
	"os"
	"sync"
)

var once sync.Once
var currentProject *TFProject
var currentProjectError error

func AnalyzeCurrentProject() (*TFProject, error) {
	once.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			currentProjectError = err
		}
		currentProject, currentProjectError = AnalyzeTerraformProject(wd)
	})

	if currentProjectError != nil {
		return nil, currentProjectError
	}

	result := *currentProject
	return &result, nil
}
