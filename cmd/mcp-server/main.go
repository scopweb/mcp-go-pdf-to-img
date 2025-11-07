package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	localmcp "github.com/tu-usuario/pdf2img/mcp"
)

func main() {
	log.Println("🚀 Starting PDF2IMG MCP Server")

	// Create MCP server using mark3labs SDK
	s := server.NewMCPServer(
		"pdf2img",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Create our local MCP server for tool implementations
	localServer, err := localmcp.NewMCPServer()
	if err != nil {
		log.Fatalf("Failed to create local MCP server: %v", err)
	}
	defer localServer.Close()

	// Register tools
	if err := registerTools(s, localServer); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	log.Println("✅ Server ready - Waiting for connections...")

	// Start the stdio server using mark3labs SDK
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// registerTools registers all PDF2IMG tools with the MCP server
func registerTools(s *server.MCPServer, localServer *localmcp.MCPServer) error {
	// Get available tools from local server
	localTools := localServer.GetTools()

	// Register pdf_to_images tool
	pdfToImagesTool := mcp.NewTool("pdf_to_images",
		mcp.WithDescription("Convert PDF pages to PNG or JPG images"),
		mcp.WithString("pdf_path", mcp.Required(), mcp.Description("Path to the PDF file to convert")),
		mcp.WithString("output_dir", mcp.Required(), mcp.Description("Directory where images will be saved")),
		mcp.WithString("format", mcp.Description("Output format: 'png' or 'jpg' (default: png)")),
		mcp.WithNumber("dpi", mcp.Description("DPI for rendering (default: 150)")),
		mcp.WithNumber("start_page", mcp.Description("Start page number (1-indexed, 0 for first page)")),
		mcp.WithNumber("end_page", mcp.Description("End page number (1-indexed, 0 for last page)")),
		mcp.WithString("prefix", mcp.Description("Prefix for output filenames (default: page_)")),
	)

	s.AddTool(pdfToImagesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pdfPath, err := request.RequireString("pdf_path")
		if err != nil {
			return mcp.NewToolResultError("pdf_path is required"), nil
		}

		outputDir, err := request.RequireString("output_dir")
		if err != nil {
			return mcp.NewToolResultError("output_dir is required"), nil
		}

		// Extract optional parameters from Arguments map
		format := "png"
		dpi := 150.0
		startPage := 0
		endPage := 0
		prefix := "page_"

		if args, ok := request.Params.Arguments.(map[string]interface{}); ok {
			if f, ok := args["format"].(string); ok && f != "" {
				format = f
			}
			if d, ok := args["dpi"].(float64); ok && d > 0 {
				dpi = d
			}
			if sp, ok := args["start_page"].(float64); ok {
				startPage = int(sp)
			}
			if ep, ok := args["end_page"].(float64); ok {
				endPage = int(ep)
			}
			if p, ok := args["prefix"].(string); ok && p != "" {
				prefix = p
			}
		}

		// Execute tool using local server
		result, err := localServer.ExecuteTool("pdf_to_images", []byte(fmt.Sprintf(`{
			"pdf_path": "%s",
			"output_dir": "%s",
			"format": "%s",
			"dpi": %f,
			"start_page": %d,
			"end_page": %d,
			"prefix": "%s"
		}`, pdfPath, outputDir, format, dpi, startPage, endPage, prefix)))

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Conversion error: %v", err)), nil
		}

		return mcp.NewToolResultText(result.Content), nil
	})

	// Register pdf_info tool
	pdfInfoTool := mcp.NewTool("pdf_info",
		mcp.WithDescription("Get information about a PDF file (pages, size, dimensions)"),
		mcp.WithString("pdf_path", mcp.Required(), mcp.Description("Path to the PDF file")),
	)

	s.AddTool(pdfInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pdfPath, err := request.RequireString("pdf_path")
		if err != nil {
			return mcp.NewToolResultError("pdf_path is required"), nil
		}

		// Execute tool using local server
		result, err := localServer.ExecuteTool("pdf_info", []byte(fmt.Sprintf(`{"pdf_path": "%s"}`, pdfPath)))

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText(result.Content), nil
	})

	log.Printf("📚 Registered %d tools", len(localTools))
	for _, tool := range localTools {
		log.Printf("   - %s: %s", tool.Name, tool.Description)
	}

	return nil
}
