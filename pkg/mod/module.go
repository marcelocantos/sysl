package mod

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var SyslModules = os.Getenv("SYSL_MODULES") != "off"

type Module struct {
	Name    string
	Dir     string
	Version string
}

type Modules []*Module

func (m *Modules) Add(v *Module) {
	*m = append(*m, v)
}

func (m *Modules) GetByFilename(filename, ver string) *Module {
	for _, mod := range *m {
		if hasPathPrefix(mod.Name, filename) {
			if ver != "" && mod.Version != ver {
				continue
			}
			return mod
		}
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Find(name string, ver string) (*Module, error) {
	if !fileExists("go.mod") {
		return nil, errors.New("no go.mod file, run `go mod init` first")
	}

	err := goGetByFilepath(name, ver)
	if err != nil {
		return nil, fmt.Errorf("%s not found", name)
	}

	var modules Modules
	if err = modules.Load(); err != nil {
		return nil, fmt.Errorf("error loading modules: %s", err.Error())
	}

	if modules == nil {
		return nil, fmt.Errorf("modules list is empty")
	}

	mod := modules.GetByFilename(name, ver)
	if mod == nil {
		return nil, fmt.Errorf("error finding module of file %s", name)
	}
	return mod, nil
}

func hasPathPrefix(prefix, s string) bool {
	prefix = filepath.Clean(prefix)
	s = filepath.Clean(s)

	if len(s) > len(prefix) {
		return s[len(prefix)] == filepath.Separator && s[:len(prefix)] == prefix
	}

	return s == prefix
}
