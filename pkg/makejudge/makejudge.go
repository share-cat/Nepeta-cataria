package makejudge

const (
	Python     = "python"
	Cpp        = "cpp"
	JavaScript = "javascript"
)

const (
	SampleInputFile  = "sample.in"
	SampleOutputFile = "sample.out"
	TestDataDir      = "data"
	DisplayFileName  = "display.json"
)

type ProblemConfig struct {
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	ReadmeFile     string   `yaml:"readme"`
	Languages      []string `yaml:"lang"`
	Level          string   `yaml:"level"`
	Tags           []string `yaml:"tags"`
	Type           string   `yaml:"type"`
	MemoryLimitStr string   `yaml:"mem"`  // e.g., "1024", "1m", "512b" (parsed to MB)
	TimeLimitStr   string   `yaml:"time"` // e.g., "1000", "1s", "0.5s" (parsed to seconds)
	CheckType      string   `yaml:"check_type"`
	Checker        struct {
		File     string `yaml:"file"`
		Language string `yaml:"language"`
	} `yaml:"checker"`

	TimeLimitSec  float64 `yaml:"-"`
	MemoryLimitMb int     `yaml:"-"`
}

type FrontendDisplayData struct {
	ProblemTitle       string   `json:"problemTitle"`
	ShortDescription   string   `json:"shortDescription"`
	FullDescription    string   `json:"fullDescription"`
	SampleInput        string   `json:"sampleInput"`
	SampleOutput       string   `json:"sampleOutput"`
	SupportedLanguages []string `json:"supportedLanguages"`
	DifficultyLevel    string   `json:"difficultyLevel"`
	Tags               []string `json:"tags"`
	ProblemType        string   `json:"problemType"`
	TimeLimitSec       float64  `json:"timeLimitSec"`
	MemoryLimitMb      int      `json:"memoryLimitMb"`
	CheckerFile        string   `json:"checkerFile,omitempty"`
	CheckerLanguage    string   `json:"checkerLanguage,omitempty"`
	CheckType          string   `json:"checkType,omitempty"`
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
}
type RunnerTemplateData struct {
	UserSourceFile        string // 用戶原始碼文件的基本名稱 (例如 "solution.py", "user_code.cpp")
	UserModuleOrClassName string // 例如 "solution" (用於 Python import)
}
