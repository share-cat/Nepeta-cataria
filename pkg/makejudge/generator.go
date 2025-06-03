package makejudge

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"share-cat/Nepeta-cataria/pkg/loadjudge"
)

func GenerateMainSolutionFile(problemDir, targetLang string, config loadjudge.ProblemConfig) (string, error) {
	if targetLang == "" {
		if len(config.Languages) > 0 {
			targetLang = config.Languages[0]
		} else {
			targetLang = loadjudge.Python
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
	case loadjudge.Python:
		mainFileName = "main.py"
	case loadjudge.Cpp:
		mainFileName = "main.cpp"
	case loadjudge.JavaScript:
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

func generateFrontendDisplayData(problemDir string, config loadjudge.ProblemConfig) (FrontendDisplayData, error) {
	var displayData FrontendDisplayData
	displayData.ProblemTitle = config.Name
	displayData.ShortDescription = config.Description
	displayData.TimeLimitSec = config.TimeLimitSec
	displayData.MemoryLimitMb = config.MemoryLimitMb
	displayData.SupportedLanguages = config.Languages
	displayData.DifficultyLevel = config.Level
	displayData.Tags = config.Tags
	displayData.ProblemType = config.Type
	displayData.CheckerFile = config.Checker.File
	displayData.CheckerLanguage = config.Checker.Language
	displayData.CheckType = config.CheckType

	readmePath := filepath.Join(problemDir, config.ReadmeFile)
	if loadjudge.FileExists(readmePath) {
		readmeBytes, err := os.ReadFile(readmePath)
		if err != nil {
			log.Printf("Warning: Failed to read readme file at %s: %v", readmePath, err)
			displayData.FullDescription = fmt.Sprintf("(Could not load %s)", config.ReadmeFile)
		} else {
			displayData.FullDescription = string(readmeBytes)
		}
	} else {
		log.Printf("Warning: Readme file '%s' not found at %s", config.ReadmeFile, readmePath)
		displayData.FullDescription = "(Readme file not found)"
	}
	testData, err := loadjudge.LoadTestDataConfig(problemDir)
	if err == nil {
		log.Printf("Using test data from testdata.yaml for display")
		displayData.SampleInput = strings.TrimSpace(testData.Sample.Input)
		displayData.SampleOutput = strings.TrimSpace(testData.Sample.Output)
		displayData.TestCases = testData.TestCases

		displayData.InputFormat = struct {
			Type string `json:"type"`
		}{
			Type: testData.TestConfig.InputFormat.Type,
		}

		displayData.OutputFormat = struct {
			Type string `json:"type"`
		}{
			Type: testData.TestConfig.OutputFormat.Type,
		}
	}

	log.Printf("Successfully generated display data for problem: %s", config.Name)
	return displayData, nil
}

func GenerateFrontendDisplayJSON(problemDir string, config loadjudge.ProblemConfig) (string, error) {
	displayData, err := generateFrontendDisplayData(problemDir, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate display data: %w", err)
	}

	jsonData, err := json.MarshalIndent(displayData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal display data to JSON: %w", err)
	}

	return string(jsonData), nil
}

func GenerateFrontendDisplayFile(problemDir string, config loadjudge.ProblemConfig) error {
	jsonData, err := GenerateFrontendDisplayJSON(problemDir, config)
	if err != nil {
		return fmt.Errorf("failed to generate display JSON: %w", err)
	}

	displayFilePath := filepath.Join(problemDir, DisplayFileName)
	err = WriteStringToFile(displayFilePath, jsonData)
	if err != nil {
		return fmt.Errorf("failed to write display file to %s: %w", displayFilePath, err)
	}

	log.Printf("Successfully generated display file: %s", displayFilePath)
	return nil
}

func GetFrontendDisplayJSON(problemDir string, config loadjudge.ProblemConfig) (FrontendDisplayData, error) {
	displayData, err := generateFrontendDisplayData(problemDir, config)
	if err != nil {
		return FrontendDisplayData{}, fmt.Errorf("failed to generate display data: %w", err)
	}
	return displayData, nil

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

	err = GenerateFrontendDisplayFile(problemDir, config)
	if err != nil {
		return fmt.Errorf("could not generate frontend display file for %s: %w", problemDir, err)
	}

	log.Printf("Successfully scaffolded problem '%s'. Main template generated for %s.", config.Name, effectiveLang)
	return nil
}
