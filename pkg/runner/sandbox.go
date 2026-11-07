package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const MaxOutputBytes = 1024 * 1024 // 1 MB buffer limit

type LocalSandboxRunner struct {
	BaseTempDir string
}

func NewLocalSandboxRunner() *LocalSandboxRunner {
	return &LocalSandboxRunner{
		BaseTempDir: os.TempDir(),
	}
}

func (r *LocalSandboxRunner) Execute(req ExecutionRequest) (*ExecutionResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp(r.BaseTempDir, "coderunner-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create sandbox environment: %w", err)
	}
	defer os.RemoveAll(tempDir)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout(req.TimeoutMs))
	defer cancel()

	start := time.Now()
	res := &ExecutionResult{}

	var cmd *exec.Cmd

	switch req.Language {
	case LangPython:
		sourcePath := filepath.Join(tempDir, "main.py")
		if err := os.WriteFile(sourcePath, []byte(req.Code), 0600); err != nil {
			return nil, err
		}
		cmd = exec.CommandContext(ctx, "python3", sourcePath)

	case LangJavaScript:
		sourcePath := filepath.Join(tempDir, "index.js")
		if err := os.WriteFile(sourcePath, []byte(req.Code), 0600); err != nil {
			return nil, err
		}
		cmd = exec.CommandContext(ctx, "node", sourcePath)

	case LangBash:
		sourcePath := filepath.Join(tempDir, "script.sh")
		if err := os.WriteFile(sourcePath, []byte(req.Code), 0700); err != nil {
			return nil, err
		}
		cmd = exec.CommandContext(ctx, "bash", sourcePath)

	case LangGo:
		sourcePath := filepath.Join(tempDir, "main.go")
		if err := os.WriteFile(sourcePath, []byte(req.Code), 0600); err != nil {
			return nil, err
		}
		cmd = exec.CommandContext(ctx, "go", "run", sourcePath)

	case LangCPP:
		sourcePath := filepath.Join(tempDir, "main.cpp")
		binaryPath := filepath.Join(tempDir, "prog")
		if err := os.WriteFile(sourcePath, []byte(req.Code), 0600); err != nil {
			return nil, err
		}
		compileCmd := exec.CommandContext(ctx, "clang++", "-O2", "-std=c++17", sourcePath, "-o", binaryPath)
		var compileErr bytes.Buffer
		compileCmd.Stderr = &compileErr
		if err := compileCmd.Run(); err != nil {
			res.Stderr = compileErr.String()
			res.ExitCode = 1
			res.DurationMs = time.Since(start).Milliseconds()
			return res, nil
		}
		cmd = exec.CommandContext(ctx, binaryPath)

	default:
		return nil, ErrUnsupportedLanguage
	}

	cmd.Dir = tempDir
	if req.Stdin != "" {
		cmd.Stdin = strings.NewReader(req.Stdin)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.LimitReader(&stdoutBuf, MaxOutputBytes).(io.Reader)
	cmd.Stderr = io.LimitReader(&stderrBuf, MaxOutputBytes).(io.Reader)

	// Direct buffers to capture output
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()
	res.DurationMs = time.Since(start).Milliseconds()

	res.Stdout = stdoutBuf.String()
	res.Stderr = stderrBuf.String()

	if ctx.Err() == context.DeadlineExceeded {
		res.TimedOut = true
		res.ExitCode = -1
		res.Error = ErrExecutionTimedOut.Error()
		return res, nil
	}

	if execErr != nil {
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			res.ExitCode = exitErr.ExitCode()
		} else {
			res.ExitCode = 1
			res.Error = execErr.Error()
		}
	} else {
		res.ExitCode = 0
	}

	return res, nil
}
