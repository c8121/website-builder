package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/c8121/go_vue_spa"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

func main() {

	rootDir := flag.String("root", "", "Root directory")
	listen := flag.String("listen", "localhost:8080", "Listen Address (ip:port)")
	flag.Parse()

	if *rootDir == "" || !util.IsExistingDirectory(*rootDir, slog.LevelError) {
		fmt.Println("Root directory does not exist")
		return
	}

	if *listen == "" {
		fmt.Println("You must specify a listen address")
		return
	}

	router := gin.Default()
	util.LogError(router.SetTrustedProxies(nil))

	resource := go_vue_spa.Resource{
		RootFs:          os.DirFS(*rootDir),
		DefaultDocument: "index.html",
	}
	router.NoRoute(func(c *gin.Context) {
		if go_vue_spa.ServeFromResource(c, resource) {
			return
		}
		http.NotFound(c.Writer, c.Request)
	})

	util.PanicOnError(router.Run(*listen), "Failed to start server")
}
