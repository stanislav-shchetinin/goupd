package resolver

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		current, latest string
		want            UpdateType
	}{
		{"v1.2.3", "v1.2.4", Patch},
		{"v1.2.3", "v1.3.0", Minor},
		{"v1.2.3", "v2.0.0", Major},
		{"v0.9.0", "v0.9.1", Patch},
		{"v0.9.0", "v0.10.0", Minor},
		{"v6.15.9+incompatible", "v6.15.10+incompatible", Patch},
	}
	for _, tt := range tests {
		if got := Classify(tt.current, tt.latest); got != tt.want {
			t.Errorf("Classify(%q, %q) = %q, want %q", tt.current, tt.latest, got, tt.want)
		}
	}
}

func TestMajorSuffix(t *testing.T) {
	tests := []struct {
		path   string
		wantN  int
		wantOK bool
	}{
		{"github.com/foo/bar", 0, false},
		{"github.com/foo/bar/v2", 2, true},
		{"github.com/foo/bar/v10", 10, true},
		{"github.com/foo/v1", 0, false},
		{"github.com/foo/version", 0, false},
	}
	for _, tt := range tests {
		n, ok := majorSuffix(tt.path)
		if n != tt.wantN || ok != tt.wantOK {
			t.Errorf("majorSuffix(%q) = (%d, %v), want (%d, %v)", tt.path, n, ok, tt.wantN, tt.wantOK)
		}
	}
}

func TestStripMajorSuffix(t *testing.T) {
	if got := stripMajorSuffix("github.com/foo/bar/v3"); got != "github.com/foo/bar" {
		t.Errorf("stripMajorSuffix = %q", got)
	}
	if got := stripMajorSuffix("github.com/foo/bar"); got != "github.com/foo/bar" {
		t.Errorf("stripMajorSuffix = %q", got)
	}
}

func TestMajorNumber(t *testing.T) {
	if got := majorNumber("github.com/foo/bar", "v1.2.3"); got != 1 {
		t.Errorf("majorNumber = %d, want 1", got)
	}
	if got := majorNumber("github.com/foo/bar/v3", "v3.1.0"); got != 3 {
		t.Errorf("majorNumber = %d, want 3", got)
	}
	if got := majorNumber("github.com/foo/bar/v4", ""); got != 4 {
		t.Errorf("majorNumber from path = %d, want 4", got)
	}
}
