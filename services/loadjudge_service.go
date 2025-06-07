package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
	"share-cat/Nepeta-cataria/pkg/makejudge"
	"share-cat/Nepeta-cataria/repositories"
	"strings"
	"time"
)

type LoadJudgeService struct {
	judgeRepo     *repositories.JudgeRepository
	judgesBaseDir string
}

func NewLoadJudgeService(
	judgeRepo *repositories.JudgeRepository,

) *LoadJudgeService {
	workDir, _ := os.Getwd()
	judgesBaseDir := filepath.Join(workDir, "judges")

	return &LoadJudgeService{
		judgeRepo:     judgeRepo,
		judgesBaseDir: judgesBaseDir,
	}
}

type LoadJudgeRequest struct {
	ZipFilePath string `json:"zip_file_path"`
	CreatorID   int    `json:"creator_id"`
	JudgeName   string `json:"judge_name,omitempty"`
}

type LoadJudgeResponse struct {
	JudgeID       uint   `json:"judge_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermanentPath string `json:"permanent_path"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
}

func (s *LoadJudgeService) LoadJudge(request LoadJudgeRequest) (*LoadJudgeResponse, error) {
	log.Printf("開始載入評判題目: %s", request.ZipFilePath)

	if err := os.MkdirAll(s.judgesBaseDir, 0755); err != nil {
		return nil, fmt.Errorf("創建 judges 目錄失敗: %w", err)
	}

	judgeDirName := s.generateJudgeDirName(request.JudgeName)
	judgeDir := filepath.Join(s.judgesBaseDir, judgeDirName)

	judgeDir = s.ensureUniqueJudgeDir(judgeDir)

	if err := s.extractZipFile(request.ZipFilePath, judgeDir); err != nil {

		os.RemoveAll(judgeDir)
		return nil, fmt.Errorf("解壓縮檔案失敗: %w", err)
	}

	config, configDir, err := s.loadProblemConfig(judgeDir)
	if err != nil {

		os.RemoveAll(judgeDir)
		return nil, fmt.Errorf("載入問題配置失敗: %w", err)
	}

	for _, lang := range config.Languages {
		switch strings.ToLower(lang) {
		case loadjudge.Python:
			makejudge.GenerateMainSolutionFile(configDir, loadjudge.Python, config)
		case loadjudge.Cpp:
			makejudge.GenerateMainSolutionFile(configDir, loadjudge.Cpp, config)
		case loadjudge.JavaScript:
			makejudge.GenerateMainSolutionFile(configDir, loadjudge.JavaScript, config)
		default:
			log.Printf("警告: 不支援的語言 '%s'，將被忽略", lang)
			continue
		}
	}

	judge, err := s.createJudgeRecord(config, request.CreatorID, configDir)
	if err != nil {

		os.RemoveAll(judgeDir)
		return nil, fmt.Errorf("創建資料庫記錄失敗: %w", err)
	}

	response := &LoadJudgeResponse{
		JudgeID:       judge.ID,
		Name:          judge.Name,
		Description:   judge.Description,
		PermanentPath: configDir,
		Success:       true,
		Message:       "評判題目載入成功",
	}

	log.Printf("評判題目載入成功: ID=%d, Name=%s, Path=%s", judge.ID, judge.Name, configDir)
	return response, nil
}

func (s *LoadJudgeService) generateJudgeDirName(customName string) string {
	if customName != "" {

		safeName := strings.ReplaceAll(customName, "/", "_")
		safeName = strings.ReplaceAll(safeName, "\\", "_")
		safeName = strings.ReplaceAll(safeName, "..", "_")
		return safeName
	}

	return fmt.Sprintf("judge_%d", time.Now().Unix())
}

func (s *LoadJudgeService) ensureUniqueJudgeDir(baseDir string) string {
	originalDir := baseDir
	counter := 1

	for {
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			return baseDir
		}
		baseDir = fmt.Sprintf("%s_%d", originalDir, counter)
		counter++
	}
}

func (s *LoadJudgeService) extractZipFile(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打開 ZIP 檔案失敗: %w", err)
	}
	defer reader.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("創建目標目錄失敗: %w", err)
	}

	for _, file := range reader.File {
		if err := s.extractSingleFile(file, destDir); err != nil {
			return fmt.Errorf("解壓縮檔案 %s 失敗: %w", file.Name, err)
		}
	}

	return nil
}

func (s *LoadJudgeService) extractSingleFile(file *zip.File, destDir string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(destDir, file.Name)

	if !strings.HasPrefix(path, filepath.Clean(destDir)+string(os.PathSeparator)) {
		return fmt.Errorf("無效的檔案路徑: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, file.FileInfo().Mode())
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

func (s *LoadJudgeService) loadProblemConfig(extractDir string) (loadjudge.ProblemConfig, string, error) {

	configPath := filepath.Join(extractDir, "config.yaml")
	configDir := extractDir

	if !loadjudge.FileExists(configPath) {

		entries, err := os.ReadDir(extractDir)
		if err != nil {
			return loadjudge.ProblemConfig{}, "", fmt.Errorf("讀取目錄失敗: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				subConfigPath := filepath.Join(extractDir, entry.Name(), "config.yaml")
				if loadjudge.FileExists(subConfigPath) {
					configPath = subConfigPath
					configDir = filepath.Join(extractDir, entry.Name())
					break
				}
			}
		}
	}

	if !loadjudge.FileExists(configPath) {
		return loadjudge.ProblemConfig{}, "", fmt.Errorf("找不到 config.yaml 檔案")
	}

	config, err := loadjudge.LoadProblemConfig(configDir)
	if err != nil {
		return loadjudge.ProblemConfig{}, "", fmt.Errorf("載入配置失敗: %w", err)
	}

	return config, configDir, nil
}

func (s *LoadJudgeService) createJudgeRecord(config loadjudge.ProblemConfig, creatorID int, configDir string) (*models.Judge, error) {

	readme := ""
	readmePath := filepath.Join(configDir, config.ReadmeFile)
	if config.ReadmeFile != "" && loadjudge.FileExists(readmePath) {
		readmeBytes, err := os.ReadFile(readmePath)
		if err != nil {
			log.Printf("警告: 讀取 README 檔案失敗: %v", err)
		} else {
			readme = string(readmeBytes)
		}
	}

	settingJSON, err := makejudge.GenerateSettingDataJSON(configDir, config)
	if err != nil {
		return nil, fmt.Errorf("獲取JSON 失敗: %w", err)
	}

	fullSetting := map[string]interface{}{
		"config":         config,
		"problem_path":   configDir,
		"test_data_path": filepath.Join(configDir, loadjudge.TestDataDir),
		"checker_file":   config.Checker.File,
		"checker_lang":   config.Checker.Language,
		"fullsetting":    json.RawMessage(settingJSON),
		"created_at":     time.Now().Unix(),
	}

	fullSettingJSON, err := json.Marshal(fullSetting)
	if err != nil {
		return nil, fmt.Errorf("序列化完整設定失敗: %w", err)
	}

	level := 1
	switch strings.ToLower(config.Level) {
	case "easy", "簡單", "1":
		level = 1
	case "medium", "中等", "2":
		level = 2
	case "hard", "困難", "3":
		level = 3
	default:
		level = 1
	}

	judge := models.Judge{
		Name:        config.Name,
		Description: config.Description,
		Readme:      readme,
		CreatorID:   creatorID,
		Tags:        config.Tags,
		Langs:       config.Languages,
		Type:        config.Type,
		MemoryLimit: config.MemoryLimitMb,
		TimeLimit:   int(config.TimeLimitSec),
		CheckType:   config.CheckType,
		IsActive:    true,
		Difficulty:  level,
	}

	if err := judge.FullSetting.Set(fullSettingJSON); err != nil {
		return nil, fmt.Errorf("設定 FullSetting 失敗: %w", err)
	}

	createdJudge, err := s.judgeRepo.Create(judge)
	if err != nil {
		return nil, fmt.Errorf("儲存 Judge 記錄失敗: %w", err)
	}

	return &createdJudge, nil
}

func (s *LoadJudgeService) ValidateZipFile(zipPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("無法打開 ZIP 檔案: %w", err)
	}
	defer reader.Close()

	hasConfig := false
	hasTestData := false

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "config.yaml") {
			hasConfig = true
		}
		if strings.Contains(file.Name, "testdata.yaml") {
			hasTestData = true
		}
	}

	if !hasConfig {
		return fmt.Errorf("ZIP 檔案中缺少 config.yaml")
	}

	if !hasTestData {
		return fmt.Errorf("ZIP 檔案中缺少測試資料")
	}

	return nil
}

func (s *LoadJudgeService) GetJudgeDirectory(judgeID uint) (string, error) {

	judge, err := s.judgeRepo.GetByID(judgeID)
	if err != nil {
		return "", fmt.Errorf("找不到 Judge 記錄: %w", err)
	}

	var fullSetting map[string]interface{}
	if err := judge.FullSetting.AssignTo(&fullSetting); err != nil {
		return "", fmt.Errorf("解析 FullSetting 失敗: %w", err)
	}

	if problemPath, ok := fullSetting["problem_path"].(string); ok {
		return problemPath, nil
	}

	return "", fmt.Errorf("無法從 FullSetting 中找到問題路徑")
}

func (s *LoadJudgeService) GetJudgesBaseDir() string {
	return s.judgesBaseDir
}
