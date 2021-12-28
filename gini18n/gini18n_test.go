package gini18n

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func doTest(t *testing.T, req *http.Request, want string) {
	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(NewLocalizer("ja", "zh").Handler())
	router.Any("/", func(c *gin.Context) {
		c.String(200, GetLocale(c))
	})

	router.ServeHTTP(w, req)

	if w.Body.String() != want {
		t.Errorf("%v = %q, want %q", req.URL.String(), w.Body.String(), want)
	}
}

func TestAcceptLanguages1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Languages", "en;en-US")

	doTest(t, req, "ja")
}

func TestAcceptLanguages2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Languages", "ja;zh")

	doTest(t, req, "ja")
}

func TestHttpHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add(HeaderName, "zh;ja")

	doTest(t, req, "zh")
}

func TestQueryString(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?__locale=zh", nil)

	doTest(t, req, "zh")
}

func TestPostForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	req.PostForm = url.Values{}
	req.PostForm.Add("__locale", "zh")

	doTest(t, req, "zh")
}

func TestCookie(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: "zh",
	})

	doTest(t, req, "zh")
}
