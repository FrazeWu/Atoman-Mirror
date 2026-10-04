package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"hubproxy/utils"
)

func TestNormalizeRepository(t *testing.T) {
	official := &Repository{Name: "nginx", IsOfficial: true}
	normalizeRepository(official)
	if official.Namespace != "library" || official.Name != "library/nginx" {
		t.Fatalf("official normalized to %#v", official)
	}

	userRepo := &Repository{Name: "owner/app", RepoOwner: "owner"}
	normalizeRepository(userRepo)
	if userRepo.Namespace != "owner" || userRepo.Name != "app" {
		t.Fatalf("user repo normalized to %#v", userRepo)
	}
}

func TestParsePaginationParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?page=3&page_size=50", nil)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	page, pageSize := parsePaginationParams(c, 25)
	if page != 3 || pageSize != 50 {
		t.Fatalf("pagination = %d %d", page, pageSize)
	}
}

func TestParsePaginationParamsClampsBounds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?page=0&page_size=1000", nil)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	page, pageSize := parsePaginationParams(c, 25)
	if page != 1 || pageSize != 100 {
		t.Fatalf("pagination = %d %d, want 1 100", page, pageSize)
	}
}

func TestFetchTagPageHonorsCanceledContext(t *testing.T) {
	utils.InitHTTPClients()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	started := time.Now()
	_, err := fetchTagPage(ctx, "http://127.0.0.1:1", 3)
	if err == nil {
		t.Fatal("expected canceled request error")
	}
	if time.Since(started) > 200*time.Millisecond {
		t.Fatalf("canceled request took too long: %v", time.Since(started))
	}
}

func TestSearchCacheExpires(t *testing.T) {
	cache := &Cache{data: make(map[string]cacheEntry), maxSize: 10}
	cache.SetWithTTL("k", "v", -time.Second)

	if got, ok := cache.Get("k"); ok || got != nil {
		t.Fatalf("expired cache returned: %#v", got)
	}
}
