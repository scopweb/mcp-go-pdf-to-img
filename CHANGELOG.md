# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- [ ] OCR de imágenes renderizadas
- [ ] Procesamiento paralelo de páginas
- [ ] Caché de imágenes renderizadas
- [ ] API HTTP REST
- [ ] Soporte para PDF con contraseña
- [ ] Soporte para marcas de agua
- [ ] Docker image

---

## [1.0.0] - 2025-11-07

### Added
- Initial release of PDF2IMG
- **CLI Application** (`pdf2img`)
  - Convert PDF pages to PNG/JPG images
  - Configurable DPI (quality)
  - Select specific page ranges
  - Custom output prefixes
  - Verbose logging option
  - `info` subcommand for PDF metadata

- **MCP Server** (`mcp-server`)
  - Model Context Protocol integration
  - Tools: `pdf_to_images`, `pdf_info`
  - HTTP and stdio support
  - JSON-based tool execution

- **Core Library** (`pkg/converter`)
  - PDF rendering engine
  - Multi-format support (PNG, JPG)
  - Error handling and reporting
  - Metadata extraction

- **Documentation**
  - README with usage guide
  - QUICKSTART guide (5 minutes)
  - EXAMPLES with practical cases
  - PROJECT_STRUCTURE explanation
  - DEVELOPMENT setup guide
  - CONTRIBUTING guidelines
  - SECURITY policy
  - CODE_OF_CONDUCT

- **Testing & Quality**
  - 19 Security tests (100% passing)
  - Path traversal protection (CWE-22)
  - Command injection protection (CWE-78)
  - Input validation tests
  - Error handling verification
  - Module integrity verification

- **Build & Deployment**
  - Makefile for build automation
  - Cross-platform support (Windows/Linux/macOS)
  - Installation scripts (Unix/Windows)
  - Single self-contained binary (~18 MB)

### Technical Details

#### Technology Stack
- **Language**: Go 1.21+
- **PDF Processing**: go-pdfium v1.17.2 (WebAssembly)
- **CLI Framework**: Cobra v1.10.1
- **WebAssembly Runtime**: Wazero v1.9.0
- **Image Processing**: Go standard library (image/png, image/jpeg)

#### Architecture
- Pure Go implementation (no CGO)
- WebAssembly-based PDF rendering
- Modular package structure
- Thread-safe PDF processing

#### Performance
- Build time: ~2-3 seconds
- Conversion time per page: <1 second
- Memory usage: ~50-100 MB during conversion
- Binary size: ~18 MB (all dependencies included)

### Security
- ✅ 0 Critical vulnerabilities
- ✅ 0 High vulnerabilities
- ✅ go-pdfium v1.17.2 (stable, no CVEs)
- ✅ Cobra v1.10.1 (latest, patched)
- ✅ All dependencies verified
- ✅ Input validation implemented
- ✅ Error handling complete
- ✅ No hardcoded secrets
- ✅ No unsafe imports

### Quality Metrics
- Security Tests: 19/19 ✅
- Code Coverage: Established baseline
- Dependency Verification: Passed ✅
- Module Integrity: Verified ✅
- Input Validation: Implemented ✅
- Error Handling: Complete ✅

### Dependencies
```
Direct:
  - github.com/klippa-app/go-pdfium v1.17.2 (Apache 2.0)
  - github.com/spf13/cobra v1.10.1 (Apache 2.0)

Transitive (14 minor updates available, non-critical)
```

### Breaking Changes
None (initial release)

### Known Issues
None reported

### Deprecations
None

---

## How to Upgrade

### From 0.x to 1.0.0
First release - no previous versions.

### Future Versions
To upgrade to a newer version:

```bash
# Download latest source
git clone https://github.com/tu-usuario/pdf2img
cd pdf2img

# Build locally
go build -o pdf2img ./cmd/pdf2img

# Or install globally
go install ./cmd/pdf2img
```

---

## Version Timeline

| Version | Release Date | Status | EOL Date |
| ------- | ------------ | ------ | -------- |
| 1.0.0 | 2025-11-07 | Current | 2027-11-07 |
| 0.x | Never | N/A | N/A |

---

## Contribution

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Contributors to 1.0.0
- Initial development and testing

---

## Security

For security vulnerabilities, please see [SECURITY.md](SECURITY.md).

Do not create public issues for security problems.

---

## References

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [Go Release Policy](https://golang.org/doc/devel/release)

---

**Last Updated**: 2025-11-07
**Current Version**: 1.0.0
**License**: Apache 2.0
