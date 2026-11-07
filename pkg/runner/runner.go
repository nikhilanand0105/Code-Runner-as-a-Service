package runner

import (
	"errors"
	"time"
)

type SupportedLanguage string

const (
	LangPython     SupportedLanguage = "python"
	LangJavaScript SupportedLanguage = "javascript"
	LangGo         SupportedLanguage = "go"
	LangCPP        SupportedLanguage = "cpp"
	LangBash       SupportedLanguage = "bash"
)

var (
	ErrUnsupportedLanguage = errors.New("unsupported programming language")
	ErrEmptyCode           = errors.New("code snippet cannot be empty")
	ErrExecutionTimedOut   = errors.New("execution time limit exceeded")
)

type ExecutionRequest struct {
	Language  SupportedLanguage `json:"language"`
	Code      string            `json:"code"`
	Stdin     string            `json:"stdin,omitempty"`
	TimeoutMs int               `json:"timeout_ms,omitempty"`
}

type ExecutionResult struct {
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	ExitCode   int           `json:"exit_code"`
	DurationMs int64         `json:"duration_ms"`
	TimedOut   bool          `json:"timed_out"`
	Error      string        `json:"error,omitempty"`
}

type Runner interface {
	Execute(req ExecutionRequest) (*ExecutionResult, error)
}

func (req *ExecutionRequest) Validate() error {
	if req.Code == "" {
		return ErrEmptyCode
	}
	switch req.Language {
	case LangPython, LangJavaScript, LangGo, LangCPP, LangBash:
		// Valid
	default:
		return ErrUnsupportedLanguage
	}

	if req.TimeoutMs <= 0 || req.TimeoutMs > 30000 {
		req.TimeoutMs = 5000 // Default 5 seconds
	}
	return nil
}

func DefaultTimeout(ms int) time.Duration {
	if ms <= 0 {
		return 5 * time.Second
	}
	return time.Duration(ms) * time.Millisecond
}
