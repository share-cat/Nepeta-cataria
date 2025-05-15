package models

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Judge struct {
	gorm.Model
	Name         string        `gorm:"type:varchar(255);not null"`
	Description  string        `gorm:"type:text"`
	CreatorID    int           `gorm:"not null"`
	Tags         []string      `gorm:"type:varchar(255)[];column:tag"`
	Type         string        `gorm:"type:varchar(255);not null"`
	IsActive     bool          `gorm:"default:true"`
	Difficulty   int           `gorm:"default:0"`
	Creator      User          `gorm:"foreignKey:CreatorID"`
	Solutions    []Solution    `gorm:"foreignKey:JudgeID"`
	ChatSessions []ChatSession `gorm:"foreignKey:JudgeID"`
}

type Permission struct {
	gorm.Model

	IsAdmin   bool `gorm:"default:false"`
	IsTeacher bool `gorm:"default:false"`
	IsStudent bool `gorm:"default:true"`

	Users []User `gorm:"foreignKey:PermissionID"`
}

type School struct {
	gorm.Model

	SchoolName string `gorm:"type:varchar(255);not null"`
	Department string `gorm:"type:varchar(255)"`

	UserProfiles []UserProfile `gorm:"foreignKey:SchoolID"`
}

type Social struct {
	gorm.Model

	Github    string `gorm:"type:varchar(255)"`
	Twitter   string `gorm:"type:varchar(255)"`
	Bluesky   string `gorm:"type:varchar(255)"`
	Instagram string `gorm:"type:varchar(255)"`
	Web       string `gorm:"type:varchar(255)"`

	UserProfile *UserProfile `gorm:"foreignKey:SocialID"`
}

type Solution struct {
	gorm.Model

	UserID        int       `gorm:"not null;index"`
	JudgeID       int       `gorm:"not null;index"`
	SubmittedAt   time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Code          string    `gorm:"type:text;not null"`
	LanguageID    int       `gorm:"not null"`
	Status        string    `gorm:"type:varchar(255);not null;default:'pending'"`
	ExecutionTime int       `gorm:"default:0"`
	MemoryUsage   int       `gorm:"default:0"`
	Score         int       `gorm:"default:0"`

	User     User        `gorm:"foreignKey:UserID"`
	Judge    Judge       `gorm:"foreignKey:JudgeID"`
	Language SupportLang `gorm:"foreignKey:LanguageID"`
}

type SupportLang struct {
	// 支援的語言
	gorm.Model
	Name           string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Version        string `gorm:"type:varchar(255);not null"`
	CompileCommand string `gorm:"type:text"`
	RunCommand     string `gorm:"type:text;not null"`
	FileExtension  string `gorm:"type:varchar(255);not null"`
	IsActive       bool   `gorm:"not null;default:true"`

	Solutions []Solution `gorm:"foreignKey:LanguageID"`
}

type SystemConfig struct {
	gorm.Model
	ConfigKey   string `gorm:"type:varchar(255);not null;uniqueIndex"`
	ConfigValue string `gorm:"type:text;not null"`
	Description string `gorm:"type:text"`
}

type User struct {
	// 使用者
	gorm.Model
	Username     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password     string    `gorm:"type:varchar(255);not null"`
	RegisterTime time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex"`
	AvatarUrl    string    `gorm:"type:varchar(255)"`
	PermissionID int

	Permission    *Permission    `gorm:"foreignKey:PermissionID"`
	UserProfile   *UserProfile   `gorm:"foreignKey:UserID"`
	Solutions     []Solution     `gorm:"foreignKey:UserID"`
	Judges        []Judge        `gorm:"foreignKey:CreatorID"`
	Discussions   []Discussion   `gorm:"foreignKey:AuthorID"`
	ChatSessions  []ChatSession  `gorm:"foreignKey:UserID"`
	Achievements  []Achievement  `gorm:"many2many:user_achievements"`
	SubmitEffects []SubmitEffect `gorm:"many2many:user_submit_effects"`
	RpgProfile    *RpgProfile    `gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	// 使用者個人資料
	gorm.Model
	UserID    int `gorm:"primaryKey"`
	SchoolID  int
	SocialID  int    `gorm:"not null"`
	AvatarURL string `gorm:"type:varchar(255)"`

	User   *User   `gorm:"foreignKey:UserID"`
	School *School `gorm:"foreignKey:SchoolID"`
	Social *Social `gorm:"foreignKey:SocialID"`
}

type ChatSession struct {
	// Ai 聊天紀錄
	gorm.Model
	UserID      int          `gorm:"not null;index"`
	JudgeID     int          `gorm:"not null;index"`
	Title       string       `gorm:"type:varchar(255);not null"`
	ChatHistory pgtype.JSONB `gorm:"type:jsonb"`

	User  User  `gorm:"foreignKey:UserID"`
	Judge Judge `gorm:"foreignKey:JudgeID"`
}

type Discussion struct {
	// 討論區
	gorm.Model
	Title      string   `gorm:"type:varchar(255);not null"`
	Content    string   `gorm:"type:text;not null"`
	Type       string   `gorm:"type:varchar(255);not null"`
	Tags       []string `gorm:"type:varchar(255)[]"`
	ParentID   int      `gorm:"column:pid"`                     // Parent ID
	ParentType string   `gorm:"type:varchar(255);column:ptype"` // Parent Type
	Views      int      `gorm:"default:0"`
	AuthorID   int      `gorm:"not null;index"`

	Author User `gorm:"foreignKey:AuthorID"`
}

type Achievement struct {
	// 成就 - 徽章
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Rarelevel   int    `gorm:"default:0"`
	Icon        string `gorm:"type:varchar(255);not null"`
}

type SubmitEffect struct {
	// 提交特效
	gorm.Model
	Name              string `gorm:"type:varchar(255);not null"`
	Description       string `gorm:"type:text"`
	Rarelevel         int    `gorm:"default:0"`
	Icon              string `gorm:"type:varchar(255);not null"`
	EffectId          string `gorm:"type:varchar(255);not null"`
	UserSubmitEffects []User `gorm:"many2many:user_submit_effects"`
}

type RpgProfile struct {
	// RPG user profile
	gorm.Model
	UserID int `gorm:"primaryKey"`
	Exp    int `gorm:"default:0"`
	Level  int `gorm:"default:0"`
}
