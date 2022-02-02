package ginfile

import (
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed testdata
var fsdata embed.FS

func Example() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// static serve path: "/" -> "./testdir" with "private" cache-control
	Static(&router.RouterGroup, "/", "./testdir", "private")

	// static serve file: "/r1.txt" -> "./file.txt" with "public" cache-control
	StaticFile(&router.RouterGroup, "/r1.txt", "./file.txt", "public")

	// static serve FS path: "/fs" -> "fs:/fsdir" with "public" cache-control
	StaticFS(&router.RouterGroup, "/fs", "/fsdir", http.FS(fsdata), "public")

	// static serve FS file: "/r2.txt" -> "fs:/fsdir/r2.txt" with "public" cache-control
	StaticFSFile(&router.RouterGroup, "/r2.txt", "fsdir/r2.txt", http.FS(fsdata), "public")

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/t1.txt", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
