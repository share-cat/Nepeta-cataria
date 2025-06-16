package controllers

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/pkg/catwheel"
	"share-cat/Nepeta-cataria/pkg/makejudge"
	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type SolutionController struct {
	solutionService *services.SolutionService
	judgeService    *services.JudgeService
	langService     *services.SupportLangService
}

func NewSolutionController(solutionService *services.SolutionService, judgeService *services.JudgeService, langService *services.SupportLangService) *SolutionController {
	return &SolutionController{
		solutionService: solutionService,
		judgeService:    judgeService,
		langService:     langService,
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

type SolutionSubmission struct {
	Code     string `json:"code" binding:"required"`
	JudgeID  int    `json:"judge_id" binding:"required"`
	Language string `json:"language" binding:"required"`
}

// SubmitSolution 提交新答案
// @Summary 提交新答案
// @Description 提交一個新答案
// @Tags solutions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param solution body SolutionSubmission true "Solution Info"
// @Success 201 {object} models.Solution
// @Failure 400 {object} map[string]string
// @Router /solutions [post]
func (c *SolutionController) SubmitSolution(ctx *gin.Context) {
	var submission SolutionSubmission
	var solution models.Solution

	if err := ctx.ShouldBindJSON(&submission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")

	if !exists {
		// 應該不太可能:P
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授權的請求",
		})
		return
	}

	judgedata, err := c.judgeService.GetJudgeByID(uint(submission.JudgeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "題目不存在",
		})
		return
	}

	language, err := c.langService.GetSupportLangByName(submission.Language)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{

			"error": "無效的語言: " + submission.Language,
		})
		return
	}
	if language.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的語言: " + submission.Language,
		})
		return
	}

	solution = models.Solution{
		Code:       submission.Code,
		JudgeID:    submission.JudgeID,
		UserID:     userID.(int),
		LanguageID: int(language.ID),
	}

	createdSolution, err := c.solutionService.SubmitSolution(solution)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	client := catwheel.NewClient("http://localhost:3200")
	judgedatasetting := judgedata.GetSetting()

	//TODO: 暫時這樣處理 一定要重構!!!
	var solutionex *catwheel.ExecuteResponse
	var totalExecutionTime int64 = 0
	var totalMemoryUsed int64 = 0
	var overallSuccess bool = true
	var overallStatus string = "COMPLETED"

	unb64code, err := base64.StdEncoding.DecodeString(submission.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "轉換b64失敗: " + err.Error(),
		})
		return
	}

	for _, testcase := range judgedatasetting.TestCases {
		executeCode, err := makejudge.GenerateExecutableCodeWithSettingData(judgedatasetting, testcase, string(unb64code), submission.Language)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "生成可執行代碼失敗: " + err.Error(), //TODO: 以後記得刪除
			})
			return
		}

		var langOption string
		switch submission.Language {
		case "python":
			langOption = catwheel.LanguagePython
		case "cpp", "c++":
			langOption = catwheel.LanguageCpp
		case "javascript", "nodejs":
			langOption = catwheel.LanguageJavaScript
		default:
			langOption = catwheel.LanguagePython
		}

		options := catwheel.DefaultExecuteOptions().
			WithLanguage(langOption).
			WithTimeout(10).
			WithMemoryLimit(128 * 1024 * 1024)

		solution_res, err := client.ExecuteWithOptions(executeCode, options)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if solutionex == nil {
			solutionex = solution_res
		}

		totalExecutionTime += solution_res.Result.ExecutionTime
		totalMemoryUsed += solution_res.Result.MemoryUsed

		if !solution_res.Success {
			overallSuccess = false
		}

		switch solution_res.Result.Status {
		case "Failed", "Timeout", "MemoryLimit", "RuntimeError":
			overallSuccess = false
			overallStatus = solution_res.Result.Status
		}

		if solution_res.Result.ExitCode != 0 {
			overallSuccess = false
			if overallStatus == "COMPLETED" {
				overallStatus = "RUNTIME_ERROR"
			}
		}

		if solution_res.Result.GetStdout() != testcase.Output {
			overallSuccess = false
			overallStatus = "Answer Error"
		}

		solutionex = solution_res
	}

	if overallSuccess {
		overallStatus = "ACCEPTED"
	}

	if solutionex == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "執行結果為空",
		})
		return
	}

	aggregatedResult := catwheel.ExecuteResponse{
		ExecutionID: solutionex.ExecutionID,
		Success:     overallSuccess,
		Result: catwheel.ExecutionResult{
			ID:            solutionex.Result.ID,
			Status:        overallStatus,
			Stdout:        solutionex.Result.Stdout,
			Stderr:        solutionex.Result.Stderr,
			ExitCode:      solutionex.Result.ExitCode,
			ExecutionTime: totalExecutionTime,
			CPUTime:       solutionex.Result.CPUTime,
			MemoryUsed:    totalMemoryUsed,
			ErrorMessage:  solutionex.Result.ErrorMessage,
			CreatedAt:     solutionex.Result.CreatedAt,
			CompletedAt:   solutionex.Result.CompletedAt,
		},
	}

	updateSolution, err := c.solutionService.UpdateSolutionFromExecution(aggregatedResult, createdSolution.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新答案結果失敗: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, updateSolution)
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
