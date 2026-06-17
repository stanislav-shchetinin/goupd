package repo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Clone(ctx context.Context, src, ref string) (dir string, cleanup func(), err error) {
	noop := func() {}

	if info, statErr := os.Stat(src); statErr == nil && info.IsDir() {
		return src, noop, nil
	}

	tmp, err := os.MkdirTemp("", "goupd-*")
	if err != nil {
		return "", noop, fmt.Errorf("create temp dir: %w", err)
	}
	cleanup = func() { _ = os.RemoveAll(tmp) }

	args := []string{"clone", "--depth", "1"}
	if ref != "" {
		args = append(args, "--branch", ref)
	}
	args = append(args, src, tmp)

	cmd := exec.CommandContext(ctx, "git", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if runErr := cmd.Run(); runErr != nil {
		cleanup()
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = runErr.Error()
		}
		return "", noop, fmt.Errorf("git clone %q: %s", src, msg)
	}

	return tmp, cleanup, nil
}
