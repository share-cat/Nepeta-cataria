package loadjudge

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Python     = "python"
	Cpp        = "cpp"
	JavaScript = "javascript"
)

func LoadProblemConfig(problemDir string) (ProblemConfig, error) {
	var config ProblemConfig
	configPath := filepath.Join(problemDir, ConfigFileName)

	if !FileExists(configPath) {
		return config, fmt.Errorf("config file not found at %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config YAML from %s: %w", configPath, err)
	}

	if config.TimeLimitStr != "" {
		parsedTime, errTime := ParseTimeLimit(config.TimeLimitStr)
		if errTime != nil {
			log.Printf("Warning: Could not parse time limit '%s': %v. Using default.", config.TimeLimitStr, errTime)
			config.TimeLimitSec = 1.0
		} else {
			config.TimeLimitSec = parsedTime
		}
	} else {
		config.TimeLimitSec = 1.0
	}

	if config.MemoryLimitStr != "" {
		parsedMem, errMem := ParseMemoryLimit(config.MemoryLimitStr)
		if errMem != nil {
			log.Printf("Warning: Could not parse memory limit '%s': %v. Using default.", config.MemoryLimitStr, errMem)
			config.MemoryLimitMb = 256
		} else {
			config.MemoryLimitMb = parsedMem
		}
	} else {
		config.MemoryLimitMb = 256
	}

	if config.ReadmeFile == "" {
		config.ReadmeFile = "README.md"
	}

	if config.Checker.File != "" && config.Checker.Language == "" {
		ext := strings.ToLower(filepath.Ext(config.Checker.File))
		if ext == ".py" {
			config.Checker.Language = Python
		} else if ext == ".cpp" || ext == ".cc" {
			config.Checker.Language = Cpp
		} else {
			log.Printf("Warning: Checker language for '%s' not specified and cannot be reliably inferred from extension. Assuming it's executable or Python if script-like.", config.Checker.File)
		}
	}
	if config.Checker.Language != "" && config.Checker.Language != Python && config.Checker.Language != Cpp {
		log.Printf("Warning: Checker language '%s' for '%s' is not explicitly supported by this tool's assumptions (Python/Cpp for CatTeaser context). Ensure judge server handles it.", config.Checker.Language, config.Checker.File)
	}

	return config, nil
}

func LoadTestDataConfig(problemDir string) (TestDataConfig, error) {
	var testData TestDataConfig
	testDataPath := filepath.Join(problemDir, TestDataYamlFile)

	if !FileExists(testDataPath) {
		return testData, fmt.Errorf("testdata.yaml not found at %s", testDataPath)
	}

	data, err := os.ReadFile(testDataPath)
	if err != nil {
		return testData, fmt.Errorf("failed to read testdata file %s: %w", testDataPath, err)
	}

	err = yaml.Unmarshal(data, &testData)
	if err != nil {
		return testData, fmt.Errorf("failed to unmarshal testdata YAML from %s: %w", testDataPath, err)
	}

	log.Printf("Successfully loaded test data config with %d test cases", len(testData.TestCases))
	return testData, nil
}

func GetSampleInputOutput(problemDir string) (string, string, error) {
	testData, err := LoadTestDataConfig(problemDir)
	if err != nil {
		log.Printf("Failed to load testdata.yaml: %v", err)
	}
	log.Printf("Using sample data from testdata.yaml")
	return strings.TrimSpace(testData.Sample.Input), strings.TrimSpace(testData.Sample.Output), nil

}
