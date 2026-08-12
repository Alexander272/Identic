package transport

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNegotiateEncoding(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "empty", header: "", want: ""},
		{name: "gzip", header: "gzip", want: "gzip"},
		{name: "gzip deflate br", header: "gzip, deflate, br", want: "br"},
		{name: "br", header: "br", want: "br"},
		{name: "gzip preferred over br", header: "br;q=0.8, gzip;q=1", want: "gzip"},
		{name: "both disabled", header: "gzip;q=0, br;q=0", want: ""},
		{name: "identity disabled", header: "identity;q=0", want: ""},
		{name: "wildcard", header: "*", want: "gzip"},
		{name: "deflate only", header: "deflate", want: ""},
		{name: "gzip half br full", header: "gzip;q=0.5, br;q=1", want: "br"},
		{name: "wildcard highest", header: "gzip;q=0.5, br;q=0.5, *;q=1", want: "gzip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, negotiateEncoding(tt.header))
		})
	}
}

func newFrontendFS() fs.FS {
	return fstest.MapFS{
		"frontend/index.html":       {Data: []byte("<!doctype html><title>Identic</title>")},
		"frontend/assets/app.js":    {Data: []byte("console.log('app');")},
		"frontend/assets/app.js.gz": {Data: []byte("gzip-bytes")},
		"frontend/assets/app.js.br": {Data: []byte("br-bytes")},
	}
}

func serveFrontendRouter(t *testing.T, fsys fs.FS) http.Handler {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(func(c *gin.Context) {
		(&Handler{}).serveFrontend(c, fsys)
	})
	return router
}

func TestServeFrontend_NotFoundForMissingAssets(t *testing.T) {
	router := serveFrontendRouter(t, newFrontendFS())

	// Отсутствующий JS-чанк должен давать 404, а не index.html (text/html)
	req := httptest.NewRequest(http.MethodGet, "/assets/missing-abc123.js", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "missing asset should be 404, got body: %s", w.Body.String())

	// Отсутствующий favicon тоже 404
	req = httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestServeFrontend_SpaFallback(t *testing.T) {
	router := serveFrontendRouter(t, newFrontendFS())

	for _, path := range []string{"/", "/orders/abc-def-123", "/search"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "SPA path %s", path)
		assert.Contains(t, w.Body.String(), "Identic", "SPA path %s should serve index.html", path)
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	}
}

func TestServeFrontend_ServesAsset(t *testing.T) {
	router := serveFrontendRouter(t, newFrontendFS())

	req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "console.log('app');", w.Body.String())
	assert.Equal(t, "public, max-age=31536000, immutable", w.Header().Get("Cache-Control"))
}

func TestServeFrontend_ContentEncodings(t *testing.T) {
	router := serveFrontendRouter(t, newFrontendFS())

	cases := []struct {
		name     string
		accept   string
		wantEnc  string
		wantBody string
	}{
		{name: "brotli", accept: "gzip, deflate, br", wantEnc: "br", wantBody: "br-bytes"},
		{name: "gzip", accept: "gzip", wantEnc: "gzip", wantBody: "gzip-bytes"},
		{name: "identity", accept: "identity;q=0, *;q=0", wantEnc: "", wantBody: "console.log('app');"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
			req.Header.Set("Accept-Encoding", tt.accept)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
			assert.Equal(t, tt.wantEnc, w.Header().Get("Content-Encoding"))
		})
	}
}

func TestServeFrontend_API404(t *testing.T) {
	router := serveFrontendRouter(t, newFrontendFS())

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
