package controllers

import (
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SocialController struct {
	schoolService *services.SocialService
}

func NewSocialController(schoolService *services.SocialService) *SocialController {
	return &SocialController{
		schoolService: schoolService,
	}
}

// GetSocials 獲取所有社交欄位
// @Summary 獲取所有社交欄位
// @Description 獲取所有社交欄位的列表
// @Tags socials
// @Accept json
// @Produce json
// @Success 200 {array} models.Social
// @Failure 500 {object} map[string]string
// @Router /socials [get]
func (c *SocialController) GetSocials(ctx *gin.Context) {
	socials, err := c.schoolService.GetAllSocials()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, socials)
}

// GetSocial 根據ID獲取社交欄位
// @Summary 獲取特定社交欄位
// @Description 根據ID獲取特定社交欄位
// @Tags socials
// @Accept json
// @Produce json
// @Param id path int true "Social ID"
// @Success 200 {object} models.Social
// @Failure 400,404 {object} map[string]string
// @Router /socials/{id} [get]
func (c *SocialController) GetSocial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid ID",
		})
		return
	}

	social, err := c.schoolService.GetSocialByID(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": "social not found",
		})
		return
	}

	ctx.JSON(200, social)
}

// CreateSocial 創建社交欄位
// @Summary 創建社交欄位
// @Description 創建新的社交欄位
// @Tags socials
// @Accept json
// @Produce json
// @Param social body models.Social true "Social"
// @Success 201 {object} models.Social
// @Failure 400 {object} map[string]string
// @Router /socials [post]
func (c *SocialController) CreateSocial(ctx *gin.Context) {
	var social models.Social

	if err := ctx.ShouldBindJSON(&social); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdSocial, err := c.schoolService.CreateSocial(social)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(201, createdSocial)
}

// UpdateSocial 更新社交欄位
// @Summary 更新社交欄位
// @Description 更新社交欄位的詳細信息
// @Tags socials
// @Accept json
// @Produce json
// @Param id path int true "Social ID"
// @Param social body models.Social true "Social"
// @Success 200 {object} models.Social
// @Failure 400,404 {object} map[string]string
// @Router /socials/{id} [put]
func (c *SocialController) UpdateSocial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var social models.Social
	if err := ctx.ShouldBindJSON(&social); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	social.ID = uint(id)

	if err := c.schoolService.UpdateSocial(social); err != nil {
		ctx.JSON(404, gin.H{
			"error": "social not found",
		})
		return
	}

	ctx.JSON(200, social)
}

// DeleteSocial 刪除社交欄位
// @Summary 刪除社交欄位
// @Description 根據ID刪除社交欄位
// @Tags socials
// @Accept json
// @Produce json
// @Param id path int true "Social ID"
// @Success 204
// @Failure 400,404 {object} map[string]string
// @Router /socials/{id} [delete]
func (c *SocialController) DeleteSocial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid ID",
		})
		return
	}

	if err := c.schoolService.DeleteSocial(uint(id)); err != nil {
		ctx.JSON(404, gin.H{
			"error": "social not found",
		})
		return
	}

	ctx.JSON(204, nil)
}
