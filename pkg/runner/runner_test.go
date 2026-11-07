package runner

import (
	"testing"
)

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     ExecutionRequest
		wantErr bool
	}{
		{
			name: "Valid Python Request",
			req: ExecutionRequest{
				Language: LangPython,
				Code:     "print('hello world')",
			},
			wantErr: false,
		},
		{
			name: "Empty Code Error",
			req: ExecutionRequest{
				Language: LangPython,
				Code:     "",
			},
			wantErr: true,
		},
		{
			name: "Unsupported Language",
			req: ExecutionRequest{
				Language: "cobol",
				Code:     "DISPLAY 'HELLO'",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultTimeout(t *testing.T) {
	if DefaultTimeout(0).Seconds() != 5 {
		t.Errorf("Expected 5s default timeout, got %v", DefaultTimeout(0))
	}
	if DefaultTimeout(2000).Milliseconds() != 2000 {
		t.Errorf("Expected 2000ms timeout, got %v", DefaultTimeout(2000))
	}
}
