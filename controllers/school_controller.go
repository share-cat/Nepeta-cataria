package controllers

import (
	"net/http"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SchoolController struct {
	schoolService *services.SchoolService
}

func NewSchoolController(schoolService *services.SchoolService) *SchoolController {
	return &SchoolController{
		schoolService: schoolService,
	}
}

// GetSchools 獲取所有學校
// @Summary 獲取所有學校
// @Description 獲取所有學校的列表
// @Tags schools
// @Accept json
// @Produce json
// @Success 200 {array} models.School
// @Failure 500 {object} map[string]string
// @Router /schools [get]
func (c *SchoolController) GetSchools(ctx *gin.Context) {
	schools, err := c.schoolService.GetAllSchools()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, schools)
}

// GetSchool 根據ID獲取學校
// @Summary 獲取特定學校
// @Description 根據ID獲取特定學校
// @Tags schools
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} models.School
// @Failure 400,404 {object} map[string]string
// @Router /schools/{id} [get]
func (c *SchoolController) GetSchool(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	school, err := c.schoolService.GetSchoolByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "school not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, school)
}

// GetSchoolByName 根據名稱獲取學校
// @Summary 根據名稱獲取學校
// @Description 根據名稱獲取學校
// @Tags schools
// @Accept json
// @Produce json
// @Param name query string true "School Name"
// @Success 200 {object} models.School
// @Failure 400,404 {object} map[string]string
// @Router /schools/name [get]
func (c *SchoolController) GetSchoolByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "name is required",
		})
		return
	}

	school, err := c.schoolService.GetSchoolByName(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "school not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, school)
}

// CreateSchool 創建學校
// @Summary 創建學校
// @Description 創建新的學校
// @Tags schools
// @Accept json
// @Produce json
// @Param school body models.School true "School"
// @Success 201 {object} models.School
// @Failure 400 {object} map[string]string
// @Router /schools [post]
func (c *SchoolController) CreateSchool(ctx *gin.Context) {
	var school models.School

	if err := ctx.ShouldBindJSON(&school); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdSchool, err := c.schoolService.CreateSchool(school)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdSchool)
}

// UpdateSchool 更新學校
// @Summary 更新學校
// @Description 更新學校的詳細信息
// @Tags schools
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Param school body models.School true "School"
// @Success 200 {object} models.School
// @Failure 400,404 {object} map[string]string
// @Router /schools/{id} [put]
func (c *SchoolController) UpdateSchool(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var school models.School
	if err := ctx.ShouldBindJSON(&school); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	school.ID = uint(id)

	if err := c.schoolService.UpdateSchool(school); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "school not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, school)
}

// DeleteSchool 刪除學校
// @Summary 刪除學校
// @Description 刪除學校
// @Tags schools
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Success 204
// @Failure 400,404 {object} map[string]string
// @Router /schools/{id} [delete]
func (c *SchoolController) DeleteSchool(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	if err := c.schoolService.DeleteSchool(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "school not found",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
