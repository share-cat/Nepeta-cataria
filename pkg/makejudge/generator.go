package makejudge

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GenerateMainSolutionFile(problemDir, targetLang string, config ProblemConfig) (string, error) {
	if targetLang == "" {

		if len(config.Languages) > 0 {
			targetLang = config.Languages[0]
		} else {
			targetLang = Python
			log.Printf("No target language specified and config.lang is empty, defaulting to %s", targetLang)
		}
	}
	log.Printf("Generating solution file for language: %s in dir: %s", targetLang, problemDir)

	sampleInPath := filepath.Join(problemDir, TestDataDir, SampleInputFile)
	sampleOutPath := filepath.Join(problemDir, TestDataDir, SampleOutputFile)

	if !FileExists(sampleInPath) {
		files, err := filepath.Glob(filepath.Join(problemDir, TestDataDir, "*.in"))
		if err != nil || len(files) == 0 {

			log.Printf("No sample input ('%s') or other .in files found in %s/%s. Template will have minimal hints.", SampleInputFile, problemDir, TestDataDir)
			sampleInPath = ""
		} else {
			sampleInPath = files[0]
			log.Printf("Using %s as sample input.", sampleInPath)
		}
	}
	if !FileExists(sampleOutPath) {
		files, err := filepath.Glob(filepath.Join(problemDir, TestDataDir, "*.out"))
		if err != nil || len(files) == 0 {

			log.Printf("No sample output ('%s') or other .out files found in %s/%s. Template will have minimal hints.", SampleOutputFile, problemDir, TestDataDir)
			sampleOutPath = ""
		} else {
			sampleOutPath = files[0]
			log.Printf("Using %s as sample output.", sampleOutPath)
		}
	}

	var td TemplateData
	td.ProblemName = config.Name

	if sampleInPath != "" && FileExists(sampleInPath) {
		sampleInLines, err := ReadLines(sampleInPath)
		if err != nil {
			log.Printf("Warning: Could not read sample input %s: %v.", sampleInPath, err)
			td.SampleInput = fmt.Sprintf("(Could not read %s)", sampleInPath)
		} else {
			td.SampleInput = strings.Join(sampleInLines, "\n")
			if len(sampleInLines) > 0 {
				tokens := strings.Fields(sampleInLines[0])
				for _, token := range tokens {
					td.InferredInputTypes = append(td.InferredInputTypes, InferType(token))
				}
			}
		}
	} else {
		td.SampleInput = "(No sample input provided for template hints)"
	}

	if sampleOutPath != "" && FileExists(sampleOutPath) {
		sampleOutLines, err := ReadLines(sampleOutPath)
		if err != nil {
			log.Printf("Warning: Could not read sample output %s: %v.", sampleOutPath, err)
			td.SampleOutput = fmt.Sprintf("(Could not read %s)", sampleOutPath)
			td.InferredOutputType = "string"
		} else {
			td.SampleOutput = strings.Join(sampleOutLines, "\n")
			if len(sampleOutLines) > 0 {
				outputTokens := strings.Fields(sampleOutLines[0])
				if len(outputTokens) > 0 {
					td.InferredOutputType = InferType(outputTokens[0])
				} else {
					td.InferredOutputType = "string" // Empty line
				}
			} else {
				td.InferredOutputType = "string" // Empty file
			}
		}
	} else {
		td.SampleOutput = "(No sample output provided for template hints)"
		td.InferredOutputType = "string"
	}

	fillTemplateDataForJudge(&td, targetLang, td.InferredInputTypes)

	generatedCode, err := executeCodeTemplate(targetLang, td)
	if err != nil {
		return "", fmt.Errorf("failed to generate code from template for %s: %w", targetLang, err)
	}

	var mainFileName string
	switch targetLang {
	case Python:
		mainFileName = "main.py"
	case Cpp:
		mainFileName = "main.cpp"
	case JavaScript:
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

func generateFrontendDisplayData(problemDir string, config ProblemConfig) (FrontendDisplayData, error) {
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
	if FileExists(readmePath) {
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

	sampleInPath := filepath.Join(problemDir, TestDataDir, SampleInputFile)
	if !FileExists(sampleInPath) {
		files, errGlob := filepath.Glob(filepath.Join(problemDir, TestDataDir, "*.in"))
		if errGlob == nil && len(files) > 0 {
			sampleInPath = files[0]
		} else {
			sampleInPath = ""
		}
	}

	if sampleInPath != "" && FileExists(sampleInPath) {
		sampleInBytes, err := os.ReadFile(sampleInPath)
		if err != nil {
			log.Printf("Warning: Failed to read sample input from %s for display data: %v", sampleInPath, err)
			displayData.SampleInput = "(Could not load sample input)"
		} else {
			displayData.SampleInput = string(sampleInBytes)
		}
	} else {
		log.Printf("Warning: No sample input file found for display data in %s", problemDir)
		displayData.SampleInput = "(No sample input available)"
	}

	sampleOutPath := ""
	files, errGlob := filepath.Glob(filepath.Join(problemDir, TestDataDir, "*.out"))
	if errGlob == nil && len(files) > 0 {
		sampleOutPath = files[0]
	} else {
		sampleOutPath = ""
	}

	if sampleOutPath != "" && FileExists(sampleOutPath) {
		sampleOutBytes, err := os.ReadFile(sampleOutPath)
		if err != nil {
			log.Printf("Warning: Failed to read sample output from %s for display data: %v", sampleOutPath, err)
			displayData.SampleOutput = "(Could not load sample output)"
		} else {
			displayData.SampleOutput = string(sampleOutBytes)
		}
	} else {
		log.Printf("Warning: No corresponding sample output file found for display data (tried %s)", sampleOutPath)
		displayData.SampleOutput = "(No sample output available)"
	}

	log.Printf("Successfully generated display data for problem: %s", config.Name)
	return displayData, nil
}

func GenerateFrontendDisplayJSON(problemDir string, config ProblemConfig) (string, error) {
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

func GenerateFrontendDisplayFile(problemDir string, config ProblemConfig) error {
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

func GetFrontendDisplayJSON(problemDir string, config ProblemConfig) (FrontendDisplayData, error) {
	displayData, err := generateFrontendDisplayData(problemDir, config)
	if err != nil {
		return FrontendDisplayData{}, fmt.Errorf("failed to generate display data: %w", err)
	}
	return displayData, nil

}

func ScaffoldProblem(problemDir, targetLang string) error {
	log.Printf("Scaffolding problem in directory: %s. Target language for template: %s", problemDir, targetLang)
	config, err := LoadProblemConfig(problemDir)
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
