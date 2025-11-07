package converter

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidateOptions tests the options validation
func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConvertOptions
		wantErr bool
	}{
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "missing input path",
			opts: &ConvertOptions{
				OutputDir: "/tmp",
			},
			wantErr: true,
		},
		{
			name: "missing output directory",
			opts: &ConvertOptions{
				InputPath: "test.pdf",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			opts: &ConvertOptions{
				InputPath: "test.pdf",
				OutputDir: "/tmp",
				Format:    "bmp",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptions(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDeterminePageRange tests page range determination
func TestDeterminePageRange(t *testing.T) {
	tests := []struct {
		name      string
		startPage int
		endPage   int
		totalPage int
		want      [2]int
	}{
		{
			name:      "all pages",
			startPage: 0,
			endPage:   0,
			totalPage: 10,
			want:      [2]int{1, 10},
		},
		{
			name:      "specific range",
			startPage: 2,
			endPage:   5,
			totalPage: 10,
			want:      [2]int{2, 5},
		},
		{
			name:      "end page beyond total",
			startPage: 5,
			endPage:   20,
			totalPage: 10,
			want:      [2]int{5, 10},
		},
		{
			name:      "start page greater than end",
			startPage: 10,
			endPage:   5,
			totalPage: 15,
			want:      [2]int{5, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := determinePageRange(tt.startPage, tt.endPage, tt.totalPage)
			if start != tt.want[0] || end != tt.want[1] {
				t.Errorf("determinePageRange() = [%d, %d], want [%d, %d]", start, end, tt.want[0], tt.want[1])
			}
		})
	}
}

// TestFileSize tests file size formatting
func TestGetFileSize(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Create file with 1024 bytes
	err := os.WriteFile(testFile, make([]byte, 1024), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	size := getFileSize(testFile)
	if size == "unknown" {
		t.Errorf("getFileSize() returned 'unknown'")
	}
}
