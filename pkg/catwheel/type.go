package catwheel

const (
	LanguagePython     = "Python"
	LanguageJavaScript = "JavaScript"
	LanguageCpp        = "C++"
	LanguageJava       = "Java"
	LanguageGo         = "Go"
	LanguageRust       = "Rust"

	StatusSuccess = "Success"
	StatusFailed  = "Failed"
	StatusTimeout = "Timeout"
	StatusError   = "Error"
	StatusRunning = "Running"
	StatusPending = "Pending"

	DefaultCPULimit    = 10000
	DefaultMemoryLimit = 67108864
	DefaultTimeoutSec  = 5
)

type ExecuteOptions struct {
	Language    string
	Timeout     int
	CPULimit    int64
	MemoryLimit int64
}

func DefaultExecuteOptions() ExecuteOptions {
	return ExecuteOptions{
		Language:    LanguagePython,
		Timeout:     DefaultTimeoutSec,
		CPULimit:    DefaultCPULimit,
		MemoryLimit: DefaultMemoryLimit,
	}
}

func (e ExecuteOptions) WithLanguage(lang string) ExecuteOptions {
	e.Language = lang
	return e
}

func (e ExecuteOptions) WithTimeout(timeout int) ExecuteOptions {
	e.Timeout = timeout
	return e
}

func (e ExecuteOptions) WithCPULimit(limit int64) ExecuteOptions {
	e.CPULimit = limit
	return e
}

func (e ExecuteOptions) WithMemoryLimit(limit int64) ExecuteOptions {
	e.MemoryLimit = limit
	return e
}

func (c *Client) ExecuteWithOptions(code string, options ExecuteOptions) (*ExecuteResponse, error) {
	return c.ExecuteCode(code, options.Language, options.Timeout, options.CPULimit, options.MemoryLimit)
}
