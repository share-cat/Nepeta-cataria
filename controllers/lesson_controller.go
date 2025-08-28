package controllers

import (
	"net/http"
	"strconv"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type LessonController struct {
	lessonService *services.LessonService
}

func NewLessonController(lessonService *services.LessonService) *LessonController {
	return &LessonController{
		lessonService: lessonService,
	}
}

// CreateLesson 創建課程
// @Summary 創建課程
// @Description 教師創建新課程
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lesson body models.Lesson true "課程資訊"
// @Success 201 {object} map[string]interface{} "data: models.Lesson"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Router /lessons [post]
func (lr *LessonController) CreateLesson(c *gin.Context) {
	var lesson models.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 從 JWT 或 session 中獲取當前用戶 ID
	teacherID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := lr.lessonService.CreateLesson(&lesson, teacherID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": lesson})
}

// GetLessonByID 獲取課程詳情
// @Summary 獲取課程詳情
// @Description 根據課程ID獲取課程詳細資訊
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "課程ID"
// @Success 200 {object} map[string]interface{} "data: models.Lesson"
// @Failure 400 {object} map[string]string "error: invalid lesson id"
// @Failure 404 {object} map[string]string "error: 課程不存在"
// @Router /lessons/{id} [get]
func (lr *LessonController) GetLessonByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	lesson, err := lr.lessonService.GetLessonByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lesson})
}

// GetTeacherLessons 獲取教師課程
// @Summary 獲取教師課程
// @Description 獲取指定教師的所有課程
// @Tags lessons
// @Accept json
// @Produce json
// @Param teacherId path int true "教師ID"
// @Success 200 {object} map[string]interface{} "data: []models.Lesson"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Router /lessons/teacher/{id} [get]
func (lr *LessonController) GetTeacherLessons(c *gin.Context) {
	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher id"})
		return
	}

	lessons, err := lr.lessonService.GetTeacherLessons(uint(teacherID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lessons})
}

// GetStudentLessons 獲取學生課程
// @Summary 獲取學生課程
// @Description 獲取指定學生加入的所有課程
// @Tags lessons
// @Accept json
// @Produce json
// @Param studentId path int true "學生ID"
// @Success 200 {object} map[string]interface{} "data: []models.Lesson"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Router /lessons/student/{id} [get]
func (lr *LessonController) GetStudentLessons(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	lessons, err := lr.lessonService.GetStudentLessons(uint(studentID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lessons})
}

// UpdateLesson 更新課程
// @Summary 更新課程
// @Description 更新課程資訊（僅課程創建者可更新）
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "課程ID"
// @Param lesson body models.Lesson true "課程資訊"
// @Success 200 {object} map[string]interface{} "data: models.Lesson"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Router /lessons/{id} [put]
func (lr *LessonController) UpdateLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var lesson models.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson.ID = uint(id)

	teacherID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := lr.lessonService.UpdateLesson(&lesson, teacherID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lesson})
}

// DeleteLesson 刪除課程
// @Summary 刪除課程
// @Description 刪除課程（僅課程創建者可刪除）
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "課程ID"
// @Success 200 {object} map[string]string "message: lesson deleted successfully"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Router /lessons/{id} [delete]
func (lr *LessonController) DeleteLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	teacherID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := lr.lessonService.DeleteLesson(uint(id), teacherID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lesson deleted successfully"})
}

// JoinLesson 加入課程
// @Summary 學生加入課程
// @Description 學生加入指定課程
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "課程ID"
// @Success 200 {object} map[string]string "message: joined lesson successfully"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Router /lessons/{id}/join [post]
func (lr *LessonController) JoinLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	studentID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := lr.lessonService.JoinLesson(uint(id), studentID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined lesson successfully"})
}

// LeaveLesson 退出課程
// @Summary 學生退出課程
// @Description 學生退出指定課程
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "課程ID"
// @Success 200 {object} map[string]string "message: left lesson successfully"
// @Failure 400 {object} map[string]string "error: 錯誤訊息"
// @Failure 401 {object} map[string]string "error: unauthorized"
// @Router /lessons/{id}/leave [delete]
func (lr *LessonController) LeaveLesson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	studentID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := lr.lessonService.LeaveLesson(uint(id), studentID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left lesson successfully"})
}

// GetAllLessons 獲取所有課程
// @Summary 獲取所有課程
// @Description 分頁獲取所有課程列表
// @Tags lessons
// @Accept json
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(10)
// @Success 200 {object} map[string]interface{} "data: []models.Lesson, total: int, page: int, page_size: int"
// @Failure 500 {object} map[string]string "error: 內部伺服器錯誤"
// @Router /lessons [get]
func (lr *LessonController) GetAllLessons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	lessons, total, err := lr.lessonService.GetAllLessons(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      lessons,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
