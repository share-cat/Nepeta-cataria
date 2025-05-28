package routes

import (
	"share-cat/Nepeta-cataria/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupSwaggerRoutes(router *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
