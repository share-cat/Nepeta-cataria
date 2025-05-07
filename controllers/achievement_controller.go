package controllers

import (
	"net/http"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AchievementController struct {
	achievementService *services.AchievementService
}

func NewAchievementController(achievementService *services.AchievementService) *AchievementController {
	return &AchievementController{
		achievementService: achievementService,
	}
}

// GetAchievements 獲取所有成就
// @Summary 獲取所有成就
// @Description 獲取所有成就的列表
// @Tags achievements
// @Accept json
// @Produce json
// @Success 200 {array} models.Achievement
// @Failure 500 {object} map[string]string
// @Router /achievements [get]
func (c *AchievementController) GetAchievements(ctx *gin.Context) {
	achievements, err := c.achievementService.GetAllAchievements()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, achievements)
}

// GetAchievement 根據ID獲取成就
// @Summary 獲取特定成就
// @Description 根據ID獲取特定成就
// @Tags achievements
// @Accept json
// @Produce json
// @Param id path int true "Achievement ID"
// @Success 200 {object} models.Achievement
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements/{id} [get]
func (c *AchievementController) GetAchievement(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	achievement, err := c.achievementService.GetAchievementByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, achievement)
}

// CreateAchievement 創建新成就
// @Summary 創建新成就
// @Description 創建新的成就
// @Tags achievements
// @Accept json
// @Produce json
// @Param achievement body models.Achievement true "Achievement"
// @Success 201 {object} models.Achievement
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements [post]
func (c *AchievementController) CreateAchievement(ctx *gin.Context) {
	var achievement models.Achievement
	if err := ctx.ShouldBindJSON(&achievement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAchievement, err := c.achievementService.CreateAchievement(achievement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, newAchievement)
}

// UpdateAchievement 更新成就
// @Summary 更新成就
// @Description 更新指定ID的成就
// @Tags achievements
// @Accept json
// @Produce json
// @Param achievement body models.Achievement true "Achievement"
// @Success 200 {object} models.Achievement
// @Failure 400,404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements [put]
func (c *AchievementController) UpdateAchievement(ctx *gin.Context) {

	var achievement models.Achievement
	if err := ctx.ShouldBindJSON(&achievement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.achievementService.UpdateAchievement(achievement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, achievement)
}

// DeleteAchievement 刪除成就
// @Summary 刪除成就
// @Description 根據ID刪除成就
// @Tags achievements
// @Accept json
// @Produce json
// @Param id path int true "Achievement ID"
// @Success 204
// @Failure 400,404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements/{id} [delete]
func (c *AchievementController) DeleteAchievement(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.achievementService.DeleteAchievement(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
