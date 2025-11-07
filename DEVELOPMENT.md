# Guía de Desarrollo

## Configuración del entorno de desarrollo

### Requisitos
- Go 1.21+
- Git
- Make (opcional, pero recomendado)

### Setup inicial

```bash
# Clonar el repositorio
git clone <repository-url>
cd pdf2img

# Descargar dependencias
go mod download
go mod verify

# Construir el binario
make build
# o
go build -o pdf2img ./cmd/pdf2img
```

## Estructura del código

### `pkg/converter/converter.go`
Módulo principal que encapsula la lógica de conversión PDF:
- `Converter`: Estructura principal que gestiona la instancia de PDFium
- `ConvertOptions`: Opciones para configurar la conversión
- `ConvertResult`: Resultado de la conversión
- Funciones helper para validación, rango de páginas, guardado de imágenes

### `cmd/pdf2img/main.go`
Aplicación CLI usando Cobra:
- Comando principal `convert` (por defecto)
- Subcomando `info` para obtener metadatos del PDF
- Manejo de flags y opciones de línea de comandos
- Formateo de salida para el usuario

### `mcp/server.go`
Implementación del servidor Model Context Protocol:
- `MCPServer`: Estructura que expone herramientas MCP
- `GetTools()`: Lista herramientas disponibles
- `ExecuteTool()`: Ejecuta una herramienta con entrada JSON
- Handlers para `pdf_to_images` y `pdf_info`

### `cmd/mcp-server/main.go`
Servidor HTTP/stdio para MCP:
- Modo HTTP: API REST en localhost
- Modo stdio: Para integración con Claude

## Compilación y pruebas

### Compilar
```bash
# Compilación optimizada (release)
make build

# Compilación con símbolos de debug
make build-dev

# Compilación específica del CLI
go build -o pdf2img ./cmd/pdf2img

# Compilación específica del MCP server
go build -o mcp-server ./cmd/mcp-server
```

### Pruebas
```bash
# Ejecutar todas las pruebas
make test
# o
go test -v ./...

# Pruebas de un paquete específico
go test -v ./pkg/converter
```

### Linting y formato
```bash
# Formatear código
make fmt

# Lint (requiere golangci-lint instalado)
make lint
```

## Trabajar con go-pdfium

### Inicialización
```go
conv, err := converter.New()
defer conv.Close()
```

### Obtener información del PDF
```go
info, err := conv.GetPDFInfo("documento.pdf")
// info["pages"], info["width"], info["height"], etc.
```

### Convertir PDF
```go
result, err := conv.Convert(&converter.ConvertOptions{
    InputPath: "documento.pdf",
    OutputDir: "./output",
    Format:    "png",
    DPI:       150,
    Prefix:    "page_",
})
```

## Extensiones futuras

### 1. Soporte para PDF con contraseña
```go
// En converter.go
doc, err := c.pdfium.OpenFile(pdfPath, &requests.OpenFile{
    Password: "tu-contraseña",
})
```

### 2. Procesamiento paralelo
```go
// Renderizar múltiples páginas concurrentemente
var wg sync.WaitGroup
for pageNum := startPage; pageNum <= endPage; pageNum++ {
    wg.Add(1)
    go func(p int) {
        defer wg.Done()
        // Renderizar página
    }(pageNum)
}
wg.Wait()
```

### 3. API REST completa
Extender `cmd/mcp-server/main.go` con más endpoints:
```go
http.HandleFunc("/api/v1/convert", handleConvert)
http.HandleFunc("/api/v1/info", handleInfo)
http.HandleFunc("/api/v1/status", handleStatus)
```

### 4. OCR de imágenes
Integrar con `github.com/otiai10/gosseract`:
```go
client, _ := gosseract.NewClient()
defer client.Close()
text, _ := client.Text()
```

### 5. Caché de renderizado
Almacenar imágenes generadas para evitar re-renderizar:
```go
type Cache struct {
    mu    sync.RWMutex
    items map[string][]byte
}
```

## Debugging

### Con logs detallados
```bash
pdf2img -i documento.pdf -o ./output -v
```

### En código
```go
import "log"

log.Printf("Debug: %+v", variable)
```

### Profiling de Go
```bash
# CPU profiling
go run -cpuprofile=cpu.prof ./cmd/pdf2img/main.go

# Memory profiling
go run -memprofile=mem.prof ./cmd/pdf2img/main.go

# Análisis
go tool pprof cpu.prof
```

## Dependencias

| Librería | Versión | Propósito | Licencia |
|----------|---------|----------|----------|
| go-pdfium | v1.12.0 | Renderizar PDFs | Apache 2.0 |
| Cobra | v1.7.0 | CLI framework | Apache 2.0 |

Para actualizar dependencias:
```bash
go get -u github.com/klippa-app/go-pdfium
go get -u github.com/spf13/cobra
go mod tidy
```

## Contribución

1. Fork el proyecto
2. Crea una rama (`git checkout -b feature/AmazingFeature`)
3. Haz commits claros (`git commit -m 'Add feature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Apache 2.0 - Ver LICENSE
