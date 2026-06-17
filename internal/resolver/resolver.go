package resolver

import (
	"context"
	"sort"

	"golang.org/x/mod/semver"

	"goupd/internal/gomod"
)

type UpdateType string

const (
	Patch UpdateType = "patch"
	Minor UpdateType = "minor"
	Major UpdateType = "major"
)

type Update struct {
	Path     string     `json:"path"`
	Current  string     `json:"current"`
	Latest   string     `json:"latest"`
	Type     UpdateType `json:"type"`
	Indirect bool       `json:"indirect"`
	LatestPath string `json:"latestPath,omitempty"`
}

type Options struct {
	Dir string
	IncludeMajor bool
	DirectOnly bool
	Proxy *ProxyClient
}

func Resolve(ctx context.Context, mod *gomod.Module, opts Options) ([]Update, error) {
	mods, err := goListUpdates(ctx, opts.Dir)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]Update)

	for _, m := range mods {
		if m.Main {
			continue
		}
		if opts.DirectOnly && m.Indirect {
			continue
		}
		if m.Update == nil || m.Update.Version == "" || m.Update.Version == m.Version {
			continue
		}
		updates[m.Path] = Update{
			Path:     m.Path,
			Current:  m.Version,
			Latest:   m.Update.Version,
			Type:     Classify(m.Version, m.Update.Version),
			Indirect: m.Indirect,
		}
	}

	if opts.IncludeMajor && opts.Proxy != nil {
		for _, m := range mods {
			if m.Main || m.Indirect {
				continue
			}
			res, err := opts.Proxy.LatestMajor(ctx, m.Path, m.Version)
			if err != nil || !res.Found {
				continue
			}
			updates[m.Path] = Update{
				Path:       m.Path,
				Current:    m.Version,
				Latest:     res.Version,
				Type:       Major,
				Indirect:   m.Indirect,
				LatestPath: res.Path,
			}
		}
	}

	out := make([]Update, 0, len(updates))
	for _, u := range updates {
		out = append(out, u)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Indirect != out[j].Indirect {
			return !out[i].Indirect
		}
		return out[i].Path < out[j].Path
	})
	return out, nil
}

func Classify(current, latest string) UpdateType {
	switch {
	case semver.Major(current) != semver.Major(latest):
		return Major
	case semver.MajorMinor(current) != semver.MajorMinor(latest):
		return Minor
	default:
		return Patch
	}
}
