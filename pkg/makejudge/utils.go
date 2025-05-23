package makejudge

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetLanguageFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".py":
		return "python"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".java":
		return "java"
	case ".js":
		return "javascript"
	default:
		return ""
	}
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteStringToFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}

func InferType(token string) string {
	if token == "" {
		return "string"
	}
	if _, err := strconv.ParseInt(token, 10, 64); err == nil {
		return "int"
	}
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return "float"
	}
	return "string"
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ParseMemoryLimit(memStr string) (int, error) {
	if memStr == "" {
		return 0, fmt.Errorf("memory limit string is empty")
	}

	memStr = strings.ToLower(strings.TrimSpace(memStr))

	valueOnlyRegex := regexp.MustCompile(`^(\d+)$`)
	valueUnitRegex := regexp.MustCompile(`^(\d+)([kmgb])$`)

	if matches := valueOnlyRegex.FindStringSubmatch(memStr); len(matches) == 2 {
		val, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("invalid memory value '%s': %w", matches[1], err)
		}
		return val, nil
	}

	matches := valueUnitRegex.FindStringSubmatch(memStr)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid memory format: '%s'. Expected format like '1024', '512m', '1g', '2048k', '1024b'", memStr)
	}

	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory value '%s': %w", matches[1], err)
	}

	unit := matches[2]
	var valueBytes int64

	switch unit {
	case "b":
		valueBytes = value
	case "k":
		valueBytes = value * 1024
	case "m":
		valueBytes = value * 1024 * 1024
	case "g":
		valueBytes = value * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown memory unit: '%s'", unit)
	}

	mbValue := int(valueBytes / (1024 * 1024))
	if valueBytes > 0 && mbValue == 0 {
		return 1, nil
	}
	return mbValue, nil
}

func ParseTimeLimit(timeStr string) (float64, error) {
	if timeStr == "" {
		return 0, fmt.Errorf("time limit string is empty")
	}
	timeStr = strings.ToLower(strings.TrimSpace(timeStr))

	valueOnlyRegex := regexp.MustCompile(`^(\d+(\.\d+)?)$`)
	valueUnitRegex := regexp.MustCompile(`^(\d+(\.\d+)?)(s|ms)$`)

	if matches := valueUnitRegex.FindStringSubmatch(timeStr); len(matches) == 4 {
		value, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid time value '%s': %w", matches[1], err)
		}
		unit := matches[3]
		if unit == "s" {
			return value, nil
		} else if unit == "ms" {
			return value / 1000.0, nil
		}
	} else if matches := valueOnlyRegex.FindStringSubmatch(timeStr); len(matches) >= 2 {
		val, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid time value '%s': %w", matches[1], err)
		}
		return val / 1000.0, nil
	}

	return 0, fmt.Errorf("invalid time format: '%s'. Expected format like '1000', '1s', '1500ms', '0.5s'", timeStr)
}
