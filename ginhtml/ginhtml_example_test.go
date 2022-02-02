package ginhtml

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
	ghe := NewEngine()
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
