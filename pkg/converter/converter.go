package converter

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

// Converter manages PDF to image conversion
type Converter struct {
	pool     pdfium.Pool
	instance pdfium.Pdfium
}

// ConvertOptions specifies conversion parameters
type ConvertOptions struct {
	InputPath    string  // Path to PDF file
	OutputDir    string  // Output directory
	Format       string  // "png" or "jpg"
	DPI          float64 // DPI for rendering (default 150)
	StartPage    int     // Start page (1-indexed, 0 = all)
	EndPage      int     // End page (1-indexed, 0 = all)
	Prefix       string  // Prefix for output files
	RetryFailed  bool    // Retry failed pages with reduced DPI (default false)
	MaxPoolSize  int     // Max PDFium instances in pool (default 2, increase for large PDFs)
	RefreshEvery int     // Refresh PDFium instance every N pages (0 = disable, default 50)
}

// ConvertResult contains conversion results
type ConvertResult struct {
	TotalPages  int
	Successful  int
	Failed      int
	OutputFiles []string
	Errors      []string
	WarningPages []int // Pages with unreachable errors (may need manual inspection)
}

// New creates a new Converter instance using WebAssembly PDFium
func New() (*Converter, error) {
	return NewWithPoolSize(2)
}

// NewWithPoolSize creates a new Converter with a specific pool size
func NewWithPoolSize(poolSize int) (*Converter, error) {
	if poolSize < 1 {
		poolSize = 2
	}

	// Initialize PDFium pool (WebAssembly - pure Go)
	// Larger pool helps prevent memory issues with large PDFs
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1,
		MaxIdle:  poolSize,
		MaxTotal: poolSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PDFium pool: %w", err)
	}

	// Get an instance from the pool
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get PDFium instance: %w", err)
	}

	return &Converter{
		pool:     pool,
		instance: instance,
	}, nil
}

// Close releases resources
func (c *Converter) Close() error {
	// Close the instance (returns it to the pool)
	if c.instance != nil {
		c.instance.Close()
	}
	// Close the pool
	if c.pool != nil {
		c.pool.Close()
	}
	return nil
}

// refreshInstance returns the current instance to the pool and gets a fresh one
// This helps prevent memory accumulation during processing of large PDFs
func (c *Converter) refreshInstance() error {
	if c.instance != nil {
		c.instance.Close()
	}

	// Get a fresh instance from the pool
	instance, err := c.pool.GetInstance(time.Second * 30)
	if err != nil {
		return fmt.Errorf("failed to get fresh PDFium instance: %w", err)
	}

	c.instance = instance
	return nil
}

// GetPDFInfo returns information about a PDF file
func (c *Converter) GetPDFInfo(pdfPath string) (map[string]interface{}, error) {
	if _, err := os.Stat(pdfPath); err != nil {
		return nil, fmt.Errorf("PDF file not found: %w", err)
	}

	// Read PDF file into memory (WebAssembly requires bytes, not file paths)
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Open document
	doc, err := c.instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer c.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Get page count
	pageCountReq := &requests.FPDF_GetPageCount{
		Document: doc.Document,
	}
	pageCountRes, err := c.instance.FPDF_GetPageCount(pageCountReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	info := map[string]interface{}{
		"file":       filepath.Base(pdfPath),
		"pages":      pageCountRes.PageCount,
		"file_size":  getFileSize(pdfPath),
	}

	// Try to get first page dimensions
	if pageCountRes.PageCount > 0 {
		pageReq := &requests.FPDFPage_GetMediaBox{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    0,
				},
			},
		}
		pageRes, err := c.instance.FPDFPage_GetMediaBox(pageReq)
		if err == nil {
			info["width"] = pageRes.Right
			info["height"] = pageRes.Bottom
		}
	}

	return info, nil
}

// Convert renders PDF pages to images
func (c *Converter) Convert(opts *ConvertOptions) (*ConvertResult, error) {
	// Validate options
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	// Read PDF file into memory (WebAssembly requires bytes, not file paths)
	pdfBytes, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Open document
	doc, err := c.instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer c.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Get page count
	pageCountReq := &requests.FPDF_GetPageCount{
		Document: doc.Document,
	}
	pageCountRes, err := c.instance.FPDF_GetPageCount(pageCountReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	pageCount := pageCountRes.PageCount

	// Determine page range
	startPage, endPage := determinePageRange(opts.StartPage, opts.EndPage, pageCount)

	result := &ConvertResult{
		TotalPages:  pageCount,
		OutputFiles: []string{},
		Errors:      []string{},
		WarningPages: []int{},
	}

	// Set DPI
	dpi := opts.DPI
	if dpi <= 0 {
		dpi = 150
	}

	// Set refresh interval for large PDFs
	refreshEvery := opts.RefreshEvery
	if refreshEvery <= 0 {
		refreshEvery = 50 // Default: refresh every 50 pages
	}

	// Set pool size for converter
	poolSize := opts.MaxPoolSize
	if poolSize <= 0 {
		poolSize = 2
	}

	// Track failed pages for retry
	failedPages := make(map[int]string)
	pagesProcessed := 0

	// Note: Instance refresh has been disabled because it closes the document.
	// Instead, we rely on the pool with larger MaxTotal to handle memory better.
	// Users can configure pool-size for their hardware needs.
	_ = refreshEvery  // Keep the parameter but don't use it
	_ = poolSize      // Keep the parameter but don't use it

	// Render each page
	for pageNum := startPage; pageNum <= endPage; pageNum++ {
		// Render page to image
		pageRender, err := c.instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI: int(dpi),
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    pageNum - 1,
				},
			},
		})

		if err != nil {
			result.Failed++
			failedPages[pageNum] = err.Error()
			result.Errors = append(result.Errors, fmt.Sprintf("Page %d: %v", pageNum, err))

			// Mark pages with WASM errors for potential retry
			if strings.Contains(err.Error(), "unreachable") || strings.Contains(err.Error(), "wasm") {
				result.WarningPages = append(result.WarningPages, pageNum)
			}
			continue
		}

		// Clean up render resources
		if pageRender != nil && pageRender.Cleanup != nil {
			defer pageRender.Cleanup()
		}

		// Get the image from result
		if pageRender == nil || pageRender.Result.Image == nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Page %d: no image generated", pageNum))
			continue
		}

		// Save image
		outputPath := filepath.Join(
			opts.OutputDir,
			fmt.Sprintf("%s%04d.%s", opts.Prefix, pageNum, opts.Format),
		)

		if err := saveImage(pageRender.Result.Image, outputPath, opts.Format); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Page %d save: %v", pageNum, err))
			continue
		}

		result.Successful++
		result.OutputFiles = append(result.OutputFiles, outputPath)
		pagesProcessed++
	}

	// Retry failed pages with reduced DPI if requested
	if opts.RetryFailed && len(failedPages) > 0 && dpi > 72 {
		retryDPI := dpi * 0.75 // Reduce DPI by 25%
		for pageNum, _ := range failedPages {
			// Try rendering with reduced DPI
			pageRender, err := c.instance.RenderPageInDPI(&requests.RenderPageInDPI{
				DPI: int(retryDPI),
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc.Document,
						Index:    pageNum - 1,
					},
				},
			})

			if err == nil && pageRender != nil && pageRender.Result.Image != nil {
				outputPath := filepath.Join(
					opts.OutputDir,
					fmt.Sprintf("%s%04d.%s", opts.Prefix, pageNum, opts.Format),
				)

				if err := saveImage(pageRender.Result.Image, outputPath, opts.Format); err == nil {
					result.Failed--
					result.Successful++
					result.OutputFiles = append(result.OutputFiles, outputPath)
					// Remove from errors list
					newErrors := []string{}
					for _, e := range result.Errors {
						if !strings.Contains(e, fmt.Sprintf("Page %d", pageNum)) {
							newErrors = append(newErrors, e)
						}
					}
					result.Errors = newErrors
				}

				// Clean up resources
				if pageRender.Cleanup != nil {
					pageRender.Cleanup()
				}
			}
		}
	}

	return result, nil
}

// Helper functions

func validateOptions(opts *ConvertOptions) error {
	if opts == nil {
		return fmt.Errorf("options cannot be nil")
	}

	if opts.InputPath == "" {
		return fmt.Errorf("input path is required")
	}

	if opts.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	if opts.Format == "" {
		opts.Format = "png"
	}

	opts.Format = strings.ToLower(opts.Format)
	if opts.Format != "png" && opts.Format != "jpg" && opts.Format != "jpeg" {
		return fmt.Errorf("format must be 'png' or 'jpg'")
	}

	if opts.Prefix == "" {
		opts.Prefix = "page_"
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if _, err := os.Stat(opts.InputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
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

func saveImage(img image.Image, path string, format string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	switch format {
	case "png":
		if err := png.Encode(file, img); err != nil {
			return fmt.Errorf("failed to encode PNG: %w", err)
		}
	case "jpg", "jpeg":
		if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("failed to encode JPEG: %w", err)
		}
	}

	return nil
}

func getFileSize(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "unknown"
	}

	size := info.Size()
	units := []string{"B", "KB", "MB", "GB"}
	divisor := 1.0

	for i := 0; i < len(units) && size > 1024; i++ {
		size = size / 1024
		divisor = divisor * 1024
	}

	return fmt.Sprintf("%.2f %s", float64(info.Size())/divisor, units[int(divisor/1024)])
}
