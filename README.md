# PDF2IMG - Convert PDF to Images

A complete tool in Go to render PDF pages as PNG or JPG images. Includes both a CLI application and an MCP (Model Context Protocol) server.

## Features

- ✨ **High-quality rendering** using go-pdfium v1.17.2 (based on PDFium)
- 🖼️ **Multiple formats**: PNG and JPG
- 🎯 **Granular control**: Customizable DPI, configurable page ranges
- 💻 **Complete CLI**: Command-line interface with multiple options
- 🔌 **MCP Server**: Integration with Model Context Protocol
- 📊 **PDF Information**: Command to get PDF file metadata
- 🗜️ **PDF Compression**: Optimize and compress PDF files to reduce size
- 🚀 **Pure Go**: No CGO, embedded WebAssembly, single binary (~18 MB)
- 🔒 **Cross-platform**: Windows, Linux, macOS

## Installation

### Prerequisites

- Go 1.21 or higher
- **Note**: No external PDFium required, everything is embedded in the binary

### From source code

```bash
git clone <repository-url>
cd pdf2img
go mod download
go build -o pdf2img ./cmd/pdf2img
```

### Global installation

```bash
go install github.com/tu-usuario/pdf2img/cmd/pdf2img@latest
```

## Usage

### CLI - Basic conversion

```bash
# Convert all pages to PNG
pdf2img -i document.pdf -o ./output

# Convert to JPG with custom DPI
pdf2img -i document.pdf -o ./output -f jpg -d 300

# Convert only pages 1-10
pdf2img -i document.pdf -o ./output --start 1 --end 10

# With custom prefix
pdf2img -i document.pdf -o ./output --prefix img_
```

### CLI - Get PDF information

```bash
pdf2img info document.pdf
```

Output:
```
PDF Information
===============
File: document.pdf
Pages: 25
Size: 2.50 MB
Width: 612.00 pt
Height: 792.00 pt
```

### CLI Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--input` | `-i` | Path to PDF (required) | - |
| `--output` | `-o` | Output directory | `.` |
| `--format` | `-f` | Format: png or jpg | `png` |
| `--dpi` | `-d` | DPI for rendering | `150` |
| `--start` | - | Start page (1-indexed) | `0` (first) |
| `--end` | - | End page (1-indexed) | `0` (last) |
| `--prefix` | - | Prefix for output files | `page_` |
| `--verbose` | `-v` | Detailed output | `false` |

### MCP Server

The pdf2img MCP Server enables integration with Claude Desktop and other systems supporting the Model Context Protocol.

#### Configuration with Claude Desktop (Recommended)

**⚡ Quick guide**: Read [MCP_CLAUDE_DESKTOP.md](MCP_CLAUDE_DESKTOP.md) for step-by-step instructions.

**Summary**:
1. Compile: `go build -o mcp-server.exe ./cmd/mcp-server`
2. Edit `claude_desktop_config.json` (in `%APPDATA%\Claude\`)
3. Add:
```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\path\\to\\mcp-server.exe",
      "args": ["--stdio"]
    }
  }
}
```
4. Restart Claude Desktop
5. Done! Now you can use pdf2img from Claude

**Having issues?** Read [MCP_TROUBLESHOOTING.md](MCP_TROUBLESHOOTING.md) for troubleshooting common problems.

#### Programmatic usage

To use the MCP Server in Go code:

```go
package main

import (
	"log"
	"github.com/tu-usuario/pdf2img/mcp"
)

func main() {
	server, err := mcp.NewMCPServer()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	// Use server.ExecuteTool() to execute tools
	tools := server.GetTools()
	for _, tool := range tools {
		log.Printf("Available: %s", tool.Name)
	}
}
```

#### Available MCP tools

##### `pdf_to_images`

Converts PDF pages to PNG or JPG images.

```json
{
  "pdf_path": "document.pdf",
  "output_dir": "./output",
  "format": "png",
  "dpi": 150,
  "start_page": 0,
  "end_page": 0,
  "prefix": "page_"
}
```

**Note**: `start_page` and `end_page` with value 0 mean "all pages".

##### `pdf_info`

Gets PDF information (pages, size, dimensions).

```json
{
  "pdf_path": "document.pdf"
}
```

#### Operating modes

The MCP Server supports two modes:

- **stdio** (default for Claude Desktop): `mcp-server --stdio`
- **HTTP**: `mcp-server --port 8080` (for testing)

More examples in [EXAMPLES.md](EXAMPLES.md#mcp-server---ejemplos-de-integración).

## Project structure

```
pdf2img/
├── cmd/
│   └── pdf2img/
│       └── main.go          # CLI application
├── mcp/
│   └── server.go            # MCP server
├── pkg/
│   └── converter/
│       └── converter.go      # Shared conversion logic
├── go.mod                   # Dependencies
└── README.md               # This file
```

## Dependencies

- **go-pdfium v1.17.2** (Apache 2.0): PDF rendering using PDFium WebAssembly
- **cobra v1.10.1** (Apache 2.0): CLI framework
- **mark3labs/mcp-go v0.43.0** (Apache 2.0): Official MCP SDK
- **wazero v1.9.0** (Apache 2.0): WebAssembly runtime for Go (included in go-pdfium)

## Examples

### Example 1: Convert entire document to high resolution

```bash
pdf2img -i book.pdf -o ./images -f jpg -d 300
```

Creates files like: `page_0001.jpg`, `page_0002.jpg`, etc.

### Example 2: Extract only first 5 pages

```bash
pdf2img -i report.pdf -o ./thumbs --start 1 --end 5 --prefix thumb_
```

Creates: `thumb_0001.png`, `thumb_0002.png`, etc.

### Example 3: PDF Information

```bash
pdf2img info form.pdf
```

### Example 4: From Go code

```go
package main

import (
	"log"
	"github.com/tu-usuario/pdf2img/pkg/converter"
)

func main() {
	conv, err := converter.New()
	if err != nil {
		log.Fatal(err)
	}
	defer conv.Close()

	result, err := conv.Convert(&converter.ConvertOptions{
		InputPath: "document.pdf",
		OutputDir: "./output",
		Format:    "png",
		DPI:       150,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successful: %d, Failed: %d", result.Successful, result.Failed)
}
```

## Technology

### WebAssembly Implementation

This project uses **go-pdfium v1.17.2 with WebAssembly**, which provides:

- **Pure Go**: No CGO dependencies, no need for C compiler on host
- **Single Binary**: Everything embedded in one executable (~18 MB)
- **Cross-platform**: Same binary works on Windows, Linux, macOS
- **Performance**: ~2x faster than multi-threaded CGO version
- **Security**: Isolated execution in WebAssembly sandbox

### Why WebAssembly

WebAssembly was chosen over CGO because:
- ❌ CGO requires C compiler on host (problematic on different systems)
- ❌ go-pdfium v1.12.0 CGO had severe compatibility bugs
- ✅ WebAssembly is more robust and portable
- ✅ Requires no external dependencies after compilation
- ✅ Better resource isolation

## 📚 Documentation

### Public Documentation
- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 5 minutes
- **[EXAMPLES.md](EXAMPLES.md)** - Practical use cases
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Development guide
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Code structure
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - How to contribute
- **[MCP_CLAUDE_DESKTOP.md](MCP_CLAUDE_DESKTOP.md)** - Claude Desktop integration
- **[MCP_TROUBLESHOOTING.md](MCP_TROUBLESHOOTING.md)** - MCP troubleshooting
- **[CHANGELOG.md](CHANGELOG.md)** - Change history
- **[SECURITY.md](SECURITY.md)** - Security policy
- **[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)** - Code of conduct

### Internal Documentation
For maintainers and detailed technical information, see [docs/README.md](docs/README.md)

---

## Troubleshooting

### Error: "PDF file not found"
Verify the PDF path is correct:
```bash
ls -la document.pdf
```

### Error: "failed to initialize PDFium"
Make sure Go can download dependencies:
```bash
go mod download
go mod verify
```

### Images come out in low resolution
Increase the DPI value:
```bash
pdf2img -i document.pdf -o ./output -d 300  # Higher = better quality
```

### Processing is slow
Reduce DPI for faster processing:
```bash
pdf2img -i document.pdf -o ./output -d 96
```

## License

This project uses go-pdfium under Apache 2.0 license.

## Contributing

Contributions are welcome. Please:

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Roadmap

- [ ] Watermark support
- [ ] OCR of rendered images
- [ ] Parallel page processing
- [ ] Rendered image caching
- [ ] REST HTTP API
- [ ] Password-protected PDF support

## References

- [go-pdfium](https://github.com/klippa-app/go-pdfium)
- [PDFium](https://pdfium.googlesource.com/)
- [Cobra CLI](https://cobra.dev/)
- [Model Context Protocol](https://modelcontextprotocol.io/)
