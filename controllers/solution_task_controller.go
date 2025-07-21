package controllers

import (
	"net/http"
	"strconv"

	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type SolutionTaskController struct {
	solutionTaskService *services.SolutionTaskService
}

func NewSolutionTaskController(solutionTaskService *services.SolutionTaskService) *SolutionTaskController {
	return &SolutionTaskController{
		solutionTaskService: solutionTaskService,
	}
}

// GetSolutionTasks 獲取解答的所有任務
// @Summary 獲取解答的所有任務
// @Description 獲取特定解答的所有任務
// @Tags solution-tasks
// @Accept json
// @Produce json
// @Param solution_id path string true "Solution ID"
// @Success 200 {array} models.SolutionTask
// @Failure 400,404 {object} map[string]string
// @Router /solutions/{solution_id}/tasks [get]
func (c *SolutionTaskController) GetSolutionTasks(ctx *gin.Context) {
	solutionID, err := strconv.ParseUint(ctx.Param("solution_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的解答ID",
		})
		return
	}

	tasks, err := c.solutionTaskService.GetTasksBySolutionID(uint(solutionID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// GetSolutionTask 獲取特定任務
// @Summary 獲取特定任務
// @Description 根據ID獲取特定任務
// @Tags solution-tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.SolutionTask
// @Failure 400,404 {object} map[string]string
// @Router /solution-tasks/{id} [get]
func (c *SolutionTaskController) GetSolutionTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的任務ID",
		})
		return
	}

	task, err := c.solutionTaskService.GetTaskByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "任務不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// GetSolutionSummary 獲取解答匯總資訊
// @Summary 獲取解答匯總資訊
// @Description 獲取解答的匯總統計資訊
// @Tags solution-tasks
// @Accept json
// @Produce json
// @Param solution_id path string true "Solution ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]string
// @Router /solutions/{solution_id}/summary [get]
func (c *SolutionTaskController) GetSolutionSummary(ctx *gin.Context) {
	solutionID, err := strconv.ParseUint(ctx.Param("solution_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的解答ID",
		})
		return
	}

	summary, err := c.solutionTaskService.GetSolutionSummary(uint(solutionID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// UpdateTaskStatus 更新任務狀態
// @Summary 更新任務狀態
// @Description 更新特定任務的狀態和結果
// @Tags solution-tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param status body object true "Task Status"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /solution-tasks/{id}/status [put]
func (c *SolutionTaskController) UpdateTaskStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的任務ID",
		})
		return
	}

	var status struct {
		Status        string `json:"status" binding:"required"`
		ExecutionTime int64  `json:"execution_time"`
		MemoryUsage   int64  `json:"memory_usage"`
		Score         int    `json:"score"`
	}

	if err := ctx.ShouldBindJSON(&status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.solutionTaskService.UpdateTaskStatus(
		uint(id),
		status.Status,
		status.ExecutionTime,
		status.MemoryUsage,
		status.Score,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "任務狀態更新成功",
	})
}

// DeleteTask 刪除任務
// @Summary 刪除任務
// @Description 刪除特定任務
// @Tags solution-tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /solution-tasks/{id} [delete]
func (c *SolutionTaskController) DeleteTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的任務ID",
		})
		return
	}

	err = c.solutionTaskService.DeleteTask(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "任務刪除成功",
	})
}
