package gomod

import "testing"

const sample = `module example.com/acme/foo

go 1.22

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.7.0
)

require golang.org/x/sys v0.18.0 // indirect

replace github.com/pkg/errors => github.com/acme/errors v0.9.2
`

func TestParseBytes(t *testing.T) {
	m, err := ParseBytes("go.mod", []byte(sample))
	if err != nil {
		t.Fatalf("ParseBytes: %v", err)
	}
	if m.Path != "example.com/acme/foo" {
		t.Errorf("Path = %q, want example.com/acme/foo", m.Path)
	}
	if m.GoVersion != "1.22" {
		t.Errorf("GoVersion = %q, want 1.22", m.GoVersion)
	}
	if len(m.Requires) != 3 {
		t.Fatalf("len(Requires) = %d, want 3", len(m.Requires))
	}

	byPath := map[string]Require{}
	for _, r := range m.Requires {
		byPath[r.Path] = r
	}
	if r := byPath["golang.org/x/sys"]; !r.Indirect {
		t.Errorf("golang.org/x/sys should be indirect")
	}
	if r := byPath["github.com/spf13/cobra"]; r.Indirect {
		t.Errorf("cobra should be direct")
	}
	if r, ok := byPath["github.com/pkg/errors"]; !ok || r.Version != "v0.9.1" {
		t.Errorf("errors require = %+v, want v0.9.1", r)
	}

	if len(m.Replaces) != 1 {
		t.Fatalf("len(Replaces) = %d, want 1", len(m.Replaces))
	}
	if m.Replaces[0].NewPath != "github.com/acme/errors" {
		t.Errorf("replace NewPath = %q", m.Replaces[0].NewPath)
	}
}

func TestParseBytesNoModule(t *testing.T) {
	if _, err := ParseBytes("go.mod", []byte("go 1.22\n")); err == nil {
		t.Fatal("expected error for go.mod without module path")
	}
}
