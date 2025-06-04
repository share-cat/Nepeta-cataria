package makejudge

import "share-cat/Nepeta-cataria/pkg/loadjudge"

const (
	DisplayFileName = "display.json"
)

type SettingData struct {
	ProblemTitle       string               `json:"problemTitle"`
	ShortDescription   string               `json:"shortDescription"`
	FullDescription    string               `json:"fullDescription"`
	SampleInput        string               `json:"sampleInput"`
	SampleOutput       string               `json:"sampleOutput"`
	SupportedLanguages []string             `json:"supportedLanguages"`
	DifficultyLevel    string               `json:"difficultyLevel"`
	Tags               []string             `json:"tags"`
	ProblemType        string               `json:"problemType"`
	TimeLimitSec       float64              `json:"timeLimitSec"`
	MemoryLimitMb      int                  `json:"memoryLimitMb"`
	CheckerFile        string               `json:"checkerFile,omitempty"`
	CheckerLanguage    string               `json:"checkerLanguage,omitempty"`
	CheckType          string               `json:"checkType,omitempty"`
	TestCases          []loadjudge.TestCase `json:"testCases,omitempty"`
	InputFormat        struct {
		Type string `json:"type"`
	} `json:"inputFormat,omitempty"`
	OutputFormat struct {
		Type string `json:"type"`
	} `json:"outputFormat,omitempty"`
}

type TemplateData struct {
	ProblemName          string
	SampleInput          string
	SampleOutput         string
	InferredInputTypes   []string
	InferredOutputType   string
	ParsingLogicHint     string
	FunctionSignature    string
	ExecutionBoilerplate string
	UserCodePlaceholder  string
	RequiredImports      string
	RequiredIncludes     string
	InputType            *loadjudge.TypeInfo
	OutputType           *loadjudge.TypeInfo
}
