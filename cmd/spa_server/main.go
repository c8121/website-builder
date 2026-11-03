package main

import (
	"github.com/c8121/asset-storage/shared/rest_api"
	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/spa_server"
	boostrap "github.com/c8121/website-builder/shared/bootstrap"
	"github.com/c8121/website-builder/shared/config"
	"github.com/gin-gonic/gin"
)

func main() {

	boostrap.InitWebsiteBuiderApp()
	defer boostrap.FinalizeWebsiteBuiderApp()

	router := gin.Default()
	util.LogError(router.SetTrustedProxies(nil))

	rest_api.CreateRoutes(router)
	spa_server.CreateRoutes(router)

	if config.AppConfig.CertFile != "" {
		util.PanicOnError(router.RunTLS(config.AppConfig.ListenAddress, config.AppConfig.CertFile, config.AppConfig.KeyFile), "Failed to start server")
	} else {
		util.PanicOnError(router.Run(config.AppConfig.ListenAddress), "Failed to start server")
	}
}
