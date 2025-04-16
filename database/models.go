package database

import (
	"time"

	"github.com/jackc/pgtype"
)

type Judge struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`
	CreatorID   int    `gorm:"not null"`
	IsActive    bool

	Creator   User       `gorm:"foreignKey:CreatorID"`
	Solutions []Solution `gorm:"foreignKey:JudgeID"`
}

type Permission struct {
	ID        int `gorm:"primaryKey"`
	IsAdmin   bool
	IsTeacher bool
	IsStudent bool

	Users []User `gorm:"foreignKey:PermissionID"`
}

type School struct {
	ID         int    `gorm:"primaryKey"`
	SchoolName string `gorm:"type:varchar(255)"`
	Department string `gorm:"type:varchar(255)"`

	UserProfiles []UserProfile `gorm:"foreignKey:SchoolID"`
}

type Social struct {
	ID        int    `gorm:"primaryKey"`
	Github    string `gorm:"type:varchar(255)"`
	Twitter   string `gorm:"type:varchar(255)"`
	Bluesky   string `gorm:"type:varchar(255)"`
	Instagram string `gorm:"type:varchar(255)"`
	Web       string `gorm:"type:varchar(255)"`

	UserProfile *UserProfile `gorm:"foreignKey:SocialID"`
}

type Solution struct {
	ID            int       `gorm:"primaryKey"`
	UserID        int       `gorm:"not null"`
	JudgeID       int       `gorm:"not null"`
	SubmittedAt   time.Time `gorm:"type:date;not null"`
	Code          string    `gorm:"type:text;not null"`
	LanguageID    int       `gorm:"not null"`
	Status        string    `gorm:"type:varchar(255);not null"`
	ExecutionTime int
	MemoryUsage   int

	User     User        `gorm:"foreignKey:UserID"`
	Judge    Judge       `gorm:"foreignKey:JudgeID"`
	Language SupportLang `gorm:"foreignKey:LanguageID"`
}

type SupportLang struct {
	ID             int    `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(255)"`
	Version        string `gorm:"type:varchar(255)"`
	CompileCommand string `gorm:"type:text"`
	RunCommand     string `gorm:"type:text"`
	FileExtension  string `gorm:"type:varchar(255)"`
	IsActive       bool   `gorm:"not null;default:true"`

	Solutions []Solution `gorm:"foreignKey:LanguageID"`
}

type SystemConfig struct {
	ID          int       `gorm:"primaryKey"`
	ConfigKey   string    `gorm:"type:varchar(255);not null"`
	ConfigValue string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

type User struct {
	ID           int       `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(255);not null"`
	Password     string    `gorm:"type:varchar(255);not null"`
	RegisterTime time.Time `gorm:"type:date;not null"`
	Email        string    `gorm:"type:varchar(255)"`
	PermissionID *int

	Permission  *Permission  `gorm:"foreignKey:PermissionID"`
	UserProfile *UserProfile `gorm:"foreignKey:UserID"`
	Solutions   []Solution   `gorm:"foreignKey:UserID"`
	Judges      []Judge      `gorm:"foreignKey:CreatorID"`
}

type UserProfile struct {
	UserID    int    `gorm:"primaryKey"`
	SchoolID  int    `gorm:"not null"`
	SocialID  int    `gorm:"not null"`
	AvatarURL string `gorm:"type:varchar(255)"`

	User   *User   `gorm:"foreignKey:UserID"`
	School *School `gorm:"foreignKey:SchoolID"`
	Social *Social `gorm:"foreignKey:SocialID"`
}

type ChatSession struct {
	ID          int          `gorm:"primaryKey"`
	UserID      int          `gorm:"not null"`
	JudgeID     int          `gorm:"not null"`
	Title       string       `gorm:"type:varchar(255)"`
	ChatHistory pgtype.JSONB `gorm:"type:jsonb;json:chat_history"`

	User *User `gorm:"foreignKey:UserID"`
}
