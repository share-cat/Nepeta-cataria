package routes

import (
	"share-cat/Nepeta-cataria/controllers"
	"share-cat/Nepeta-cataria/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSchoolRoutes(router *gin.RouterGroup, controller *controllers.SchoolController) {

	publicSchools := router.Group("/schools")
	{
		publicSchools.GET("", controller.GetSchools)
		publicSchools.GET("/:id", controller.GetSchool)
		publicSchools.GET("/name", controller.GetSchoolByName)
	}

	protectedSchools := router.Group("/schools", middleware.JwtAuth())
	{
		protectedSchools.POST("", controller.CreateSchool)
		protectedSchools.PUT("/:id", controller.UpdateSchool)
		protectedSchools.DELETE("/:id", controller.DeleteSchool)
	}
}
