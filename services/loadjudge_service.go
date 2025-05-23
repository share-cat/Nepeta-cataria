package services

import (
	"share-cat/Nepeta-cataria/repositories"
)

// JudgeConfig 題目設定檔結構
type JudgeConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	README      string   `yaml:"readme"`
	Lang        []string `yaml:"lang"`
	Level       string   `yaml:"level"`
	Tags        []string `yaml:"tags"`
	Type        string   `yaml:"type"`
	Memory      string   `yaml:"mem"`
	Time        string   `yaml:"time"`
	CheckType   string   `yaml:"check_type"`
	Checker     string   `yaml:"checker"`
	Generator   string   `yaml:"generator"`
}

// TestCase 測試案例結構
type TestCase struct {
	ID         int    `json:"id"`
	InputFile  string `json:"input_file"`
	OutputFile string `json:"output_file"`
	Input      string `json:"input"`
	Output     string `json:"output"`
}

// JudgeData 完整的題目資料結構
type JudgeData struct {
	Config      *JudgeConfig      `json:"config"`
	README      string            `json:"readme"`
	Checker     string            `json:"checker"`
	Generator   string            `json:"generator"`
	TestCases   []TestCase        `json:"test_cases"`
	PublicFiles map[string]string `json:"public_files"` // 檔名 -> 內容
	StarterCode map[string]string `json:"starter_code"` // 語言 -> 起始程式碼
}

// LoadJudgeService 載入題目服務
type LoadJudgeService struct {
	judgeService    *JudgeService
	supportLangRepo *repositories.SupportLangRepository
}

// NewLoadJudgeService 創建載入題目服務
func NewLoadJudgeService(
	judgeService *JudgeService,
	supportLangRepo *repositories.SupportLangRepository,
) *LoadJudgeService {
	return &LoadJudgeService{
		judgeService:    judgeService,
		supportLangRepo: supportLangRepo,
	}
}
