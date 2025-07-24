package models

import (
	"encoding/json"
	"log"
	"time"

	"share-cat/Nepeta-cataria/pkg/makejudge"

	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Judge struct {
	gorm.Model
	Name         string         `json:"name" gorm:"type:varchar(255);not null"`
	Description  string         `json:"description" gorm:"type:text"`
	Readme       string         `json:"readme" gorm:"type:text"`
	CreatorID    int            `json:"creator_id" gorm:"not null"`
	Tags         pq.StringArray `json:"tags" gorm:"type:text[];column:tag"`
	Langs        pq.StringArray `json:"langs" gorm:"type:text[]"`
	Type         string         `json:"type" gorm:"type:varchar(255);not null"`
	MemoryLimit  int            `json:"memory_limit" gorm:"default:0"`
	TimeLimit    int            `json:"time_limit" gorm:"default:0"`
	CheckType    string         `json:"check_type" gorm:"type:varchar(255);not null"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	Difficulty   int            `json:"difficulty" gorm:"default:0"`
	FullSetting  pgtype.JSONB   `json:"-" gorm:"type:jsonb;default:'{}'"`
	Creator      User           `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	Solutions    []Solution     `json:"solutions,omitempty" gorm:"foreignKey:JudgeID"`
	ChatSessions []ChatSession  `json:"chat_sessions,omitempty" gorm:"foreignKey:JudgeID"`
}

func (e *Judge) GetSetting() makejudge.SettingData {
	var setting makejudge.SettingData

	if e.FullSetting.Status != pgtype.Present {
		log.Printf("FullSetting is not present, status: %v", e.FullSetting.Status)
		return makejudge.SettingData{}
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(e.FullSetting.Bytes, &rawData); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return makejudge.SettingData{}
	}

	var settingData interface{}
	if fullsetting, exists := rawData["fullsetting"]; exists {
		settingData = fullsetting
	} else {
		settingData = rawData
	}

	settingBytes, err := json.Marshal(settingData)
	if err != nil {
		log.Printf("Error marshaling setting data: %v", err)
		return makejudge.SettingData{}
	}

	if err := json.Unmarshal(settingBytes, &setting); err != nil {
		log.Printf("Error unmarshaling to SettingData: %v", err)
		return makejudge.SettingData{}
	}
	return setting
}

type Permission struct {
	gorm.Model

	IsAdmin   bool `json:"is_admin" gorm:"default:false"`
	IsTeacher bool `json:"is_teacher" gorm:"default:false"`
	IsStudent bool `json:"is_student" gorm:"default:true"`

	Users []User `json:"users,omitempty" gorm:"foreignKey:PermissionID"`
}

type School struct {
	gorm.Model

	SchoolName string `json:"school_name" gorm:"type:varchar(255);not null"`
	Department string `json:"department" gorm:"type:varchar(255)"`

	UserProfiles []UserProfile `json:"user_profiles,omitempty" gorm:"foreignKey:SchoolID"`
}

type Social struct {
	gorm.Model

	Github    string `json:"github" gorm:"type:varchar(255)"`
	Twitter   string `json:"twitter" gorm:"type:varchar(255)"`
	Bluesky   string `json:"bluesky" gorm:"type:varchar(255)"`
	Instagram string `json:"instagram" gorm:"type:varchar(255)"`
	Web       string `json:"web" gorm:"type:varchar(255)"`

	UserProfile *UserProfile `json:"user_profile,omitempty" gorm:"foreignKey:SocialID"`
}

type Solution struct {
	gorm.Model

	UserID      int       `json:"user_id" gorm:"not null;index"`
	JudgeID     int       `json:"judge_id" gorm:"not null;index"`
	SubmittedAt time.Time `json:"submitted_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Code        string    `json:"code" gorm:"type:text;not null"`
	LanguageID  int       `json:"language_id" gorm:"not null"`
	Status      string    `json:"status" gorm:"type:varchar(255);not null;default:'pending'"`

	User     User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Judge    Judge          `json:"judge,omitempty" gorm:"foreignKey:JudgeID"`
	Language SupportLang    `json:"language,omitempty" gorm:"foreignKey:LanguageID"`
	Tasks    []SolutionTask `json:"tasks,omitempty" gorm:"foreignKey:SolutionID"`
}

type SolutionTask struct {
	gorm.Model
	SolutionID    uint   `json:"solution_id" gorm:"not null;index"`
	TaskIndex     int    `json:"task_index" gorm:"not null"`
	Status        string `json:"status" gorm:"type:varchar(255);not null;default:'pending'"`
	ExecutionTime int64  `json:"execution_time" gorm:"default:0"`
	MemoryUsage   int64  `json:"memory_usage" gorm:"default:0"`
	Score         int    `json:"score" gorm:"default:0"`
	Stdout        string `json:"stdout" gorm:"type:text"`
	Stderr        string `json:"stderr" gorm:"type:text"`
	ExitCode      int    `json:"exit_code" gorm:"default:0"`

	Solution Solution `json:"solution,omitempty" gorm:"foreignKey:SolutionID"`
}

type SupportLang struct {
	// 支援的語言
	gorm.Model
	Name           string `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Version        string `json:"version" gorm:"type:varchar(255);not null"`
	CompileCommand string `json:"compile_command" gorm:"type:text"`
	RunCommand     string `json:"run_command" gorm:"type:text;not null"`
	FileExtension  string `json:"file_extension" gorm:"type:varchar(255);not null"`
	IsActive       bool   `json:"is_active" gorm:"not null;default:true"`

	Solutions []Solution `json:"solutions,omitempty" gorm:"foreignKey:LanguageID"`
}

type SystemConfig struct {
	gorm.Model
	ConfigKey   string `json:"config_key" gorm:"type:varchar(255);not null;uniqueIndex"`
	ConfigValue string `json:"config_value" gorm:"type:text;not null"`
	Description string `json:"description" gorm:"type:text"`
}

type User struct {
	// 使用者
	gorm.Model
	Username     string    `json:"username" gorm:"type:varchar(255);not null;uniqueIndex"`
	Password     string    `json:"-" gorm:"type:varchar(255);not null"`
	RegisterTime time.Time `json:"register_time" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Email        string    `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	AvatarUrl    string    `json:"avatar_url" gorm:"type:varchar(255)"`
	PermissionID int       `json:"permission_id"`

	Permission    *Permission    `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
	UserProfile   *UserProfile   `json:"user_profile,omitempty" gorm:"foreignKey:UserID"`
	Solutions     []Solution     `json:"solutions,omitempty" gorm:"foreignKey:UserID"`
	Judges        []Judge        `json:"judges,omitempty" gorm:"foreignKey:CreatorID"`
	Discussions   []Discussion   `json:"discussions,omitempty" gorm:"foreignKey:AuthorID"`
	ChatSessions  []ChatSession  `json:"chat_sessions,omitempty" gorm:"foreignKey:UserID"`
	Achievements  []Achievement  `json:"achievements,omitempty" gorm:"many2many:user_achievements"`
	SubmitEffects []SubmitEffect `json:"submit_effects,omitempty" gorm:"many2many:user_submit_effects"`
	RpgProfile    *RpgProfile    `json:"rpg_profile,omitempty" gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	// 使用者個人資料
	gorm.Model
	UserID    int    `json:"user_id" gorm:"not null;uniqueIndex"`
	SchoolID  *int   `json:"school_id"`
	SocialID  int    `json:"social_id" gorm:"not null"`
	AvatarURL string `json:"avatar_url" gorm:"type:varchar(255)"`

	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	School *School `json:"school,omitempty" gorm:"foreignKey:SchoolID"`
	Social *Social `json:"social,omitempty" gorm:"foreignKey:SocialID"`
}

type ChatSession struct {
	// Ai 聊天紀錄
	gorm.Model
	UserID      int          `json:"user_id" gorm:"not null;index"`
	JudgeID     int          `json:"judge_id" gorm:"not null;index"`
	Title       string       `json:"title" gorm:"type:varchar(255);not null"`
	ChatHistory pgtype.JSONB `json:"-" gorm:"type:jsonb"`

	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Judge Judge `json:"judge,omitempty" gorm:"foreignKey:JudgeID"`
}

type Discussion struct {
	// 討論區
	gorm.Model
	Title      string         `json:"title" gorm:"type:varchar(255);not null"`
	Content    string         `json:"content" gorm:"type:text;not null"`
	Type       string         `json:"type" gorm:"type:varchar(255);not null"`
	Tags       pq.StringArray `json:"tags" gorm:"type:text[]"`
	ParentID   int            `json:"parent_id" gorm:"column:pid"`                       // Parent ID
	ParentType string         `json:"parent_type" gorm:"type:varchar(255);column:ptype"` // Parent Type
	Views      int            `json:"views" gorm:"default:0"`
	AuthorID   int            `json:"author_id" gorm:"not null;index"`

	Author User `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
}

type Achievement struct {
	// 成就 - 徽章
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	Rarelevel   int    `json:"rare_level" gorm:"default:0"`
	Icon        string `json:"icon" gorm:"type:varchar(255);not null"`
}

type SubmitEffect struct {
	// 提交特效
	gorm.Model
	Name              string `json:"name" gorm:"type:varchar(255);not null"`
	Description       string `json:"description" gorm:"type:text"`
	Rarelevel         int    `json:"rare_level" gorm:"default:0"`
	Icon              string `json:"icon" gorm:"type:varchar(255);not null"`
	EffectId          string `json:"effect_id" gorm:"type:varchar(255);not null"`
	UserSubmitEffects []User `json:"user_submit_effects,omitempty" gorm:"many2many:user_submit_effects"`
}

type RpgProfile struct {
	// RPG user profile
	gorm.Model
	UserID int `json:"user_id" gorm:"primaryKey"`
	Exp    int `json:"exp" gorm:"default:0"`
	Level  int `json:"level" gorm:"default:0"`
}
