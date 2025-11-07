package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/tu-usuario/pdf2img/pkg/converter"
)

// MCPServer implements the Model Context Protocol server
type MCPServer struct {
	converter *converter.Converter
}

// Tool represents an available MCP tool
type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() (*MCPServer, error) {
	conv, err := converter.New()
	if err != nil {
		return nil, err
	}

	return &MCPServer{
		converter: conv,
	}, nil
}

// GetTools returns available tools
func (s *MCPServer) GetTools() []Tool {
	return []Tool{
		{
			Name:        "pdf_to_images",
			Description: "Convert PDF pages to PNG or JPG images",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pdf_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the PDF file to convert",
					},
					"output_dir": map[string]interface{}{
						"type":        "string",
						"description": "Directory where images will be saved",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"description": "Output format: 'png' or 'jpg' (default: png)",
						"enum":        []string{"png", "jpg"},
					},
					"dpi": map[string]interface{}{
						"type":        "number",
						"description": "DPI for rendering (default: 150)",
					},
					"start_page": map[string]interface{}{
						"type":        "integer",
						"description": "Start page number (1-indexed, 0 for first page)",
					},
					"end_page": map[string]interface{}{
						"type":        "integer",
						"description": "End page number (1-indexed, 0 for last page)",
					},
					"prefix": map[string]interface{}{
						"type":        "string",
						"description": "Prefix for output filenames (default: page_)",
					},
				},
				"required": []string{"pdf_path", "output_dir"},
			},
		},
		{
			Name:        "pdf_info",
			Description: "Get information about a PDF file",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pdf_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the PDF file",
					},
				},
				"required": []string{"pdf_path"},
			},
		},
	}
}

// ExecuteTool executes a tool with the given input
func (s *MCPServer) ExecuteTool(toolName string, input json.RawMessage) (ToolResult, error) {
	switch toolName {
	case "pdf_to_images":
		return s.handlePDFToImages(input)
	case "pdf_info":
		return s.handlePDFInfo(input)
	default:
		return ToolResult{}, fmt.Errorf("unknown tool: %s", toolName)
	}
}

// Private tool handlers

func (s *MCPServer) handlePDFToImages(input json.RawMessage) (ToolResult, error) {
	var req struct {
		PDFPath   string  `json:"pdf_path"`
		OutputDir string  `json:"output_dir"`
		Format    string  `json:"format"`
		DPI       float64 `json:"dpi"`
		StartPage int     `json:"start_page"`
		EndPage   int     `json:"end_page"`
		Prefix    string  `json:"prefix"`
	}

	if err := json.Unmarshal(input, &req); err != nil {
		return ToolResult{}, fmt.Errorf("invalid input: %w", err)
	}

	if req.Format == "" {
		req.Format = "png"
	}
	if req.DPI == 0 {
		req.DPI = 150
	}
	if req.Prefix == "" {
		req.Prefix = "page_"
	}

	opts := &converter.ConvertOptions{
		InputPath:  req.PDFPath,
		OutputDir:  req.OutputDir,
		Format:     req.Format,
		DPI:        req.DPI,
		StartPage:  req.StartPage,
		EndPage:    req.EndPage,
		Prefix:     req.Prefix,
	}

	result, err := s.converter.Convert(opts)
	if err != nil {
		return ToolResult{}, err
	}

	response := map[string]interface{}{
		"total_pages": result.TotalPages,
		"successful":  result.Successful,
		"failed":      result.Failed,
		"files":       result.OutputFiles,
	}

	if len(result.Errors) > 0 {
		response["errors"] = result.Errors
	}

	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return ToolResult{
		Type:    "text",
		Content: string(responseJSON),
	}, nil
}

func (s *MCPServer) handlePDFInfo(input json.RawMessage) (ToolResult, error) {
	var req struct {
		PDFPath string `json:"pdf_path"`
	}

	if err := json.Unmarshal(input, &req); err != nil {
		return ToolResult{}, fmt.Errorf("invalid input: %w", err)
	}

	info, err := s.converter.GetPDFInfo(req.PDFPath)
	if err != nil {
		return ToolResult{}, err
	}

	responseJSON, _ := json.MarshalIndent(info, "", "  ")
	return ToolResult{
		Type:    "text",
		Content: string(responseJSON),
	}, nil
}

// Close closes the server and releases resources
func (s *MCPServer) Close() error {
	if s.converter != nil {
		return s.converter.Close()
	}
	return nil
}
