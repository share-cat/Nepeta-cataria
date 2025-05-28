package controllers

import (
	"net/http"
	"strconv"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type SolutionController struct {
	solutionService *services.SolutionService
}

func NewSolutionController(solutionService *services.SolutionService) *SolutionController {
	return &SolutionController{
		solutionService: solutionService,
	}
}

// GetSolutions 獲取所有答案
// @Summary 獲取所有答案
// @Description 獲取所有答案的列表
// @Tags solutions
// @Accept json
// @Produce json
// @Success 200 {array} models.Solution
// @Router /solutions [get]
func (c *SolutionController) GetSolutions(ctx *gin.Context) {
	// 檢查是否要獲取最新的提交
	latest := ctx.Query("latest")
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	if limit <= 0 {
		limit = 10 // 默認限制
	}

	var solutions []models.Solution
	var err error

	if latest == "true" {
		solutions, err = c.solutionService.GetLatestSolutions(limit)
	} else {
		solutions, err = c.solutionService.GetAllSolutions()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, solutions)
}

// GetSolution 根據ID獲取答案
// @Summary 獲取特定答案
// @Description 根據ID獲取特定答案
// @Tags solutions
// @Accept json
// @Produce json
// @Param id path string true "Solution ID"
// @Success 200 {object} models.Solution
// @Failure 404 {object} map[string]string
// @Router /solutions/{id} [get]
func (c *SolutionController) GetSolution(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的答案ID",
		})
		return
	}

	solution, err := c.solutionService.GetSolutionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "答案不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, solution)
}

// SubmitSolution 提交新答案
// @Summary 提交新答案
// @Description 提交一個新答案
// @Tags solutions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param solution body models.Solution true "Solution Info"
// @Success 201 {object} models.Solution
// @Failure 400 {object} map[string]string
// @Router /solutions [post]
func (c *SolutionController) SubmitSolution(ctx *gin.Context) {
	var solution models.Solution

	if err := ctx.ShouldBindJSON(&solution); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 從身份驗證中獲取用戶 ID (如果已實現身份驗證)
	userID, exists := ctx.Get("userID")
	if exists {
		solution.UserID = userID.(int)
	}

	createdSolution, err := c.solutionService.SubmitSolution(solution)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdSolution)
}

// UpdateSolutionResult 更新答案結果
// @Summary 更新答案結果
// @Description 更新特定答案的評分結果
// @Tags solutions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Solution ID"
// @Param result body object true "Solution Result"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /solutions/{id}/result [put]
func (c *SolutionController) UpdateSolutionResult(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的答案ID",
		})
		return
	}

	var result struct {
		Status        string `json:"status" binding:"required"`
		ExecutionTime int    `json:"execution_time"`
		MemoryUsage   int    `json:"memory_usage"`
		Score         int    `json:"score"`
	}

	if err := ctx.ShouldBindJSON(&result); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.solutionService.UpdateSolutionResult(
		uint(id),
		result.Status,
		result.ExecutionTime,
		result.MemoryUsage,
		result.Score,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "答案結果更新成功",
	})
}

// DeleteSolution 刪除答案
// @Summary 刪除答案
// @Description 刪除特定答案
// @Tags solutions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Solution ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /solutions/{id} [delete]
func (c *SolutionController) DeleteSolution(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的答案ID",
		})
		return
	}

	err = c.solutionService.DeleteSolution(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "答案刪除成功",
	})
}

// GetUserSolutions 獲取特定用戶的所有答案
// @Summary 獲取特定用戶的所有答案
// @Description 獲取特定用戶提交的所有答案
// @Tags solutions
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} models.Solution
// @Failure 400,404 {object} map[string]string
// @Router /users/{id}/solutions [get]
func (c *SolutionController) GetUserSolutions(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的用戶ID",
		})
		return
	}

	solutions, err := c.solutionService.GetUserSolutions(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, solutions)
}

// GetJudgeSolutions 獲取特定題目的所有答案
// @Summary 獲取特定題目的所有答案
// @Description 獲取提交到特定題目的所有答案
// @Tags solutions
// @Accept json
// @Produce json
// @Param id path string true "Judge ID"
// @Success 200 {array} models.Solution
// @Failure 400,404 {object} map[string]string
// @Router /judges/{id}/solutions [get]
func (c *SolutionController) GetJudgeSolutions(ctx *gin.Context) {
	judgeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目ID",
		})
		return
	}

	solutions, err := c.solutionService.GetJudgeSolutions(uint(judgeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, solutions)
}

// GetBestSolutions 獲取評分最高的答案
// @Summary 獲取評分最高的答案
// @Description 獲取評分最高的答案列表
// @Tags solutions
// @Accept json
// @Produce json
// @Param judge_id query string false "Judge ID"
// @Param limit query int false "Limit"
// @Success 200 {array} models.Solution
// @Failure 400 {object} map[string]string
// @Router /solutions/best [get]
func (c *SolutionController) GetBestSolutions(ctx *gin.Context) {
	judgeIDStr := ctx.Query("judge_id")
	limitStr := ctx.Query("limit")

	var judgeID uint = 0
	var limit int = 10

	if judgeIDStr != "" {
		id, err := strconv.ParseUint(judgeIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "無效的題目ID",
			})
			return
		}
		judgeID = uint(id)
	}

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
	}

	solutions, err := c.solutionService.GetBestSolutions(judgeID, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, solutions)
}
