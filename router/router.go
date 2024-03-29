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

	return &DefaultRouter{
		Router: router,
	}
}
