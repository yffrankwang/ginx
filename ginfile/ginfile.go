package ginfile

import (
	"bytes"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Public1Year Cache-Control: public, max-age=31536000
const Public1Year = "public, max-age=31536000"

func getCacheControlWriter(c *gin.Context, cacheControl string) http.ResponseWriter {
	if cacheControl == "" {
		return c.Writer
	}
	h := map[string]string{"Cache-Control": cacheControl}
	return &headerWriter{c.Writer, h}
}

// headerWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headerWriter struct {
	http.ResponseWriter
	header map[string]string
}

// WriteHeader append header when statusCode == 200
func (hw *headerWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		for k, v := range hw.header {
			if hw.Header().Get(k) == "" {
				hw.Header().Add(k, v)
			}
		}
	}
	hw.ResponseWriter.WriteHeader(statusCode)
}

// AppendPrefix returns a handler that serves HTTP requests by appending the
// given prefix from the request URL's Path (and RawPath if set) and invoking
// the handler hh.
func appendPrefix(prefix string, hh http.Handler) http.Handler {
	if prefix == "" {
		return hh
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := prefix + r.URL.Path
		rp := prefix + r.URL.RawPath
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

// URLReplace returns a handler that serves HTTP requests by replacing the
// request URL's Path (and RawPath if set) (use strings.Replace(path, src, des) and invoking
// the handler hh.
func urlReplace(src, des string, hh http.Handler) http.Handler {
	if src == "" || src == des {
		return hh
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.Replace(r.URL.Path, src, des, 1)
		rp := strings.Replace(r.URL.RawPath, src, des, 1)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

// Static serves files from the given file system root.
func Static(g *gin.RouterGroup, relativePath, localPath, cacheControl string) {
	StaticFS(g, relativePath, "", http.Dir(localPath), cacheControl)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// ginfile.StaticFSFile(gin, "favicon.ico", "./resources/favicon.ico", "public")
func StaticFile(g *gin.RouterGroup, relativePath, localPath, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := func(c *gin.Context) {
		ccw := getCacheControlWriter(c, cacheControl)
		http.ServeFile(ccw, c.Request, localPath)
	}
	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func StaticFS(g *gin.RouterGroup, relativePath string, localPath string, hfs http.FileSystem, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(g.BasePath(), relativePath)
	fileServer := http.FileServer(hfs)
	if prefix == "" || prefix == "/" {
		fileServer = appendPrefix(localPath, fileServer)
	} else if localPath == "" || localPath == "." {
		fileServer = http.StripPrefix(prefix, fileServer)
	} else {
		fileServer = urlReplace(prefix, localPath, fileServer)
	}

	handler := func(c *gin.Context) {
		ccw := getCacheControlWriter(c, cacheControl)
		fileServer.ServeHTTP(ccw, c.Request)
	}

	urlPattern := path.Join(relativePath, "/*path")

	// Register GET and HEAD handlers
	g.GET(urlPattern, handler)
	g.HEAD(urlPattern, handler)
}

// StaticFSFile registers a single route in order to serve a single file of the filesystem.
// ginfile.StaticFSFile(gin, "favicon.ico", "./resources/favicon.ico", hfs, ginfile.Public1Year)
func StaticFSFile(g *gin.RouterGroup, relativePath, filePath string, hfs http.FileSystem, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := func(c *gin.Context) {
		defer func(old string) {
			c.Request.URL.Path = old
		}(c.Request.URL.Path)

		c.Request.URL.Path = filePath
		ccw := getCacheControlWriter(c, cacheControl)

		http.FileServer(hfs).ServeHTTP(ccw, c.Request)
	}

	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}

// StaticContent registers a single route in order to serve a single file of the data.
// //go:embed favicon.ico
// var favicon []byte
// ginfile.StaticContent(gin, "favicon.ico", favicon, time.Now(), "public")
func StaticContent(g *gin.RouterGroup, relativePath string, data []byte, modtime time.Time, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	if modtime.IsZero() {
		modtime = time.Now()
	}
	handler := func(c *gin.Context) {
		if cacheControl != "" {
			c.Header("Cache-Control", cacheControl)
		}
		name := filepath.Base(c.Request.URL.Path)
		http.ServeContent(c.Writer, c.Request, name, modtime, bytes.NewReader(data))
	}
	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}
