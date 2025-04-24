package database

import (
	"time"

	"github.com/jackc/pgtype"
)

type Judge struct {
	ID          int       `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	CreatorID   int       `gorm:"not null"`
	Tags        []string  `gorm:"type:varchar(255)[];column:tag"`
	Type        string    `gorm:"type:varchar(255);not null"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	Creator      User          `gorm:"foreignKey:CreatorID"`
	Solutions    []Solution    `gorm:"foreignKey:JudgeID"`
	ChatSessions []ChatSession `gorm:"foreignKey:JudgeID"`
}

type Permission struct {
	ID        int  `gorm:"primaryKey"`
	IsAdmin   bool `gorm:"default:false"`
	IsTeacher bool `gorm:"default:false"`
	IsStudent bool `gorm:"default:true"`

	Users []User `gorm:"foreignKey:PermissionID"`
}

type School struct {
	ID         int    `gorm:"primaryKey"`
	SchoolName string `gorm:"type:varchar(255);not null"`
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
	UserID        int       `gorm:"not null;index"`
	JudgeID       int       `gorm:"not null;index"`
	SubmittedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Code          string    `gorm:"type:text;not null"`
	LanguageID    int       `gorm:"not null"`
	Status        string    `gorm:"type:varchar(255);not null;default:'pending'"`
	ExecutionTime int       `gorm:"default:0"`
	MemoryUsage   int       `gorm:"default:0"`
	Score         int       `gorm:"default:0"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	User     User        `gorm:"foreignKey:UserID"`
	Judge    Judge       `gorm:"foreignKey:JudgeID"`
	Language SupportLang `gorm:"foreignKey:LanguageID"`
}

type SupportLang struct {
	ID             int    `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Version        string `gorm:"type:varchar(255);not null"`
	CompileCommand string `gorm:"type:text"`
	RunCommand     string `gorm:"type:text;not null"`
	FileExtension  string `gorm:"type:varchar(255);not null"`
	IsActive       bool   `gorm:"not null;default:true"`

	Solutions []Solution `gorm:"foreignKey:LanguageID"`
}

type SystemConfig struct {
	ID          int       `gorm:"primaryKey"`
	ConfigKey   string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	ConfigValue string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type User struct {
	ID           int       `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password     string    `gorm:"type:varchar(255);not null"`
	RegisterTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex"`
	AvatarUrl    string    `gorm:"type:varchar(255)"`
	PermissionID *int
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	Permission   *Permission   `gorm:"foreignKey:PermissionID"`
	UserProfile  *UserProfile  `gorm:"foreignKey:UserID"`
	Solutions    []Solution    `gorm:"foreignKey:UserID"`
	Judges       []Judge       `gorm:"foreignKey:CreatorID"`
	Discussions  []Discussion  `gorm:"foreignKey:AuthorID"`
	ChatSessions []ChatSession `gorm:"foreignKey:UserID"`
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
	UserID      int          `gorm:"not null;index"`
	JudgeID     int          `gorm:"not null;index"`
	Title       string       `gorm:"type:varchar(255);not null"`
	ChatHistory pgtype.JSONB `gorm:"type:jsonb"`
	CreatedAt   time.Time    `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time    `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	User  User  `gorm:"foreignKey:UserID"`
	Judge Judge `gorm:"foreignKey:JudgeID"`
}

type Discussion struct {
	ID         int       `gorm:"primaryKey"`
	Title      string    `gorm:"type:varchar(255);not null"`
	Content    string    `gorm:"type:text;not null"`
	Type       string    `gorm:"type:varchar(255);not null"`
	Tags       []string  `gorm:"type:varchar(255)[]"`
	ParentID   int       `gorm:"column:pid"`                     // Parent ID
	ParentType string    `gorm:"type:varchar(255);column:ptype"` // Parent Type
	Views      int       `gorm:"default:0"`
	AuthorID   int       `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`

	Author User `gorm:"foreignKey:AuthorID"`
}
