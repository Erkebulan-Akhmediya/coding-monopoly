package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticHandlerServesIndexAndAssets(t *testing.T) {
	h := staticHandler()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / status %d", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	html := string(body)
	if !strings.Contains(html, "<html") {
		t.Fatalf("expected HTML index, got %q", truncate(html, 80))
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("content-type %q", ct)
	}

	// Hashed asset refs in index.html must resolve through the same handler.
	for _, prefix := range []string{`src="/assets/`, `href="/assets/`} {
		rest := html
		for {
			i := strings.Index(rest, prefix)
			if i < 0 {
				break
			}
			rest = rest[i+len(prefix):]
			end := strings.IndexAny(rest, `"'`)
			if end < 0 {
				t.Fatalf("unclosed asset ref after %q", prefix)
			}
			path := "/assets/" + rest[:end]
			assetRec := httptest.NewRecorder()
			assetReq := httptest.NewRequest(http.MethodGet, path, nil)
			h.ServeHTTP(assetRec, assetReq)
			if assetRec.Code != http.StatusOK {
				t.Fatalf("GET %s status %d (stale embed? run make build-client)", path, assetRec.Code)
			}
			rest = rest[end:]
		}
	}

	// Unknown path falls back to index.html (SPA / ?admin=1 deep links).
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	h.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("SPA fallback status %d", rec2.Code)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
