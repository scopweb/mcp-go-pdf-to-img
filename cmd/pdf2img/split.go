package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tu-usuario/pdf2img/pkg/splitter"
)

var (
	splitInputFile  string
	splitOutputFile string
	splitStartPage  int
	splitEndPage    int
	splitVerbose    bool
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Extract pages from a PDF into a new PDF file",
	Long:  "Extract a range of pages from a PDF document and save them to a new PDF file.",
	RunE:  runSplit,
}

func init() {
	splitCmd.Flags().StringVarP(&splitInputFile, "input", "i", "", "Input PDF file (required)")
	splitCmd.Flags().StringVarP(&splitOutputFile, "output", "o", "", "Output PDF file (required)")
	splitCmd.Flags().IntVar(&splitStartPage, "start", 0, "Start page number (1-indexed, 0 for first)")
	splitCmd.Flags().IntVar(&splitEndPage, "end", 0, "End page number (1-indexed, 0 for last)")
	splitCmd.Flags().BoolVarP(&splitVerbose, "verbose", "v", false, "Verbose output")

	splitCmd.MarkFlagRequired("input")
	splitCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(splitCmd)
}

func runSplit(cmd *cobra.Command, args []string) error {
	split := splitter.New()

	opts := &splitter.SplitOptions{
		InputPath:  splitInputFile,
		OutputPath: splitOutputFile,
		StartPage:  splitStartPage,
		EndPage:    splitEndPage,
	}

	if splitVerbose {
		log.Printf("Extracting pages from: %s\n", splitInputFile)
		log.Printf("Output file: %s\n", splitOutputFile)
	}

	result, err := split.Split(opts)
	if err != nil {
		return fmt.Errorf("split failed: %w", err)
	}

	// Print results
	fmt.Printf("\n✓ Split Complete\n")
	fmt.Printf("Total pages in PDF: %d\n", result.TotalPages)
	fmt.Printf("Extracted pages: %d\n", result.ExtractedPages)
	fmt.Printf("Output saved to: %s\n", result.OutputPath)

	return nil
}
