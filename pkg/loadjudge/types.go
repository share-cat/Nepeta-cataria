package loadjudge

const (
	SampleInputFile  = "sample.in"
	SampleOutputFile = "sample.out"
	ConfigFileName   = "config.yaml"
	TestDataYamlFile = "testdata.yaml"
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

type TestDataConfig struct {
	Sample struct {
		Input       []map[string]interface{} `yaml:"input"`
		Output      interface{}              `yaml:"output"`
		Description string                   `yaml:"description"`
	} `yaml:"sample"`

	TestCases []TestCase `yaml:"cases"`

	TestConfig struct {
		TotalPoints        int    `yaml:"total_points"`
		TimeLimitPerCase   string `yaml:"time_limit_per_case"`
		MemoryLimitPerCase string `yaml:"memory_limit_per_case"`

		InputFormat []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
		} `yaml:"input_format"`

		OutputFormat struct {
			Type string `yaml:"type"`
		} `yaml:"output_format"`
	} `yaml:"config"`
}

type TestCase struct {
	Name   string                   `yaml:"name" json:"name"`
	Input  []map[string]interface{} `yaml:"input" json:"input"`
	Output interface{}              `yaml:"output" json:"output"`
	Points int                      `yaml:"points" json:"points"`
}
