package gomod

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	OldPath    string
	OldVersion string
	NewPath    string
	NewVersion string
}

type Module struct {
	Path      string
	GoVersion string
	Requires  []Require
	Replaces  []Replace
}

func Parse(dir string) (*Module, error) {
	path := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("go.mod not found in %s: not a Go module", dir)
		}
		return nil, fmt.Errorf("read go.mod: %w", err)
	}
	return ParseBytes(path, data)
}

func ParseBytes(name string, data []byte) (*Module, error) {
	f, err := modfile.Parse(name, data, nil)
	if err != nil {
		return nil, fmt.Errorf("parse go.mod: %w", err)
	}

	m := &Module{}
	if f.Module != nil {
		m.Path = f.Module.Mod.Path
	}
	if m.Path == "" {
		return nil, fmt.Errorf("go.mod has no module path")
	}
	if f.Go != nil {
		m.GoVersion = f.Go.Version
	}

	for _, r := range f.Require {
		m.Requires = append(m.Requires, Require{
			Path:     r.Mod.Path,
			Version:  r.Mod.Version,
			Indirect: r.Indirect,
		})
	}
	for _, r := range f.Replace {
		m.Replaces = append(m.Replaces, Replace{
			OldPath:    r.Old.Path,
			OldVersion: r.Old.Version,
			NewPath:    r.New.Path,
			NewVersion: r.New.Version,
		})
	}
	return m, nil
}
