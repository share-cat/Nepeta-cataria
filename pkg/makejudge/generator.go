package makejudge

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"share-cat/Nepeta-cataria/pkg/language"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
)

func GenerateMainSolutionFile(problemDir, targetLang string, config loadjudge.ProblemConfig) (string, error) {
	if targetLang == "" {
		if len(config.Languages) > 0 {
			langType, err := language.StringToLanguageType(config.Languages[0])
			if err != nil {
				targetLang = string(language.LanguageTypePython)
				log.Printf("Cannot parse language '%s', defaulting to %s", config.Languages[0], targetLang)
			} else {
				targetLang = string(langType)
			}
		} else {
			targetLang = string(language.LanguageTypePython)
			log.Printf("No target language specified and config.lang is empty, defaulting to %s", targetLang)
		}
	}
	log.Printf("Generating solution file for language: %s in dir: %s", targetLang, problemDir)

	sampleInput, sampleOutput, err := loadjudge.GetSampleInputOutput(problemDir)
	if err != nil {
		log.Printf("Warning: Could not get sample input/output: %v", err)
		sampleInput = "(No sample input provided for template hints)"
		sampleOutput = "(No sample output provided for template hints)"
	}

	var td TemplateData
	td.ProblemName = config.Name
	td.SampleInput = sampleInput
	td.SampleOutput = sampleOutput

	// 嘗試從 testdata.yaml 載入類型信息
	testData, testDataErr := loadjudge.LoadTestDataConfig(problemDir)
	var inputType, outputType *loadjudge.TypeInfo

	if testDataErr == nil {
		log.Printf("Using type information from testdata.yaml")
		inputType, outputType, err = loadjudge.ParseTestDataTypes(testData)
		if err != nil {
			log.Printf("Warning: Failed to parse types from testdata.yaml: %v", err)
			log.Printf("Falling back to type inference from sample data")
		} else {
			td.InputType = inputType
			td.OutputType = outputType
		}
	}

	var generatedCode string
	err = fillTemplateDataWithTypes(&td, targetLang, td.InputType, td.OutputType)
	if err != nil {
		return "", fmt.Errorf("failed to fill template data with types: %w", err)
	}

	generatedCode, err = executeCodeTemplateWithTypes(targetLang, td)
	if err != nil {
		return "", fmt.Errorf("failed to generate code from enhanced template for %s: %w", targetLang, err)
	}
	var mainFileName string
	switch targetLang {
	case string(language.LanguageTypePython):
		mainFileName = "main.py"
	case string(language.LanguageTypeCPP):
		mainFileName = "main.cpp"
	case string(language.LanguageTypeJavaScript):
		mainFileName = "main.js"
	default:
		return "", fmt.Errorf("cannot determine main file name for unsupported language: %s", targetLang)
	}

	outputFilePath := filepath.Join(problemDir, mainFileName)
	err = WriteStringToFile(outputFilePath, generatedCode)
	if err != nil {
		return "", fmt.Errorf("failed to write generated solution file to %s: %w", outputFilePath, err)
	}

	log.Printf("Successfully generated solution template: %s", outputFilePath)
	return outputFilePath, nil
}

func generateSettingData(problemDir string, config loadjudge.ProblemConfig) (SettingData, error) {
	var settingData SettingData
	settingData.ProblemTitle = config.Name
	settingData.ShortDescription = config.Description
	settingData.TimeLimitSec = config.TimeLimitSec
	settingData.MemoryLimitMb = config.MemoryLimitMb
	settingData.SupportedLanguages = config.Languages
	settingData.DifficultyLevel = config.Level
	settingData.Tags = config.Tags
	settingData.ProblemType = config.Type
	settingData.CheckerFile = config.Checker.File
	settingData.CheckerLanguage = config.Checker.Language
	settingData.CheckType = config.CheckType

	readmePath := filepath.Join(problemDir, config.ReadmeFile)
	if loadjudge.FileExists(readmePath) {
		readmeBytes, err := os.ReadFile(readmePath)
		if err != nil {
			log.Printf("Warning: Failed to read readme file at %s: %v", readmePath, err)
			settingData.FullDescription = fmt.Sprintf("(Could not load %s)", config.ReadmeFile)
		} else {
			settingData.FullDescription = string(readmeBytes)
		}
	} else {
		log.Printf("Warning: Readme file '%s' not found at %s", config.ReadmeFile, readmePath)
		settingData.FullDescription = "(Readme file not found)"
	}
	testData, err := loadjudge.LoadTestDataConfig(problemDir)
	if err == nil {
		log.Printf("Using test data from testdata.yaml for display")

		// 將 Sample Input 轉換為 JSON 字符串
		if inputBytes, marshalErr := json.Marshal(testData.Sample.Input); marshalErr == nil {
			settingData.SampleInput = strings.TrimSpace(string(inputBytes))
		} else {
			settingData.SampleInput = "(Failed to marshal sample input)"
		}

		// 處理 Sample Output
		if outputStr, ok := testData.Sample.Output.(string); ok {
			settingData.SampleOutput = strings.TrimSpace(outputStr)
		} else if outputBytes, marshalErr := json.Marshal(testData.Sample.Output); marshalErr == nil {
			settingData.SampleOutput = strings.TrimSpace(string(outputBytes))
		} else {
			settingData.SampleOutput = "(Failed to marshal sample output)"
		}

		settingData.TestCases = testData.TestCases

		// 處理 InputFormat - 現在是數組，取第一個元素
		if len(testData.TestConfig.InputFormat) > 0 {
			settingData.InputFormat = struct {
				Type string `json:"type"`
			}{
				Type: testData.TestConfig.InputFormat[0].Type,
			}
		}

		settingData.OutputFormat = struct {
			Type string `json:"type"`
		}{
			Type: testData.TestConfig.OutputFormat.Type,
		}
	}

	log.Printf("Successfully generated Setting Data for problem: %s", config.Name)
	return settingData, nil
}

func GenerateSettingDataJSON(problemDir string, config loadjudge.ProblemConfig) (string, error) {
	settingData, err := generateSettingData(problemDir, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate Setting Data: %w", err)
	}

	jsonData, err := json.MarshalIndent(settingData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Setting Data to JSON: %w", err)
	}

	return string(jsonData), nil
}

func GenerateSettingDataFile(problemDir string, config loadjudge.ProblemConfig) error {
	jsonData, err := GenerateSettingDataJSON(problemDir, config)
	if err != nil {
		return fmt.Errorf("failed to generate setting JSON: %w", err)
	}

	settingFilePath := filepath.Join(problemDir, DisplayFileName)
	err = WriteStringToFile(settingFilePath, jsonData)
	if err != nil {
		return fmt.Errorf("failed to write setting file to %s: %w", settingFilePath, err)
	}

	log.Printf("Successfully generated setting file: %s", settingFilePath)
	return nil
}

func GetSettingDataJSON(problemDir string, config loadjudge.ProblemConfig) (SettingData, error) {
	settingData, err := generateSettingData(problemDir, config)
	if err != nil {
		return SettingData{}, fmt.Errorf("failed to generate Setting Data: %w", err)
	}
	return settingData, nil

}

func ScaffoldProblem(problemDir, targetLang string) error {
	log.Printf("Scaffolding problem in directory: %s. Target language for template: %s", problemDir, targetLang)
	config, err := loadjudge.LoadProblemConfig(problemDir)
	if err != nil {
		return fmt.Errorf("could not load problem config for %s: %w", problemDir, err)
	}

	effectiveLang := targetLang
	if effectiveLang == "" {
		if len(config.Languages) > 0 {
			effectiveLang = config.Languages[0]
			log.Printf("No specific language requested for template, using first from config: %s", effectiveLang)
		} else {
			return fmt.Errorf("no target language specified for template generation and problem config 'lang' is empty for %s", problemDir)
		}
	}

	_, err = GenerateMainSolutionFile(problemDir, effectiveLang, config)
	if err != nil {
		return fmt.Errorf("could not generate main solution file for %s: %w", problemDir, err)
	}

	err = GenerateSettingDataFile(problemDir, config)
	if err != nil {
		return fmt.Errorf("could not generate frontend display file for %s: %w", problemDir, err)
	}

	log.Printf("Successfully scaffolded problem '%s'. Main template generated for %s.", config.Name, effectiveLang)
	return nil
}
