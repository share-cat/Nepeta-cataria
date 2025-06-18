package language

import (
	"errors"
	"strings"
)

type LanguageType string

const (
	LanguageTypeUnknown    LanguageType = "unknown"
	LanguageTypePython     LanguageType = "python"
	LanguageTypeJava       LanguageType = "java"
	LanguageTypeCPP        LanguageType = "cpp"
	LanguageTypeC          LanguageType = "c"
	LanguageTypeJavaScript LanguageType = "javascript"
	LanguageTypeGo         LanguageType = "go"
	LanguageTypeRust       LanguageType = "rust"
	LanguageTypeKotlin     LanguageType = "kotlin"
	LanguageTypeSwift      LanguageType = "swift"
	LanguageTypePHP        LanguageType = "php"
	LanguageTypeRuby       LanguageType = "ruby"
	LanguageTypeCSharp     LanguageType = "csharp"
	LanguageTypeTypeScript LanguageType = "typescript"
	LanguageTypeScala      LanguageType = "scala"
	LanguageTypePerl       LanguageType = "perl"
	LanguageTypeR          LanguageType = "r"
	LanguageTypeDart       LanguageType = "dart"
	LanguageTypeElixir     LanguageType = "elixir"
	LanguageTypeHaskell    LanguageType = "haskell"
)

type LanguageInfo struct {
	Type      LanguageType `json:"type"`
	Name      string       `json:"name"`
	Extension string       `json:"extension"`
}

var languageMap = map[string]LanguageInfo{

	"python":  {LanguageTypePython, "Python", ".py"},
	"py":      {LanguageTypePython, "Python", ".py"},
	"python3": {LanguageTypePython, "Python", ".py"},
	"python2": {LanguageTypePython, "Python", ".py"},

	"java": {LanguageTypeJava, "Java", ".java"},

	"cpp": {LanguageTypeCPP, "C++", ".cpp"},
	"c++": {LanguageTypeCPP, "C++", ".cpp"},
	"cxx": {LanguageTypeCPP, "C++", ".cpp"},
	"cc":  {LanguageTypeCPP, "C++", ".cpp"},

	"c": {LanguageTypeC, "C", ".c"},

	"javascript": {LanguageTypeJavaScript, "JavaScript", ".js"},
	"js":         {LanguageTypeJavaScript, "JavaScript", ".js"},
	"node":       {LanguageTypeJavaScript, "JavaScript", ".js"},
	"nodejs":     {LanguageTypeJavaScript, "JavaScript", ".js"},

	"go":     {LanguageTypeGo, "Go", ".go"},
	"golang": {LanguageTypeGo, "Go", ".go"},

	"rust": {LanguageTypeRust, "Rust", ".rs"},
	"rs":   {LanguageTypeRust, "Rust", ".rs"},

	"kotlin": {LanguageTypeKotlin, "Kotlin", ".kt"},
	"kt":     {LanguageTypeKotlin, "Kotlin", ".kt"},

	"swift": {LanguageTypeSwift, "Swift", ".swift"},

	"php": {LanguageTypePHP, "PHP", ".php"},

	"ruby": {LanguageTypeRuby, "Ruby", ".rb"},
	"rb":   {LanguageTypeRuby, "Ruby", ".rb"},

	"csharp": {LanguageTypeCSharp, "C#", ".cs"},
	"cs":     {LanguageTypeCSharp, "C#", ".cs"},
	"c#":     {LanguageTypeCSharp, "C#", ".cs"},

	"typescript": {LanguageTypeTypeScript, "TypeScript", ".ts"},
	"ts":         {LanguageTypeTypeScript, "TypeScript", ".ts"},

	"scala": {LanguageTypeScala, "Scala", ".scala"},

	"perl": {LanguageTypePerl, "Perl", ".pl"},
	"pl":   {LanguageTypePerl, "Perl", ".pl"},

	"r": {LanguageTypeR, "R", ".r"},

	"elixir": {LanguageTypeElixir, "Elixir", ".ex"},
	"ex":     {LanguageTypeElixir, "Elixir", ".ex"},

	"haskell": {LanguageTypeHaskell, "Haskell", ".hs"},
	"hs":      {LanguageTypeHaskell, "Haskell", ".hs"},
}

func StringToLanguageType(langStr string) (LanguageType, error) {
	if langStr == "" {
		return LanguageTypeUnknown, errors.New("語言字串不能為空")
	}

	normalized := strings.ToLower(strings.TrimSpace(langStr))

	if info, exists := languageMap[normalized]; exists {
		return info.Type, nil
	}

	return LanguageTypeUnknown, errors.New("不支援的程式語言: " + langStr)
}

func StringToLanguageInfo(langStr string) (LanguageInfo, error) {
	if langStr == "" {
		return LanguageInfo{}, errors.New("語言字串不能為空")
	}

	normalized := strings.ToLower(strings.TrimSpace(langStr))

	if info, exists := languageMap[normalized]; exists {
		return info, nil
	}

	return LanguageInfo{}, errors.New("不支援的程式語言: " + langStr)
}
