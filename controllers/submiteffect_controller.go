package controllers

import (
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubmitEffectController struct {
	submitEffectService *services.SubmitEffectService
}

func NewSubmitEffectController(submitEffectService *services.SubmitEffectService) *SubmitEffectController {
	return &SubmitEffectController{
		submitEffectService: submitEffectService,
	}
}

// GetSubmitEffects 獲取所有提交效果
// @Summary 獲取所有提交效果
// @Description 獲取所有提交效果的列表
// @Tags submit_effects
// @Accept json
// @Produce json
// @Success 200 {array} models.SubmitEffect
// @Failure 500 {object} map[string]string
// @Router /submit_effects [get]
func (c *SubmitEffectController) GetSubmitEffects(ctx *gin.Context) {
	submitEffects, err := c.submitEffectService.GetAllSubmitEffects()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, submitEffects)
}

// GetSubmitEffect 根據ID獲取提交效果
// @Summary 獲取特定提交效果
// @Description 根據ID獲取特定提交效果
// @Tags submit_effects
// @Accept json
// @Produce json
// @Param id path int true "SubmitEffect ID"
// @Success 200 {object} models.SubmitEffect
// @Failure 400,404 {object} map[string]string
// @Router /submit_effects/{id} [get]
func (c *SubmitEffectController) GetSubmitEffect(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid ID",
		})
		return
	}

	submitEffect, err := c.submitEffectService.GetSubmitEffectByID(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": "submit effect not found",
		})
		return
	}

	ctx.JSON(200, submitEffect)
}

// CreateSubmitEffect 創建提交效果
// @Summary 創建提交效果
// @Description 創建新的提交效果
// @Tags submit_effects
// @Accept json
// @Produce json
// @Param submitEffect body models.SubmitEffect true "SubmitEffect"
// @Success 201 {object} models.SubmitEffect
// @Failure 400 {object} map[string]string
// @Router /submit_effects [post]
func (c *SubmitEffectController) CreateSubmitEffect(ctx *gin.Context) {
	var submitEffect models.SubmitEffect

	if err := ctx.ShouldBindJSON(&submitEffect); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdSubmitEffect, err := c.submitEffectService.CreateSubmitEffect(submitEffect)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(201, createdSubmitEffect)
}

// UpdateSubmitEffect 更新提交效果
// @Summary 更新提交效果
// @Description 更新提交效果的詳細信息
// @Tags submit_effects
// @Accept json
// @Produce json
// @Param id path int true "SubmitEffect ID"
// @Param submitEffect body models.SubmitEffect true "SubmitEffect"
// @Success 200 {object} models.SubmitEffect
// @Failure 400,404 {object} map[string]string
// @Router /submit_effects/{id} [put]
func (c *SubmitEffectController) UpdateSubmitEffect(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var submitEffect models.SubmitEffect
	if err := ctx.ShouldBindJSON(&submitEffect); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	submitEffect.ID = uint(id)

	if err := c.submitEffectService.UpdateSubmitEffect(submitEffect); err != nil {
		ctx.JSON(404, gin.H{
			"error": "submit effect not found",
		})
		return
	}

	ctx.JSON(200, submitEffect)
}

// DeleteSubmitEffect 刪除提交效果
// @Summary 刪除提交效果
// @Description 根據ID刪除提交效果
// @Tags submit_effects
// @Accept json
// @Produce json
// @Param id path int true "SubmitEffect ID"
// @Success 204
// @Failure 400,404 {object} map[string]string
// @Router /submit_effects/{id} [delete]
func (c *SubmitEffectController) DeleteSubmitEffect(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid ID",
		})
		return
	}

	if err := c.submitEffectService.DeleteSubmitEffect(uint(id)); err != nil {
		ctx.JSON(404, gin.H{
			"error": "submit effect not found",
		})
		return
	}

	ctx.Status(204)
}
