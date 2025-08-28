package routes

import (
	"share-cat/Nepeta-cataria/controllers"

	"github.com/gin-gonic/gin"
)

func SetupLessonRoutes(router *gin.RouterGroup, lr *controllers.LessonController) {
	lessons := router.Group("/lessons")
	{
		// 課程基本操作
		lessons.POST("", lr.CreateLesson)       // 創建課程
		lessons.GET("", lr.GetAllLessons)       // 獲取所有課程
		lessons.GET("/:id", lr.GetLessonByID)   // 獲取課程詳情
		lessons.PUT("/:id", lr.UpdateLesson)    // 更新課程
		lessons.DELETE("/:id", lr.DeleteLesson) // 刪除課程

		// 教師相關
		lessons.GET("/teacher/:id", lr.GetTeacherLessons) // 獲取教師課程

		// 學生相關
		lessons.GET("/student/:id", lr.GetStudentLessons) // 獲取學生課程
		lessons.POST("/:id/join", lr.JoinLesson)          // 加入課程
		lessons.DELETE("/:id/leave", lr.LeaveLesson)      // 退出課程
	}
}
