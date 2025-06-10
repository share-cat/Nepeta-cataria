package boot

import (
	"fmt"

	"share-cat/Nepeta-cataria/models"

	"gorm.io/gorm"
)

func InitSupportLangs(db *gorm.DB) error {
	supportedLangs := []models.SupportLang{
		{
			Name:           "cpp",
			Version:        "no",
			CompileCommand: "no",
			RunCommand:     "no",
			FileExtension:  ".cpp",
			IsActive:       true,
		},
		{
			Name:           "python",
			Version:        "no",
			CompileCommand: "no",
			RunCommand:     "no",
			FileExtension:  ".py",
			IsActive:       true,
		},
		{
			Name:           "nodejs",
			Version:        "no",
			CompileCommand: "no",
			RunCommand:     "no",
			FileExtension:  ".js",
			IsActive:       true,
		},
	}

	for _, lang := range supportedLangs {
		var existingLang models.SupportLang
		result := db.Where("name = ?", lang.Name).First(&existingLang)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				if err := db.Create(&lang).Error; err != nil {
					return fmt.Errorf("failed to create supported language %s: %w", lang.Name, err)
				}
			} else {
				return fmt.Errorf("failed to query supported language %s: %w", lang.Name, result.Error)
			}
		}
	}

	return nil
}
