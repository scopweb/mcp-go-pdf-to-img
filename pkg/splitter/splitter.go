package splitter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Splitter handles PDF page extraction operations
type Splitter struct{}

// SplitOptions specifies splitting parameters
type SplitOptions struct {
	InputPath  string // Path to PDF file
	OutputPath string // Output file path
	StartPage  int    // Start page (1-indexed, 0 = first)
	EndPage    int    // End page (1-indexed, 0 = last)
}

// SplitResult contains splitting results
type SplitResult struct {
	TotalPages int
	ExtractedPages int
	OutputPath string
}

// New creates a new Splitter instance
func New() *Splitter {
	return &Splitter{}
}

// Split extracts pages from a PDF into a new PDF file
func (s *Splitter) Split(opts *SplitOptions) (*SplitResult, error) {
	// Validate options
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	// Get total page count
	pageCount, err := api.PageCountFile(opts.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	// Determine page range
	startPage, endPage := determinePageRange(opts.StartPage, opts.EndPage, pageCount)

	// Build page selection string for pdfcpu
	// pdfcpu uses 1-indexed pages
	pageSelection := fmt.Sprintf("%d-%d", startPage, endPage)

	// Extract pages
	err = api.ExtractPagesFile(opts.InputPath, opts.OutputPath, []string{pageSelection}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to extract pages: %w", err)
	}

	result := &SplitResult{
		TotalPages:     pageCount,
		ExtractedPages: endPage - startPage + 1,
		OutputPath:     opts.OutputPath,
	}

	return result, nil
}

// Helper functions

func validateOptions(opts *SplitOptions) error {
	if opts == nil {
		return fmt.Errorf("options cannot be nil")
	}

	if opts.InputPath == "" {
		return fmt.Errorf("input path is required")
	}

	if opts.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	// Validate input file exists
	if _, err := os.Stat(opts.InputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Validate output directory exists or can be created
	outputDir := filepath.Dir(opts.OutputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	return nil
}

func determinePageRange(startPage, endPage, totalPages int) (int, int) {
	if startPage == 0 {
		startPage = 1
	}
	if endPage == 0 {
		endPage = totalPages
	}

	if startPage < 1 {
		startPage = 1
	}
	if endPage > totalPages {
		endPage = totalPages
	}
	if startPage > endPage {
		startPage = endPage
	}

	return startPage, endPage
}
