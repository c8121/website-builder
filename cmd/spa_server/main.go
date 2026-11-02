package main

import (
	"github.com/c8121/asset-storage/shared/rest_api"
	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/config"
	"github.com/c8121/website-builder/internal/spa_server"
	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadDefault()

	router := gin.Default()
	util.LogError(router.SetTrustedProxies(nil))

	rest_api.CreateRoutes(router)
	spa_server.CreateRoutes(router)

	if config.CertFile != "" {
		util.PanicOnError(router.RunTLS(config.ListenAddress, config.CertFile, config.KeyFile), "Failed to start server")
	} else {
		util.PanicOnError(router.Run(config.ListenAddress), "Failed to start server")
	}
}
