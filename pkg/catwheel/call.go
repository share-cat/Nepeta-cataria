package catwheel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
	DefaultBaseURL = "http://localhost:3200"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type ExecuteRequest struct {
	Code        string `json:"code"`
	Language    string `json:"language"`
	Timeout     int    `json:"timeout"`
	CPULimit    int64  `json:"cpu_limit"`
	MemoryLimit int64  `json:"memory_limit"`
}

type ExecuteResponse struct {
	ExecutionID string          `json:"execution_id"`
	Result      ExecutionResult `json:"result"`
	Success     bool            `json:"success"`
}

type ExecutionResult struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	Stdout        string    `json:"stdout"`
	Stderr        string    `json:"stderr"`
	ExitCode      int       `json:"exit_code"`
	ExecutionTime int64     `json:"execution_time"`
	CPUTime       int64     `json:"cpu_time"`
	MemoryUsed    int64     `json:"memory_used"`
	ErrorMessage  *string   `json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   time.Time `json:"completed_at"`
}

type ListExecutionsResponse struct {
	Executions []ExecutionResult `json:"executions"`
	Success    bool              `json:"success"`
	Total      int               `json:"total"`
}

type GetExecutionResponse struct {
	Result  ExecutionResult `json:"result"`
	Success bool            `json:"success"`
}

type CleanupResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp bool   `json:"timestamp"`
	Success   bool   `json:"version"`
}

type SystemInfoResponse struct {
	System SystemInfo `json:"system"`
}

type SystemInfo struct {
	CgroupV2Available bool               `json:"cgroup_v2_available"`
	UID               uint32             `json:"uid"`
	ServiceVersion    string             `json:"service_version"`
	Runtimes          map[string]Runtime `json:"runtimes"`
}

type Runtime struct {
	Available bool   `json:"available"`
	Version   string `json:"version,omitempty"`
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (c *Client) ExecuteCode(code, language string, timeout int, cpuLimit, memoryLimit int64) (*ExecuteResponse, error) {
	encodedCode := base64.StdEncoding.EncodeToString([]byte(code))

	request := ExecuteRequest{
		Code:        encodedCode,
		Language:    language,
		Timeout:     timeout,
		CPULimit:    cpuLimit,
		MemoryLimit: memoryLimit,
	}

	var response ExecuteResponse
	err := c.makeRequest("POST", "/execute", request, &response)
	if err != nil {
		return nil, fmt.Errorf("execute code failed: %w", err)
	}

	return &response, nil
}

func (c *Client) ListExecutions() (*ListExecutionsResponse, error) {
	var response ListExecutionsResponse
	err := c.makeRequest("GET", "/executions", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("list executions failed: %w", err)
	}

	return &response, nil
}

func (c *Client) GetExecution(id string) (*GetExecutionResponse, error) {
	var response GetExecutionResponse
	err := c.makeRequest("GET", fmt.Sprintf("/executions/%s", id), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("get execution failed: %w", err)
	}

	return &response, nil
}

func (c *Client) CleanupExecution(id string) (*CleanupResponse, error) {
	var response CleanupResponse
	err := c.makeRequest("DELETE", fmt.Sprintf("/executions/%s", id), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("cleanup execution failed: %w", err)
	}

	return &response, nil
}

func (c *Client) HealthCheck() (*HealthResponse, error) {
	var response HealthResponse
	err := c.makeRequest("GET", "/health", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	return &response, nil
}

func (c *Client) SystemInfo() (*SystemInfoResponse, error) {
	var response SystemInfoResponse
	err := c.makeRequest("GET", "/info", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("system info failed: %w", err)
	}

	return &response, nil
}

func (c *Client) makeRequest(method, endpoint string, requestBody interface{}, response interface{}) error {
	url := c.BaseURL + endpoint

	var body io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

func (s *SystemInfo) IsRuntimeAvailable(runtime string) bool {
	if r, exists := s.Runtimes[runtime]; exists {
		return r.Available
	}
	return false
}

func (s *SystemInfo) GetRuntimeVersion(runtime string) string {
	if r, exists := s.Runtimes[runtime]; exists && r.Available {
		return r.Version
	}
	return ""
}

func (s *SystemInfo) GetAvailableRuntimes() []string {
	var available []string
	for name, runtime := range s.Runtimes {
		if runtime.Available {
			available = append(available, name)
		}
	}
	return available
}

func (r *ExecutionResult) GetStdout() string {
	return strings.TrimSuffix(r.Stdout, "\n")
}

func (r *ExecutionResult) IsSuccess() bool {
	return r.Status == StatusSuccess && r.ExitCode == 0
}

func (r *ExecutionResult) HasError() bool {
	return r.ErrorMessage != nil && *r.ErrorMessage != ""
}

func (r *ExecutionResult) GetErrorMessage() string {
	if r.ErrorMessage != nil {
		return *r.ErrorMessage
	}
	return ""
}

func (r *ExecutionResult) Duration() time.Duration {
	return time.Duration(r.ExecutionTime) * time.Microsecond
}

func (r *ExecutionResult) CPUDuration() time.Duration {
	return time.Duration(r.CPUTime) * time.Microsecond
}

func (r *ExecutionResult) MemoryUsedMB() float64 {
	return float64(r.MemoryUsed) / 1024 / 1024
}
