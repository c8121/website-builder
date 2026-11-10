package rest_api

import (
	auth "github.com/c8121/go_vue_spa/auth/gin"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router *gin.Engine) {

	router.GET("/pages/list", auth.UserAuthRequiredHandler(ListPages))
	router.GET("/pages/list/:parent", auth.UserAuthRequiredHandler(ListPages))

	router.GET("/page/:id", auth.UserAuthRequiredHandler(GetPage))
	router.POST("/page/add", auth.UserAuthRequiredHandler(AddPage))
	router.POST("/page/delete", auth.UserAuthRequiredHandler(DeletePage))
	router.POST("/pages/reorder", auth.UserAuthRequiredHandler(ReorderPages))

	router.GET("/content/list/:page", auth.UserAuthRequiredHandler(ListPageContents))

	router.GET("/content/:id", auth.UserAuthRequiredHandler(GetPageContent))
	router.POST("/content/add", auth.UserAuthRequiredHandler(AddPageContent))
	router.POST("/content/delete", auth.UserAuthRequiredHandler(DeletePageContent))
	router.POST("/content/reorder", auth.UserAuthRequiredHandler(ReorderPageContents))

	router.GET("/sitesconfigs/list", auth.UserAuthRequiredHandler(ListSiteConfig))

	router.GET("/sitesconfig/:id", auth.UserAuthRequiredHandler(GetSiteConfig))
	router.POST("/sitesconfig/add", auth.UserAuthRequiredHandler(AddSiteConfig))
	router.POST("/sitesconfig/delete", auth.UserAuthRequiredHandler(DeleteSiteConfig))
	router.POST("/sitesconfig/reorder", auth.UserAuthRequiredHandler(ReorderSiteConfigs))

}
