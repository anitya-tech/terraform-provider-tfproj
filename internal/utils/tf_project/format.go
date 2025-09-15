package tf_project

import (
	"os"
	"path/filepath"

	"github.com/valyala/fasttemplate"
)

func (proj *TFProject) Format(template string) string {
	t := fasttemplate.New(template, "{", "}")

	homedir, _ := os.UserHomeDir()

	return t.ExecuteStringStd(map[string]any{
		"~": homedir,

		"project.path": proj.Path,

		"secret.path":   filepath.Join(proj.Path, ".secret"),
		"secret.states": filepath.Join(proj.Path, ".secret", "states"),
		"secret.store":  filepath.Join(proj.Path, ".secret", "store"),

		"module.name":   proj.CurrentModule.Name,
		"module.path":   proj.CurrentModule.Path,
		"module.states": filepath.Join(proj.Path, ".secret", "states", proj.CurrentModule.Name),
		"module.store":  filepath.Join(proj.Path, ".secret", "store", proj.CurrentModule.Name),
	})
}
