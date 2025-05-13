package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupSwaggerRoutes(router *gin.RouterGroup) {

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
