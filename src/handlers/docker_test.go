package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"hubproxy/config"
)

func TestParseRegistryPath(t *testing.T) {
	tests := []struct {
		path      string
		image     string
		apiType   string
		reference string
	}{
		{"library/nginx/manifests/latest", "library/nginx", "manifests", "latest"},
		{"library/nginx/blobs/sha256:abc", "library/nginx", "blobs", "sha256:abc"},
		{"library/nginx/tags/list", "library/nginx", "tags", "list"},
	}

	for _, tt := range tests {
		image, apiType, reference := parseRegistryPath(tt.path)
		if image != tt.image || apiType != tt.apiType || reference != tt.reference {
			t.Fatalf("parseRegistryPath(%q) = %q %q %q", tt.path, image, apiType, reference)
		}
	}
}

func TestParseRegistryPathInvalid(t *testing.T) {
	image, apiType, reference := parseRegistryPath("library/nginx/unknown/latest")
	if image != "" || apiType != "" || reference != "" {
		t.Fatalf("invalid path parsed as %q %q %q", image, apiType, reference)
	}
}

func TestResolveAuthHost(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := []byte(`
[registries."ghcr.io"]
upstream = "ghcr.io"
authHost = "ghcr.io/token"
authType = "github"
enabled = true

[registries."quay.io"]
upstream = "quay.io"
authHost = "quay.io/v2/auth"
authType = "quay"
enabled = true
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", path)
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
	}

	if got := resolveAuthHost(""); got != "" {
		t.Fatalf("empty service = %q", got)
	}
	if got := resolveAuthHost("registry.docker.io"); got != "" {
		t.Fatalf("docker hub service = %q", got)
	}
	if got := resolveAuthHost("ghcr.io"); got != "ghcr.io/token" {
		t.Fatalf("ghcr.io = %q", got)
	}
	if got := resolveAuthHost("quay.io"); got != "quay.io/v2/auth" {
		t.Fatalf("quay.io = %q", got)
	}
}

func TestBuildDockerAuthURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := []byte(`
[registries."ghcr.io"]
upstream = "ghcr.io"
authHost = "ghcr.io/token"
enabled = true
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", path)
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)

	t.Run("docker hub keeps path", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/token?service=registry.docker.io&scope=repository:library/nginx:pull", nil)
		got := buildDockerAuthURL(c)
		want := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/nginx:pull"
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})

	t.Run("ghcr uses AuthHost without duplicating path", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/token?service=ghcr.io&scope=repository:foo/bar:pull", nil)
		got := buildDockerAuthURL(c)
		want := "https://ghcr.io/token?service=ghcr.io&scope=repository:foo/bar:pull"
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})
}

func TestBuildDockerAuthCacheKeyScopesPublicRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newContext := func(rawURL, host, authorization string) *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, rawURL, nil)
		c.Request.Host = host
		if authorization != "" {
			c.Request.Header.Set("Authorization", authorization)
		}
		return c
	}

	base := buildDockerAuthCacheKey(newContext("/token?scope=repository:library/nginx:pull&service=registry.docker.io", "mirror.example", ""))
	canonical := buildDockerAuthCacheKey(newContext("/token?service=registry.docker.io&scope=repository:library/nginx:pull", "mirror.example", ""))
	if base == "" || base != canonical {
		t.Fatalf("equivalent public requests produced different keys: %q %q", base, canonical)
	}
	if base == buildDockerAuthCacheKey(newContext("/token?service=ghcr.io&scope=repository:library/nginx:pull", "mirror.example", "")) {
		t.Fatal("different registry services shared a token cache key")
	}
	if base == buildDockerAuthCacheKey(newContext("/token/alternate?scope=repository:library/nginx:pull&service=registry.docker.io", "mirror.example", "")) {
		t.Fatal("different auth paths shared a token cache key")
	}
	if base == buildDockerAuthCacheKey(newContext("/token?scope=repository:library/nginx:pull&service=registry.docker.io", "other.example", "")) {
		t.Fatal("different hosts shared a token cache key")
	}
	post := newContext("/token?scope=repository:library/nginx:pull&service=registry.docker.io", "mirror.example", "")
	post.Request.Method = http.MethodPost
	if base == buildDockerAuthCacheKey(post) {
		t.Fatal("different auth methods shared a token cache key")
	}
	if got := buildDockerAuthCacheKey(newContext("/token?scope=repository:library/nginx:pull&service=registry.docker.io", "mirror.example", "Bearer secret")); got != "" {
		t.Fatalf("authenticated request should bypass shared cache, got key %q", got)
	}
}

func TestRewriteAuthHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := []byte(`
[registries."quay.io"]
upstream = "quay.io"
authHost = "quay.io/v2/auth"
enabled = true

[registries."ghcr.io"]
upstream = "ghcr.io"
authHost = "ghcr.io/token"
enabled = true
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", path)
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
	}

	got := rewriteAuthHeader(`Bearer realm="https://quay.io/v2/auth",service="quay.io"`, "proxy.example.com")
	want := `Bearer realm="http://proxy.example.com/token",service="quay.io"`
	if got != want {
		t.Fatalf("quay rewrite: got %q want %q", got, want)
	}

	got = rewriteAuthHeader(`Bearer realm="https://ghcr.io/token",service="ghcr.io"`, "proxy.example.com")
	want = `Bearer realm="http://proxy.example.com/token",service="ghcr.io"`
	if got != want {
		t.Fatalf("ghcr rewrite: got %q want %q", got, want)
	}

	got = rewriteAuthHeader(`Bearer realm="https://auth.docker.io/token",service="registry.docker.io"`, "proxy.example.com")
	want = `Bearer realm="http://proxy.example.com/token",service="registry.docker.io"`
	if got != want {
		t.Fatalf("docker hub rewrite: got %q want %q", got, want)
	}
}

func TestRewriteAuthHeaderPreservesForwardedHTTPS(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	data := []byte(`[registries."quay.io"]
upstream = "quay.io"
authHost = "quay.io/v2/auth"
enabled = true
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", path)
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
	}

	got := rewriteAuthHeaderWithScheme(`Bearer realm="https://quay.io/v2/auth",service="quay.io"`, "https", "mirror.example")
	want := `Bearer realm="https://mirror.example/token",service="quay.io"`
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProxyOriginUsesForwardedHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/token", nil)
	c.Request.Host = "127.0.0.1:50001"
	c.Request.Header.Set("X-Forwarded-Proto", "https, http")
	c.Request.Header.Set("X-Forwarded-Host", "mirror.example, proxy.internal")

	scheme, host := proxyOrigin(c)
	if scheme != "https" || host != "mirror.example" {
		t.Fatalf("proxy origin = %s://%s, want https://mirror.example", scheme, host)
	}
}
