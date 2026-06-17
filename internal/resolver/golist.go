package resolver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type goListModule struct {
	Path     string
	Version  string
	Update   *goListUpdate
	Indirect bool
	Main     bool
}

type goListUpdate struct {
	Version string
}

func goListUpdates(ctx context.Context, dir string) ([]goListModule, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-u", "-json", "all")
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("go list -m -u all: %s", msg)
	}

	var mods []goListModule
	dec := json.NewDecoder(&stdout)
	for {
		var m goListModule
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode go list output: %w", err)
		}
		mods = append(mods, m)
	}
	return mods, nil
}
