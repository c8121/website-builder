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

}
