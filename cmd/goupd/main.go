package main

import (
	"context"
	"fmt"
	"os"

	"goupd/internal/cli"
	"goupd/internal/gomod"
	"goupd/internal/repo"
	"goupd/internal/report"
	"goupd/internal/resolver"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "goupd: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := cli.Parse(args, os.Stderr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	dir, cleanup, err := repo.Clone(ctx, cfg.Repo, cfg.Ref)
	if err != nil {
		return err
	}
	defer cleanup()

	mod, err := gomod.Parse(dir)
	if err != nil {
		return err
	}

	opts := resolver.Options{
		Dir:          dir,
		IncludeMajor: cfg.Major,
		DirectOnly:   cfg.DirectOnly,
	}
	if cfg.Major {
		opts.Proxy = resolver.NewProxyClient("", nil)
	}

	updates, err := resolver.Resolve(ctx, mod, opts)
	if err != nil {
		return err
	}

	rep := report.Report{
		Module:    mod.Path,
		GoVersion: mod.GoVersion,
		Updates:   updates,
	}

	switch cfg.Format {
	case "json":
		return report.WriteJSON(os.Stdout, rep)
	default:
		return report.WriteText(os.Stdout, rep)
	}
}
