package gindump

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/str"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	fmt.Println(strings.Repeat("-", 78))
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertContains(t *testing.T, msg string, body string, ss ...string) {
	for _, s := range ss {
		if !str.Contains(body, s) {
			t.Errorf(`%s http dump does not contains %q`, msg, s)
		}
	}
}

func TestHttpDump(t *testing.T) {
	router := gin.New()

	buffer := new(bytes.Buffer)
	writer := io.MultiWriter(buffer, os.Stdout)
	router.Use(New(writer).Handler())

	router.Any("/example", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.URL.String())
	})

	buffer.Reset()
	performRequest(router, "GET", "/example?a=100")
	assertContains(t, "GET /example?a=100", buffer.String(), "GET /example?a=100 HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "POST", "/example")
	assertContains(t, "POST /example", buffer.String(), "POST /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	assertContains(t, "PUT /example", buffer.String(), "PUT /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	assertContains(t, "DELETE /example", buffer.String(), "DELETE /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	assertContains(t, "PATCH /example", buffer.String(), "PATCH /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	assertContains(t, "HEAD /example", buffer.String(), "HEAD /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	assertContains(t, "OPTIONS /example", buffer.String(), "OPTIONS /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	assertContains(t, "GET /notfound", buffer.String(), "GET /notfound HTTP/1.1", "HTTP/1.1 404 Not Found")
}
