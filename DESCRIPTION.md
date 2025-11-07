cat > DESCRIPTION.md << 'EOF'
# mcp-go-pdf-to-img

A high-performance Go utility for converting PDF pages to PNG/JPG images 
with Claude Desktop MCP integration.

## Overview

This project provides a pure Go implementation for batch PDF to image conversion 
with seamless integration into Claude Desktop via the Model Context Protocol (MCP).

## Key Features

- **Pure Go + WebAssembly**: No CGO dependencies, fully portable
- **Single Binary**: ~16-18 MB executable for Windows, Linux, macOS
- **CLI Tool**: Command-line interface for batch operations
- **MCP Server**: Direct integration with Claude Desktop
- **Flexible Output**: PNG/JPG with customizable DPI and page ranges
- **Production Ready**: Based on go-pdfium v1.17.2 with official MCP SDK

## Technology Stack

- Go 1.21+
- go-pdfium v1.17.2 (WebAssembly-based PDF rendering)
- mark3labs/mcp-go v0.43.0 (Official MCP SDK)
- Cobra CLI framework

## Quick Links

- [README.md](README.md) - Complete documentation
- [MCP_CLAUDE_DESKTOP.md](MCP_CLAUDE_DESKTOP.md) - Claude Desktop setup
- [EXAMPLES.md](EXAMPLES.md) - Usage examples
- [README_es.md](README_es.md) - Spanish documentation

## License

Apache 2.0
EOF