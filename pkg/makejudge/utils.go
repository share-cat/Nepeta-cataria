package makejudge

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteStringToFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}
