package tf_project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Module struct {
	Name string `tfsdk:"name" json:"name"`
	Path string `tfsdk:"path" json:"path"`
}

type TFProject struct {
	Path          string             `tfsdk:"path" json:"path"`
	Modules       map[string]*Module `tfsdk:"modules" json:"modules"`
	CurrentModule *Module            `tfsdk:"current_module" json:"current_module"`
}

var ProjectRootMark = filepath.Join(".git", "HEAD")

func detectProjectRoot(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ProjectRootMark))
	return err == nil
}

func isTFModule(dir string) bool {
	tfStateStat, err := os.Stat(filepath.Join(dir, ".terraform", "terraform.tfstate"))
	if err == nil && tfStateStat.Size() > 0 {
		return true
	}

	tfStateStat, err = os.Stat(filepath.Join(dir, "terraform.tfstate"))
	if err == nil && tfStateStat.Size() > 0 {
		return true
	}

	return false
}

func AnalyzeTerraformProject(targetDri string) (*TFProject, error) {
	projectRoot := targetDri

	maxDepth := 64
	for {
		if maxDepth == 0 {
			return nil, fmt.Errorf("project root not found")
		}
		maxDepth--

		if detectProjectRoot(projectRoot) {
			break
		}
		parrent := filepath.Dir(projectRoot)
		if parrent == projectRoot {
			return nil, fmt.Errorf("project root not found")
		}
		projectRoot = parrent
	}

	modules := map[string]*Module{}
	walkerr := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if errors.Is(err, fs.ErrPermission) {
			return filepath.SkipDir
		}

		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if info.Name() == ".terraform" {
			return filepath.SkipDir
		}

		if info.Name() == ".secret" {
			return filepath.SkipDir
		}

		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !isTFModule(path) {
			return nil
		}

		relPath, err := filepath.Rel(projectRoot, path)
		if err != nil {
			return err
		}

		modules[relPath] = &Module{Path: path, Name: relPath}

		return nil
	})
	if walkerr != nil {
		return nil, walkerr
	}

	var wdModule *Module
	if wdRelPath, err := filepath.Rel(projectRoot, targetDri); err == nil {
		wdModule = modules[wdRelPath]
	}

	return &TFProject{
		Path:          projectRoot,
		Modules:       modules,
		CurrentModule: wdModule,
	}, nil
}
