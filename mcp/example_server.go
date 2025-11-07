package mcp

import (
	"encoding/json"
	"fmt"
	"log"
)

// ExampleMCPServer demonstrates how to use the MCP server
func ExampleMCPServer() {
	// Create server
	server, err := NewMCPServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Get available tools
	tools := server.GetTools()
	fmt.Println("Available Tools:")
	for _, tool := range tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// Example 1: Get PDF info
	fmt.Println("\n--- Example 1: Get PDF Info ---")
	infoInput := json.RawMessage(`{"pdf_path": "sample.pdf"}`)
	result, err := server.ExecuteTool("pdf_info", infoInput)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println(result.Content)
	}

	// Example 2: Convert PDF to images
	fmt.Println("\n--- Example 2: Convert PDF to Images ---")
	convertInput := json.RawMessage(`{
		"pdf_path": "sample.pdf",
		"output_dir": "./output",
		"format": "png",
		"dpi": 150,
		"prefix": "page_"
	}`)
	result, err = server.ExecuteTool("pdf_to_images", convertInput)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println(result.Content)
	}
}
