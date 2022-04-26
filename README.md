 GINX
=====================================================================

[![Build Status](https://github.com/yffrankwang/ginx/actions/workflows/build.yml/badge.svg)](https://github.com/yffrankwang/ginx/actions?query=branch%3Amaster) [![codecov](https://codecov.io/gh/yffrankwang/ginx/branch/master/graph/badge.svg)](https://codecov.io/gh/pandafw/pango) [![MIT](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/licenses/MIT)



GINX is a GO development utility library for [GIN](https://github.com/gin-gonic/gin).

| **Package**           | **Description**                         |
| :-------------------- | :-------------------------------------- |
| [gindump](#gindump)   | a http request/response dumper middleware for gin               |
| [ginfile](#ginfile)   | a static file handler with Cache-Control header support for gin |
| [gingzip](#gingzip)   | a gzip encoding support middleware for gin                      |
| [ginhtml](#ginhtml)   | a html template engine for gin                                  |
| [gini18n](#gini18n)   | a localizer middleware for gin                              |
| [ginlog](#ginlog)     | a access logger middleware for gin                              |

## Install:

	go get github.com/yffrankwang/ginx


## gindump
A http request/response dumper middleware for gin.


### Example:

```golang
import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/gindump"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gindump.New(os.Stdout).Handler())

	router.Any("/example", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.URL.String())
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/example?a=100", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

```

Output:
```
>>>>>>>> 2021-12-04T16:33:02.087 4264618929a66a101ca5a28625fdf29a832574c5 >>>>>>>>
GET /example?a=100 HTTP/1.1
Host: example.com



<<<<<<<< 2021-12-04T16:33:02.088 4264618929a66a101ca5a28625fdf29a832574c5 <<<<<<<<
HTTP/1.1 200 OK
Connection: close
Content-Type: text/plain; charset=utf-8

/example?a=100

```

## ginfile
A static file handler with Cache-Control header support for gin.

### Example:

```golang
import (
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/ginfile"
)

//go:embed testdata
var fsdata embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// static serve path: "/" -> "./testdir" with "private" cache-control
	ginfile.Static(&router.RouterGroup, "/", "./testdir", "private")

	// static serve file: "/r1.txt" -> "./file.txt" with "public" cache-control
	ginfile.StaticFile(&router.RouterGroup, "/r1.txt", "./file.txt", "public")

	// static serve FS path: "/fs" -> "fs:/fsdir" with "public" cache-control
	ginfile.StaticFS(&router.RouterGroup, "/fs", "/fsdir", http.FS(fsdata), "public")

	// static serve FS file: "/r2.txt" -> "fs:/fsdir/r2.txt" with "public" cache-control
	ginfile.StaticFSFile(&router.RouterGroup, "/r2.txt", "fsdir/file.txt", http.FS(fsdata), "public")

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/r1.txt", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
```


## gingzip
A gzip encoding support middleware for gin.

### Example:

```golang
import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/gingzip"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gingzip.Default().Handler())
	router.GET("/", func(c *gin.Context) {
		c.String(200, strings.Repeat("This is a Test!\n", 1000))
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
```


## ginhtml
A html template engine for gin.

### Example:

```golang
import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Title string
}

func (r *Result) Add(a, b int) int {
	return a + b
}

func Example() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// html template engine
	ghe := ginhtml.NewEngine()
	if err := ghe.Load("./testdata"); err != nil {
		panic(err)
	}

	// customize gin html render
	router.HTMLRender = ghe

	router.GET("/", func(ctx *gin.Context) {
		// render
		ctx.HTML(http.StatusOK, "index", &Result{"Index title!"})
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/", nil)
	client := &http.Client{Timeout: time.Second * 1}
	res, _ := client.Do(req)

	io.Copy(os.Stdout, res.Body)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
```

./testdata/index.html
```html
<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>{{.Add 1 2}}</body>
</html>
```


## ginlog
A access logger middleware for gin.

### Example:

```golang
import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/ginlog"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(ginlog.New(os.Stdout, ginlog.DefaultTextLogFormat).Handler())

	router.Any("/example", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.URL.String())
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/example?a=100", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
```

Output:
```
2021-12-04T14:30:30.840	200	0	-1	127.0.0.1	127.0.0.1:1234		GET	example.com	/example?a=100
```


## gini18n
A localizer middleware for gin.

### Example:

```golang
import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yffrankwang/ginx/gini18n"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gini18n.NewLocalizer("ja", "zh", "en").Handler())

	router.Any("/example", func(c *gin.Context) {
		locale := gini18n.GetLocale(c)
		c.String(http.StatusOK, locale)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/example?a=100&__locale=ja", nil)
	req.Header.Add("Accept-Languages", "ja;zh")

	client := &http.Client{Timeout: time.Second * 1}
	res, _ := client.Do(req)

	raw, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(raw))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
```

Output:
```
ja
```
