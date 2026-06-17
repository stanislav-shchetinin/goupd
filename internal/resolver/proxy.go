package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

const DefaultProxy = "https://proxy.golang.org"

const (
	maxMajorProbe = 50
	maxMajorGap   = 3
)

type ProxyClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewProxyClient(baseURL string, httpClient *http.Client) *ProxyClient {
	if baseURL == "" {
		baseURL = DefaultProxy
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &ProxyClient{BaseURL: strings.TrimRight(baseURL, "/"), HTTP: httpClient}
}

type MajorResult struct {
	Found   bool
	Path    string
	Version string
}

func (c *ProxyClient) LatestMajor(ctx context.Context, modPath, version string) (MajorResult, error) {
	cur := majorNumber(modPath, version)
	base := stripMajorSuffix(modPath)

	start := cur + 1
	if start < 2 {
		start = 2
	}

	var best MajorResult
	gap := 0
	for n := start; n < start+maxMajorProbe && gap < maxMajorGap; n++ {
		cand := base + "/v" + strconv.Itoa(n)
		v, ok, err := c.latest(ctx, cand)
		if err != nil {
			return best, err
		}
		if !ok {
			gap++
			continue
		}
		gap = 0
		best = MajorResult{Found: true, Path: cand, Version: v}
	}
	return best, nil
}

func (c *ProxyClient) latest(ctx context.Context, modPath string) (version string, ok bool, err error) {
	esc, err := module.EscapePath(modPath)
	if err != nil {
		return "", false, nil
	}
	url := c.BaseURL + "/" + esc + "/@latest"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", false, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", false, err
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		var info struct{ Version string }
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return "", false, fmt.Errorf("decode @latest for %s: %w", modPath, err)
		}
		return info.Version, info.Version != "", nil
	case http.StatusNotFound, http.StatusGone:
		return "", false, nil
	default:
		return "", false, fmt.Errorf("proxy %s returned %s", url, resp.Status)
	}
}

func majorNumber(modPath, version string) int {
	if m := semver.Major(version); m != "" {
		if n, err := strconv.Atoi(strings.TrimPrefix(m, "v")); err == nil {
			return n
		}
	}
	if n, ok := majorSuffix(modPath); ok {
		return n
	}
	return 1
}

func stripMajorSuffix(modPath string) string {
	if _, ok := majorSuffix(modPath); ok {
		if i := strings.LastIndex(modPath, "/"); i >= 0 {
			return modPath[:i]
		}
	}
	return modPath
}

func majorSuffix(modPath string) (int, bool) {
	i := strings.LastIndex(modPath, "/")
	if i < 0 {
		return 0, false
	}
	last := modPath[i+1:]
	if len(last) < 2 || last[0] != 'v' {
		return 0, false
	}
	n, err := strconv.Atoi(last[1:])
	if err != nil || n < 2 {
		return 0, false
	}
	return n, true
}
