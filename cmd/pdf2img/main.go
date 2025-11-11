package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tu-usuario/pdf2img/pkg/converter"
)

var (
	inputFile    string
	outputDir    string
	format       string
	dpi          float64
	startPage    int
	endPage      int
	prefix       string
	verbose      bool
	retryFailed  bool
	maxPoolSize  int
	refreshEvery int
)

var rootCmd = &cobra.Command{
	Use:   "pdf2img",
	Short: "Convert PDF pages to images (PNG or JPG)",
	Long:  "A command-line tool to render PDF pages as images with configurable DPI and format.",
	RunE:  runConvert,
}

var infoCmd = &cobra.Command{
	Use:   "info <pdf-file>",
	Short: "Get information about a PDF file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInfo(args[0])
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input PDF file (required)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory (default: current directory)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "png", "Output format: png or jpg (default: png)")
	rootCmd.Flags().Float64VarP(&dpi, "dpi", "d", 150, "DPI for rendering (default: 150)")
	rootCmd.Flags().IntVar(&startPage, "start", 0, "Start page number (1-indexed, 0 for first)")
	rootCmd.Flags().IntVar(&endPage, "end", 0, "End page number (1-indexed, 0 for last)")
	rootCmd.Flags().StringVar(&prefix, "prefix", "page_", "Prefix for output files (default: page_)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVar(&retryFailed, "retry", false, "Retry failed pages with reduced DPI")
	rootCmd.Flags().IntVar(&maxPoolSize, "pool-size", 2, "Max PDFium instances in pool (default: 2, increase for large PDFs)")
	rootCmd.Flags().IntVar(&refreshEvery, "refresh-every", 50, "Refresh PDFium instance every N pages (default: 50, 0 to disable)")

	rootCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(infoCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runConvert(cmd *cobra.Command, args []string) error {
	// Initialize converter
	conv, err := converter.New()
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	defer conv.Close()

	// Prepare options
	opts := &converter.ConvertOptions{
		InputPath:    inputFile,
		OutputDir:    outputDir,
		Format:       format,
		DPI:          dpi,
		StartPage:    startPage,
		EndPage:      endPage,
		Prefix:       prefix,
		RetryFailed:  retryFailed,
		MaxPoolSize:  maxPoolSize,
		RefreshEvery: refreshEvery,
	}

	if verbose {
		log.Printf("Starting conversion: %s\n", inputFile)
		log.Printf("Output directory: %s\n", outputDir)
		log.Printf("Format: %s, DPI: %.0f\n", format, dpi)
	}

	// Convert
	result, err := conv.Convert(opts)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Print results
	fmt.Printf("\n✓ Conversion Complete\n")
	fmt.Printf("Total pages: %d\n", result.TotalPages)
	fmt.Printf("Successful: %d\n", result.Successful)
	fmt.Printf("Failed: %d\n", result.Failed)

	if verbose && len(result.OutputFiles) > 0 {
		fmt.Println("\nOutput files:")
		for _, file := range result.OutputFiles {
			fmt.Printf("  - %s\n", file)
		}
	}

	if len(result.WarningPages) > 0 {
		fmt.Println("\n⚠ Pages with WASM/unreachable errors (may need manual inspection):")
		for _, pageNum := range result.WarningPages {
			fmt.Printf("  - Page %d\n", pageNum)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, errMsg := range result.Errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	if retryFailed && len(result.Errors) > 0 {
		fmt.Println("\nTip: Some pages failed. Try running with --retry flag to attempt rendering with reduced DPI.")
	}

	return nil
}

func runInfo(pdfPath string) error {
	conv, err := converter.New()
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	defer conv.Close()

	info, err := conv.GetPDFInfo(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to get PDF info: %w", err)
	}

	fmt.Printf("\nPDF Information\n")
	fmt.Printf("===============\n")
	fmt.Printf("File: %s\n", info["file"])
	fmt.Printf("Pages: %d\n", info["pages"])
	fmt.Printf("Size: %s\n", info["file_size"])

	if width, ok := info["width"]; ok {
		fmt.Printf("Width: %.2f pt\n", width)
	}
	if height, ok := info["height"]; ok {
		fmt.Printf("Height: %.2f pt\n", height)
	}

	return nil
}
