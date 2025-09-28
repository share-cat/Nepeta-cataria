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
	"share-cat/Nepeta-cataria/pkg/language"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
	"share-cat/Nepeta-cataria/pkg/makejudge"
	"share-cat/Nepeta-cataria/repositories"
	"share-cat/Nepeta-cataria/utils"
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
	ForceUpdate bool   `json:"force_update,omitempty"`
}

type LoadJudgeFromDirectoryRequest struct {
	DirectoryName string `json:"directory_name"`
	CreatorID     int    `json:"creator_id"`
	ForceUpdate   bool   `json:"force_update,omitempty"`
}

type LoadJudgeResponse struct {
	JudgeID       uint   `json:"judge_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermanentPath string `json:"permanent_path"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	IsUpdate      bool   `json:"is_update"`
	ExistingID    uint   `json:"existing_id,omitempty"`
}

func (s *LoadJudgeService) LoadJudge(request LoadJudgeRequest) (*LoadJudgeResponse, error) {
	log.Printf("開始載入評判題目: %s", request.ZipFilePath)

	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("judge_temp_%d", time.Now().UnixNano()))
	defer os.RemoveAll(tempDir)

	if err := s.extractZipFile(request.ZipFilePath, tempDir); err != nil {
		return nil, fmt.Errorf("解壓縮檔案失敗: %w", err)
	}

	config, _, err := s.loadProblemConfig(tempDir)
	if err != nil {
		return nil, fmt.Errorf("載入問題配置失敗: %w", err)
	}

	existingJudge, isDuplicate := s.checkDuplicateJudge(config, request.CreatorID)

	if isDuplicate && !request.ForceUpdate {
		return &LoadJudgeResponse{
			JudgeID:    0,
			Name:       config.Name,
			Success:    false,
			Message:    fmt.Sprintf("題目 '%s' 已存在。如需更新，請設定 force_update=true", config.Name),
			IsUpdate:   false,
			ExistingID: existingJudge.ID,
		}, nil
	}

	if err := os.MkdirAll(s.judgesBaseDir, 0755); err != nil {
		return nil, fmt.Errorf("創建 judges 目錄失敗: %w", err)
	}

	var judgeDirName string
	var judgeDir string
	var isUpdate bool

	if isDuplicate && request.ForceUpdate {
		isUpdate = true
		existingPath, err := s.getJudgeDirectoryFromFullSetting(existingJudge)
		if err != nil {
			judgeDirName = s.generateJudgeDirName(request.JudgeName)
			judgeDir = filepath.Join(s.judgesBaseDir, judgeDirName)
			judgeDir = s.ensureUniqueJudgeDir(judgeDir)
		} else {
			judgeDir = existingPath
			backupDir := existingPath + "_backup_" + time.Now().Format("20060102_150405")
			if err := utils.CopyDirectory(existingPath, backupDir); err != nil {
				log.Printf("警告: 創建備份失敗: %v", err)
			} else {
				log.Printf("已創建備份: %s", backupDir)
			}
		}
	} else {
		judgeDirName = s.generateJudgeDirName(request.JudgeName)
		judgeDir = filepath.Join(s.judgesBaseDir, judgeDirName)
		judgeDir = s.ensureUniqueJudgeDir(judgeDir)
	}

	if err := s.extractZipFile(request.ZipFilePath, judgeDir); err != nil {
		os.RemoveAll(judgeDir)
		return nil, fmt.Errorf("解壓縮檔案失敗: %w", err)
	}

	if err := s.normalizeExtractedDirectory(judgeDir); err != nil {
		log.Printf("警告: 正規化目錄結構失敗: %v", err)

	}

	config, configDir, err := s.loadProblemConfig(judgeDir)
	if err != nil {
		os.RemoveAll(judgeDir)
		return nil, fmt.Errorf("載入問題配置失敗: %w", err)
	}
	for _, lang := range config.Languages {
		langType, err := language.StringToLanguageType(lang)
		if err != nil {
			log.Printf("警告: 不支援的語言 '%s'，將被忽略: %v", lang, err)
			continue
		}
		langStr := string(langType)
		switch langType {
		case language.LanguageTypePython:
			makejudge.GenerateMainSolutionFile(configDir, langStr, config)
		case language.LanguageTypeCPP:
			makejudge.GenerateMainSolutionFile(configDir, langStr, config)
		case language.LanguageTypeJavaScript:
			makejudge.GenerateMainSolutionFile(configDir, langStr, config)
		default:
			log.Printf("警告: 語言 '%s' 尚未實作程式碼範本生成", lang)
			continue
		}
	}

	var judge *models.Judge
	if isUpdate {
		judge, err = s.updateJudgeRecord(existingJudge, config, configDir)
		if err != nil {
			os.RemoveAll(judgeDir)
			return nil, fmt.Errorf("更新資料庫記錄失敗: %w", err)
		}
	} else {
		judge, err = s.createJudgeRecord(config, request.CreatorID, configDir)
		if err != nil {
			os.RemoveAll(judgeDir)
			return nil, fmt.Errorf("創建資料庫記錄失敗: %w", err)
		}
	}

	message := "評判題目載入成功"
	if isUpdate {
		message = "評判題目更新成功"
	}

	response := &LoadJudgeResponse{
		JudgeID:       judge.ID,
		Name:          judge.Name,
		Description:   judge.Description,
		PermanentPath: configDir,
		Success:       true,
		Message:       message,
		IsUpdate:      isUpdate,
	}

	log.Printf("評判題目載入成功: ID=%d, Name=%s, Path=%s", judge.ID, judge.Name, configDir)
	return response, nil
}

func (s *LoadJudgeService) LoadJudgeFromDirectory(request LoadJudgeFromDirectoryRequest) (*LoadJudgeResponse, error) {
	log.Printf("開始從目錄載入評判題目: %s", request.DirectoryName)

	if request.DirectoryName == "" {
		return nil, fmt.Errorf("目錄名稱不能為空")
	}

	judgeDir := filepath.Join(s.judgesBaseDir, request.DirectoryName)

	if _, err := os.Stat(judgeDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("指定的目錄不存在: %s", request.DirectoryName)
	}

	config, configDir, err := s.loadProblemConfig(judgeDir)
	if err != nil {
		return nil, fmt.Errorf("載入問題配置失敗: %w", err)
	}

	existingJudge, isDuplicate := s.checkDuplicateJudge(config, request.CreatorID)

	if isDuplicate && !request.ForceUpdate {
		return &LoadJudgeResponse{
			JudgeID:    0,
			Name:       config.Name,
			Success:    false,
			Message:    fmt.Sprintf("題目 '%s' 已存在。如需更新，請設定 force_update=true", config.Name),
			IsUpdate:   false,
			ExistingID: existingJudge.ID,
		}, nil
	}

	for _, lang := range config.Languages {
		langType, err := language.StringToLanguageType(lang)
		if err != nil {
			log.Printf("警告: 不支援的語言 '%s'，將被忽略: %v", lang, err)
			continue
		}
		langStr := string(langType)
		switch langType {
		case language.LanguageTypePython:
			makejudge.GenerateMainSolutionFile(configDir, langStr, config)
		case language.LanguageTypeCPP:
			makejudge.GenerateMainSolutionFile(configDir, langStr, config)
		case language.LanguageTypeJavaScript:
			makejudge.GenerateMainSolutionFile(configDir, langStr, config)
		default:
			log.Printf("警告: 語言 '%s' 尚未實作程式碼範本生成", lang)
			continue
		}
	}

	var judge *models.Judge
	var isUpdate bool
	if isDuplicate && request.ForceUpdate {
		isUpdate = true
		judge, err = s.updateJudgeRecord(existingJudge, config, configDir)
		if err != nil {
			return nil, fmt.Errorf("更新資料庫記錄失敗: %w", err)
		}
	} else {
		judge, err = s.createJudgeRecord(config, request.CreatorID, configDir)
		if err != nil {
			return nil, fmt.Errorf("創建資料庫記錄失敗: %w", err)
		}
	}

	message := "評判題目從目錄載入成功"
	if isUpdate {
		message = "評判題目從目錄更新成功"
	}

	response := &LoadJudgeResponse{
		JudgeID:       judge.ID,
		Name:          judge.Name,
		Description:   judge.Description,
		PermanentPath: configDir,
		Success:       true,
		Message:       message,
		IsUpdate:      isUpdate,
	}

	log.Printf("評判題目從目錄載入成功: ID=%d, Name=%s, Path=%s", judge.ID, judge.Name, configDir)
	return response, nil
}

func (s *LoadJudgeService) checkDuplicateJudge(config loadjudge.ProblemConfig, creatorID int) (models.Judge, bool) {
	judges, err := s.judgeRepo.GetJudgesByCreator(uint(creatorID))
	if err != nil {
		log.Printf("檢查重複 judge 時發生錯誤: %v", err)
		return models.Judge{}, false
	}

	for _, judge := range judges {
		if judge.Name == config.Name {
			return judge, true
		}
	}

	return models.Judge{}, false
}

func (s *LoadJudgeService) updateJudgeRecord(existingJudge models.Judge, config loadjudge.ProblemConfig, configDir string) (*models.Judge, error) {
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
		"config":       config,
		"problem_path": configDir,
		"checker_file": config.Checker.File,
		"checker_lang": config.Checker.Language,
		"fullsetting":  json.RawMessage(settingJSON),
		"updated_at":   time.Now().Unix(),
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

	existingJudge.Name = config.Name
	existingJudge.Description = config.Description
	existingJudge.Readme = readme
	existingJudge.Tags = config.Tags
	existingJudge.Langs = config.Languages
	existingJudge.Type = config.Type
	existingJudge.MemoryLimit = config.MemoryLimitMb
	existingJudge.TimeLimit = int(config.TimeLimitSec)
	existingJudge.CheckType = config.CheckType
	existingJudge.Difficulty = level

	if err := existingJudge.FullSetting.Set(fullSettingJSON); err != nil {
		return nil, fmt.Errorf("設定 FullSetting 失敗: %w", err)
	}

	if err := s.judgeRepo.Update(existingJudge); err != nil {
		return nil, fmt.Errorf("更新 Judge 記錄失敗: %w", err)
	}

	return &existingJudge, nil
}

func (s *LoadJudgeService) getJudgeDirectoryFromFullSetting(judge models.Judge) (string, error) {
	var fullSetting map[string]interface{}
	if err := judge.FullSetting.AssignTo(&fullSetting); err != nil {
		return "", fmt.Errorf("解析 FullSetting 失敗: %w", err)
	}

	if problemPath, ok := fullSetting["problem_path"].(string); ok {
		return problemPath, nil
	}

	return "", fmt.Errorf("無法從 FullSetting 中找到問題路徑")
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

func (s *LoadJudgeService) normalizeExtractedDirectory(extractDir string) error {
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return fmt.Errorf("讀取解壓目錄失敗: %w", err)
	}

	var dirs []os.DirEntry
	var files []os.DirEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	if len(dirs) != 1 || len(files) > 0 {
		return nil
	}

	subDirName := dirs[0].Name()
	subDirPath := filepath.Join(extractDir, subDirName)

	configPath := filepath.Join(subDirPath, "config.yaml")
	if !loadjudge.FileExists(configPath) {
		return nil
	}

	tempDir := subDirPath + "_temp"
	if err := os.Rename(subDirPath, tempDir); err != nil {
		return fmt.Errorf("重命名子目錄失敗: %w", err)
	}

	tempEntries, err := os.ReadDir(tempDir)
	if err != nil {
		os.Rename(tempDir, subDirPath)
		return fmt.Errorf("讀取臨時目錄失敗: %w", err)
	}

	for _, entry := range tempEntries {
		srcPath := filepath.Join(tempDir, entry.Name())
		destPath := filepath.Join(extractDir, entry.Name())

		if err := os.Rename(srcPath, destPath); err != nil {
			log.Printf("警告: 移動項目 %s 失敗: %v", entry.Name(), err)
			os.Rename(tempDir, subDirPath)
			return fmt.Errorf("移動目錄內容失敗: %w", err)
		}
	}

	if err := os.Remove(tempDir); err != nil {
		log.Printf("警告: 刪除臨時目錄失敗: %v", err)
	}

	return nil
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
		"config":       config,
		"problem_path": configDir,
		"checker_file": config.Checker.File,
		"checker_lang": config.Checker.Language,
		"fullsetting":  json.RawMessage(settingJSON),
		"created_at":   time.Now().Unix(),
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

func (s *LoadJudgeService) ListAvailableJudgeDirectories() ([]string, error) {
	entries, err := os.ReadDir(s.judgesBaseDir)
	if err != nil {
		return nil, fmt.Errorf("讀取judges目錄失敗: %w", err)
	}

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			// 檢查是否包含config.yaml檔案
			configPath := filepath.Join(s.judgesBaseDir, entry.Name(), "config.yaml")
			if loadjudge.FileExists(configPath) {
				directories = append(directories, entry.Name())
			}
		}
	}

	return directories, nil
}
