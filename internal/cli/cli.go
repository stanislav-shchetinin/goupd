package cli

import (
	"flag"
	"fmt"
	"io"
	"time"
)

type Config struct {
	Repo       string
	Ref        string
	Format     string
	Major      bool
	DirectOnly bool
	Timeout    time.Duration
}

func Parse(args []string, out io.Writer) (*Config, error) {
	fs := flag.NewFlagSet("goupd", flag.ContinueOnError)
	fs.SetOutput(out)

	cfg := &Config{}
	fs.StringVar(&cfg.Ref, "ref", "", "git branch, tag or commit to check out (default: repository default branch)")
	fs.StringVar(&cfg.Format, "format", "text", "output format: text or json")
	fs.BoolVar(&cfg.Major, "major", true, "include major (v2+) upgrades discovered via the module proxy")
	fs.BoolVar(&cfg.DirectOnly, "direct-only", false, "report only direct dependencies")
	fs.DurationVar(&cfg.Timeout, "timeout", 2*time.Minute, "overall timeout for network operations")

	fs.Usage = func() {
		_, _ = fmt.Fprintf(out, "goupd - report Go module info and available dependency updates\n\n")
		_, _ = fmt.Fprintf(out, "Usage:\n  goupd [flags] <git-repo-url|local-path>\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	rest := fs.Args()
	if len(rest) == 0 {
		fs.Usage()
		return nil, fmt.Errorf("missing required argument: repository address")
	}
	if len(rest) > 1 {
		return nil, fmt.Errorf("unexpected extra arguments: %v", rest[1:])
	}
	cfg.Repo = rest[0]

	switch cfg.Format {
	case "text", "json":
	default:
		return nil, fmt.Errorf("invalid -format %q: must be \"text\" or \"json\"", cfg.Format)
	}

	return cfg, nil
}
