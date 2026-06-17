package resolver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeProxy(t *testing.T, versions map[string]string) *ProxyClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const suffix = "/@latest"
		p := r.URL.Path
		if len(p) <= len(suffix) || p[len(p)-len(suffix):] != suffix {
			http.NotFound(w, r)
			return
		}
		mod := p[1 : len(p)-len(suffix)]
		v, ok := versions[mod]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Version":"` + v + `"}`))
	}))
	t.Cleanup(srv.Close)
	return NewProxyClient(srv.URL, srv.Client())
}

func TestLatestMajorFound(t *testing.T) {
	c := fakeProxy(t, map[string]string{
		"github.com/redis/go-redis/v8": "v8.11.5",
		"github.com/redis/go-redis/v9": "v9.5.1",
	})

	res, err := c.LatestMajor(context.Background(), "github.com/redis/go-redis/v6", "v6.15.9+incompatible")
	if err != nil {
		t.Fatalf("LatestMajor: %v", err)
	}
	if !res.Found {
		t.Fatal("expected a major upgrade to be found")
	}
	if res.Path != "github.com/redis/go-redis/v9" {
		t.Errorf("Path = %q, want .../v9", res.Path)
	}
	if res.Version != "v9.5.1" {
		t.Errorf("Version = %q, want v9.5.1", res.Version)
	}
}

func TestLatestMajorNoneFromV1(t *testing.T) {
	c := fakeProxy(t, map[string]string{})

	res, err := c.LatestMajor(context.Background(), "github.com/pkg/errors", "v0.9.1")
	if err != nil {
		t.Fatalf("LatestMajor: %v", err)
	}
	if res.Found {
		t.Errorf("did not expect a major upgrade, got %+v", res)
	}
}

func TestLatestMajorFromV1(t *testing.T) {
	c := fakeProxy(t, map[string]string{
		"github.com/foo/bar/v2": "v2.1.0",
	})

	res, err := c.LatestMajor(context.Background(), "github.com/foo/bar", "v1.0.0")
	if err != nil {
		t.Fatalf("LatestMajor: %v", err)
	}
	if !res.Found || res.Path != "github.com/foo/bar/v2" || res.Version != "v2.1.0" {
		t.Errorf("got %+v, want v2 found", res)
	}
}
