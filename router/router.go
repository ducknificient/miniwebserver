package router

import (
	"miniwebserver/config"
	"miniwebserver/controller"

	"github.com/gin-gonic/gin"
)

type DefaultRouter struct {
	Router     *gin.Engine
	controller controller.Controller
}

func NewRouter(config config.Configuration, controller controller.Controller) *DefaultRouter {

	if *config.GetConfiguration().Production == `true` {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/", gin.WrapF(controller.Root))
	router.OPTIONS("/*path", gin.WrapF(controller.Options))
	router.GET("/about/version", gin.WrapF(controller.About))
	router.GET("/serve/:path/:file", gin.WrapF(controller.ServeFile))
	router.GET("/upload/", gin.WrapF(controller.UploadPage))
	router.POST("/upload/file", gin.WrapF(controller.UploadFile))

	return &DefaultRouter{
		Router: router,
	}
}
