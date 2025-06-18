package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type JudgeController struct {
	judgeService *services.JudgeService
}

func NewJudgeController(judgeService *services.JudgeService) *JudgeController {
	return &JudgeController{
		judgeService: judgeService,
	}
}

// PaginatedJudgeResponse
type PaginatedJudgeResponse struct {
	Data       []models.Judge `json:"data"`
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
	Total      int            `json:"total"`
}

// GetJudges 獲取題目列表（支援搜尋、類型篩選、分頁）
// @Summary 獲取題目列表
// @Description 獲取題目列表，支援搜尋、類型篩選、分頁和活躍狀態篩選
// @Tags judges
// @Accept json
// @Produce json
// @Param search query string false "搜尋關鍵字（支援名稱和描述搜尋）"
// @Param type query string false "題目類型篩選"
// @Param page query int false "頁碼（預設為1）" default(1)
// @Param active query string false "只顯示活躍題目（true/false）"
// @Success 200 {object} PaginatedJudgeResponse
// @Failure 400,500 {object} map[string]string
// @Router /judges [get]
func (c *JudgeController) GetJudges(ctx *gin.Context) {

	search := ctx.Query("search")
	judgeType := ctx.Query("type")
	pageStr := ctx.DefaultQuery("page", "1")
	activeOnly := ctx.Query("active") == "true"

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的頁碼",
		})
		return
	}

	if search == "" && judgeType == "" && page == 1 {
		var judges []models.Judge
		var err error

		if activeOnly {
			judges, err = c.judgeService.GetActiveJudges()
		} else {
			judges, err = c.judgeService.GetAllJudges()
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, PaginatedJudgeResponse{
			Data:       judges,
			Page:       1,
			TotalPages: 1,
			Total:      len(judges),
		})
		return
	}

	judges, total, err := c.judgeService.GetJudgesWithFilters(search, judgeType, page, activeOnly)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	pageSize := 20
	totalPages := (total + pageSize - 1) / pageSize

	ctx.JSON(http.StatusOK, PaginatedJudgeResponse{
		Data:       judges,
		Page:       page,
		TotalPages: totalPages,
		Total:      total,
	})
}

// GetJudge 根據ID獲取題目
// @Summary 獲取特定題目
// @Description 根據ID獲取特定題目
// @Tags judges
// @Accept json
// @Produce json
// @Param id path string true "Judge ID"
// @Success 200 {object} models.Judge
// @Failure 404 {object} map[string]string
// @Router /judges/{id} [get]
func (c *JudgeController) GetJudge(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目ID",
		})
		return
	}

	judge, err := c.judgeService.GetJudgeByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "題目不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, judge)
}

// GetJudgeCode 獲取題目代碼
// @Summary 獲取題目代碼
// @Description 根據題目ID和程式語言獲取題目代碼
// @Tags judges
// @Accept json
// @Produce json
// @Param id path string true "Judge ID"
// @Param language path string true "Programming Language (e.g., python, cpp)"
// @Success 200 {string} string "Code content"
// @Failure 400,404 {object} map[string]string
// @Router /judges/{id}/code/{language} [get]
func (c *JudgeController) GetJudgeCode(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目ID",
		})
		return
	}

	languageStr := ctx.Param("language")
	code, err := c.judgeService.GetJudgeCode(uint(id), languageStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "題目代碼不存在或不支援的程式語言",
		})
		return
	}

	ctx.String(http.StatusOK, code)
}

// CreateJudge 創建新題目
// @Summary 創建新題目
// @Description 創建一個新題目
// @Tags judges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param judge body models.Judge true "Judge Info"
// @Success 201 {object} models.Judge
// @Failure 400 {object} map[string]string
// @Router /judges [post]
func (c *JudgeController) CreateJudge(ctx *gin.Context) {
	var judge models.Judge

	if err := ctx.ShouldBindJSON(&judge); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 從 JWT 令牌或會話中獲取當前用戶 ID
	// 注意：這裡假設已經有一個身份驗證中間件設置了用戶ID
	// 如果您沒有使用身份驗證，請相應地調整此邏輯
	userID, exists := ctx.Get("userID")
	if exists {
		judge.CreatorID = userID.(int)
	}

	createdJudge, err := c.judgeService.CreateJudge(judge)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdJudge)
}

// UpdateJudge 更新題目
// @Summary 更新題目
// @Description 更新特定題目的信息
// @Tags judges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Judge ID"
// @Param judge body models.Judge true "Judge Info"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /judges/{id} [put]
func (c *JudgeController) UpdateJudge(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目ID",
		})
		return
	}

	var judgeUpdate models.Judge
	if err := ctx.ShouldBindJSON(&judgeUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.judgeService.UpdateJudge(uint(id), judgeUpdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "題目更新成功",
	})
}

// DeleteJudge 刪除題目
// @Summary 刪除題目
// @Description 刪除特定題目
// @Tags judges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Judge ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /judges/{id} [delete]
func (c *JudgeController) DeleteJudge(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目ID",
		})
		return
	}

	err = c.judgeService.DeleteJudge(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "題目刪除成功",
	})
}

// GetJudgesByCreator 獲取特定用戶創建的題目
// @Summary 獲取特定用戶創建的題目
// @Description 獲取特定用戶創建的所有題目
// @Tags judges
// @Accept json
// @Produce json
// @Param id path string true "Creator ID"
// @Success 200 {array} models.Judge
// @Failure 400,404 {object} map[string]string
// @Router /users/{id}/judges [get]
func (c *JudgeController) GetJudgesByCreator(ctx *gin.Context) {
	creatorID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	judges, err := c.judgeService.GetJudgesByCreator(uint(creatorID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, judges)
}

// GetJudgesByTags 根據標籤獲取題目
// @Summary 根據標籤獲取題目
// @Description 根據標籤獲取符合條件的題目
// @Tags judges
// @Accept json
// @Produce json
// @Param tags query string true "Tags (comma separated)"
// @Success 200 {array} models.Judge
// @Failure 400 {object} map[string]string
// @Router /judges/tags [get]
func (c *JudgeController) GetJudgesByTags(ctx *gin.Context) {
	tagsStr := ctx.Query("tags")
	if tagsStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "請提供至少一個標籤",
		})
		return
	}

	tags := strings.Split(tagsStr, ",")

	judges, err := c.judgeService.GetJudgesByTags(tags)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, judges)
}

// GetJudgesByDifficulty 根據難度獲取題目
// @Summary 根據難度獲取題目
// @Description 根據難度獲取符合條件的題目
// @Tags judges
// @Accept json
// @Produce json
// @Param difficulty path int true "Difficulty Level"
// @Success 200 {array} models.Judge
// @Failure 400 {object} map[string]string
// @Router /judges/difficulty/{difficulty} [get]
func (c *JudgeController) GetJudgesByDifficulty(ctx *gin.Context) {
	diffStr := ctx.Param("difficulty")
	difficulty, err := strconv.Atoi(diffStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的難度值",
		})
		return
	}

	judges, err := c.judgeService.GetJudgesByDifficulty(difficulty)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, judges)
}
