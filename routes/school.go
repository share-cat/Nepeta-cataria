package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupSchoolRoutes(router *gin.RouterGroup, controller *controllers.SchoolController) {
	schools := router.Group("/schools")
	{
		schools.GET("", controller.GetSchools)
		schools.GET("/:id", controller.GetSchool)
		schools.GET("/name", controller.GetSchoolByName)
		schools.POST("", controller.CreateSchool)
		schools.PUT("/:id", controller.UpdateSchool)
		schools.DELETE("/:id", controller.DeleteSchool)
	}
}
