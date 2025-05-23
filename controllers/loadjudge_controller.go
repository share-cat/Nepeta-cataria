package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"share-cat/Nepeta-cataria/services"

	"github.com/gin-gonic/gin"
)

type LoadJudgeController struct {
	loadJudgeService *services.LoadJudgeService
}

func NewLoadJudgeController(loadJudgeService *services.LoadJudgeService) *LoadJudgeController {
	return &LoadJudgeController{
		loadJudgeService: loadJudgeService,
	}
}

// UploadAndLoadJudge 上傳並載入題目
// @Summary 上傳並載入題目
// @Description 上傳 ZIP 檔案並載入為題目
// @Tags loadjudge
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Judge ZIP file"
// @Param judge_name formData string false "Custom judge name"
// @Success 200 {object} services.LoadJudgeResponse
// @Failure 400,500 {object} map[string]string
// @Router /load-judge/upload [post]
func (c *LoadJudgeController) UploadAndLoadJudge(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無法獲取上傳檔案",
		})
		return
	}

	if filepath.Ext(file.Filename) != ".zip" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "只支援 ZIP 檔案格式",
		})
		return
	}

	tempDir := filepath.Join(os.TempDir(), "judge_uploads")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "創建臨時目錄失敗",
		})
		return
	}

	timestamp := time.Now().Unix()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("judge_%d_%s", timestamp, file.Filename))

	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "儲存檔案失敗",
		})
		return
	}

	defer os.Remove(tempFilePath)

	if err := c.loadJudgeService.ValidateZipFile(tempFilePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("ZIP 檔案驗證失敗: %v", err),
		})
		return
	}

	var creatorID int
	if userID, exists := ctx.Get("userID"); exists {
		creatorID = userID.(int)
	} else {

		creatorIDStr := ctx.PostForm("creator_id")
		if creatorIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "缺少創建者 ID",
			})
			return
		}

		id, err := strconv.Atoi(creatorIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "無效的創建者 ID",
			})
			return
		}
		creatorID = id
	}

	judgeName := ctx.PostForm("judge_name")

	request := services.LoadJudgeRequest{
		ZipFilePath: tempFilePath,
		CreatorID:   creatorID,
		JudgeName:   judgeName,
	}

	response, err := c.loadJudgeService.LoadJudge(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("載入題目失敗: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// LoadJudgeFromPath 從指定路徑載入題目
// @Summary 從指定路徑載入題目
// @Description 從伺服器上指定的 ZIP 檔案路徑載入題目
// @Tags loadjudge
// @Accept json
// @Produce json
// @Param request body services.LoadJudgeRequest true "Load Judge Request"
// @Success 200 {object} services.LoadJudgeResponse
// @Failure 400,500 {object} map[string]string
// @Router /load-judge/path [post]
func (c *LoadJudgeController) LoadJudgeFromPath(ctx *gin.Context) {
	var request services.LoadJudgeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("請求參數錯誤: %v", err),
		})
		return
	}

	if request.ZipFilePath == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ZIP 檔案路徑不能為空",
		})
		return
	}

	if _, err := os.Stat(request.ZipFilePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "指定的 ZIP 檔案不存在",
		})
		return
	}

	if err := c.loadJudgeService.ValidateZipFile(request.ZipFilePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("ZIP 檔案驗證失敗: %v", err),
		})
		return
	}

	response, err := c.loadJudgeService.LoadJudge(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("載入題目失敗: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ValidateJudgeZip 驗證題目 ZIP 檔案
// @Summary 驗證題目 ZIP 檔案
// @Description 驗證上傳的 ZIP 檔案是否符合題目格式
// @Tags loadjudge
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Judge ZIP file"
// @Success 200 {object} map[string]string
// @Failure 400,500 {object} map[string]string
// @Router /load-judge/validate [post]
func (c *LoadJudgeController) ValidateJudgeZip(ctx *gin.Context) {
	// 獲取上傳的檔案
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無法獲取上傳檔案",
		})
		return
	}

	if filepath.Ext(file.Filename) != ".zip" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "只支援 ZIP 檔案格式",
		})
		return
	}

	tempDir := filepath.Join(os.TempDir(), "judge_validate")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "創建臨時目錄失敗",
		})
		return
	}

	timestamp := time.Now().Unix()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("validate_%d_%s", timestamp, file.Filename))

	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "儲存檔案失敗",
		})
		return
	}

	defer os.Remove(tempFilePath)

	if err := c.loadJudgeService.ValidateZipFile(tempFilePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   fmt.Sprintf("驗證失敗: %v", err),
			"valid":   false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "ZIP 檔案格式正確",
	})
}

// GetJudgeDirectory 獲取題目目錄路徑
// @Summary 獲取題目目錄路徑
// @Description 根據題目 ID 獲取其在伺服器上的目錄路徑
// @Tags loadjudge
// @Accept json
// @Produce json
// @Param id path string true "Judge ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /load-judge/directory/{id} [get]
func (c *LoadJudgeController) GetJudgeDirectory(ctx *gin.Context) {
	judgeIDStr := ctx.Param("id")
	judgeID, err := strconv.ParseUint(judgeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目 ID",
		})
		return
	}

	directory, err := c.loadJudgeService.GetJudgeDirectory(uint(judgeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("獲取目錄路徑失敗: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"judge_id":  judgeID,
		"directory": directory,
		"exists":    true,
	})
}

// DownloadJudgeFile 下載題目相關檔案
// @Summary 下載題目相關檔案
// @Description 下載指定題目的相關檔案
// @Tags loadjudge
// @Accept json
// @Produce application/octet-stream
// @Param id path string true "Judge ID"
// @Param file_path query string true "File path relative to judge directory"
// @Success 200 {file} binary
// @Failure 400,404 {object} map[string]string
// @Router /load-judge/download/{id} [get]
func (c *LoadJudgeController) DownloadJudgeFile(ctx *gin.Context) {
	judgeIDStr := ctx.Param("id")
	judgeID, err := strconv.ParseUint(judgeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題目 ID",
		})
		return
	}

	filePath := ctx.Query("file_path")
	if filePath == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少檔案路徑參數",
		})
		return
	}

	judgeDirectory, err := c.loadJudgeService.GetJudgeDirectory(uint(judgeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("獲取題目目錄失敗: %v", err),
		})
		return
	}

	fullFilePath := filepath.Join(judgeDirectory, filePath)

	if !filepath.HasPrefix(fullFilePath, judgeDirectory) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "不允許存取目錄外的檔案",
		})
		return
	}

	if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "檔案不存在",
		})
		return
	}

	file, err := os.Open(fullFilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "開啟檔案失敗",
		})
		return
	}
	defer file.Close()

	fileName := filepath.Base(fullFilePath)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "application/octet-stream")

	_, err = io.Copy(ctx.Writer, file)
	if err != nil {

		ctx.Header("X-Error", "檔案傳輸失敗")
	}
}

// GetJudgesBaseDirectory 獲取題目基礎目錄
// @Summary 獲取題目基礎目錄
// @Description 獲取所有題目的基礎存儲目錄
// @Tags loadjudge
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /load-judge/base-directory [get]
func (c *LoadJudgeController) GetJudgesBaseDirectory(ctx *gin.Context) {
	baseDir := c.loadJudgeService.GetJudgesBaseDir()

	ctx.JSON(http.StatusOK, gin.H{
		"base_directory": baseDir,
		"exists":         true,
	})
}
