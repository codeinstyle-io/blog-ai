package cmd

import (
	"os"
	"testing"

	"github.com/captain-corp/captain/utils"
)

func TestGetValidInput(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	tests := []struct {
		name      string
		inputs    []string
		validator func(string) error
		want      string
	}{
		{
			name:      "Valid input first try",
			inputs:    []string{"test\n"},
			validator: utils.ValidateFirstName,
			want:      "test",
		},
		{
			name:      "Valid input after invalid",
			inputs:    []string{"123\n", "test\n"},
			validator: utils.ValidateFirstName,
			want:      "test",
		},
		{
			name:      "Valid email",
			inputs:    []string{"test@example.com\n"},
			validator: utils.ValidateEmail,
			want:      "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe and make it stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write all inputs
			errCh := make(chan error, 1)
			go func() {
				for _, input := range tt.inputs {
					if _, err := w.Write([]byte(input)); err != nil {
						errCh <- err
						return
					}
				}
				w.Close()
				errCh <- nil
			}()

			got := getValidInput("Test: ", tt.validator)
			if err := <-errCh; err != nil {
				t.Fatalf("Failed to write test input: %v", err)
			}

			if got != tt.want {
				t.Errorf("getValidInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
