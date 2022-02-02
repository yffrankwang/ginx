package ginhtml

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/str"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertContains(t *testing.T, msg string, body string, ss ...string) {
	for _, s := range ss {
		if !str.Contains(body, s) {
			t.Errorf(`%s access log does not contains %q`, msg, s)
		}
	}
}

func TestHtml(t *testing.T) {
	router := gin.Default()

	// new template engine
	ghe := NewEngine()
	if err := ghe.Load("testdata"); err != nil {
		panic(err)
	}

	// customize gin html render
	router.HTMLRender = ghe

	router.GET("/", func(ctx *gin.Context) {
		// render
		ctx.HTML(http.StatusOK, "index", &Result{"Index title!"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	performRequest(router, "GET", "/")
	assertContains(t, "GET /", w.Body.String(), "<title>Index title!</title>", "<body>3</body>")
}
