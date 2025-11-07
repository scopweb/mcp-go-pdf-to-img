# PDF2IMG - Convert PDF to Images

Una herramienta completa en Go para renderizar páginas PDF como imágenes PNG o JPG. Incluye tanto una aplicación CLI como un servidor MCP (Model Context Protocol).

## Características

- ✨ **Renderizado de alta calidad** usando go-pdfium v1.17.2 (basado en PDFium)
- 🖼️ **Múltiples formatos**: PNG y JPG
- 🎯 **Control granular**: DPI personalizable, rango de páginas configurable
- 💻 **CLI completa**: Interfaz de línea de comandos con múltiples opciones
- 🔌 **Servidor MCP**: Integración con Model Context Protocol
- 📊 **Información de PDF**: Comando para obtener metadatos de archivos PDF
- 🚀 **Pure Go**: Sin CGO, WebAssembly embebido, single binary (~18 MB)
- 🔒 **Multiplataforma**: Windows, Linux, macOS

## Instalación

### Requisitos previos

- Go 1.21 o superior
- **Nota**: No requiere PDFium externo, todo está embebido en el binario

### Desde el código fuente

```bash
git clone <repository-url>
cd pdf2img
go mod download
go build -o pdf2img ./cmd/pdf2img
```

### Instalación global

```bash
go install github.com/tu-usuario/pdf2img/cmd/pdf2img@latest
```

## Uso

### CLI - Conversión básica

```bash
# Convertir todas las páginas a PNG
pdf2img -i documento.pdf -o ./output

# Convertir a JPG con DPI personalizado
pdf2img -i documento.pdf -o ./output -f jpg -d 300

# Convertir solo páginas 1-10
pdf2img -i documento.pdf -o ./output --start 1 --end 10

# Con prefix personalizado
pdf2img -i documento.pdf -o ./output --prefix img_
```

### CLI - Obtener información del PDF

```bash
pdf2img info documento.pdf
```

Salida:
```
PDF Information
===============
File: documento.pdf
Pages: 25
Size: 2.50 MB
Width: 612.00 pt
Height: 792.00 pt
```

### Opciones de la CLI

| Opción | Corto | Descripción | Default |
|--------|-------|-------------|---------|
| `--input` | `-i` | Ruta del PDF (requerido) | - |
| `--output` | `-o` | Directorio de salida | `.` |
| `--format` | `-f` | Formato: png o jpg | `png` |
| `--dpi` | `-d` | DPI para renderizado | `150` |
| `--start` | - | Página inicial (1-indexed) | `0` (primera) |
| `--end` | - | Página final (1-indexed) | `0` (última) |
| `--prefix` | - | Prefijo de archivos de salida | `page_` |
| `--verbose` | `-v` | Salida detallada | `false` |

### Servidor MCP

El MCP Server de pdf2img permite integración con Claude Desktop y otros sistemas que soporten el Model Context Protocol.

#### Configuración con Claude Desktop (Recomendado)

**⚡ Guía rápida**: Lee [MCP_CLAUDE_DESKTOP.md](MCP_CLAUDE_DESKTOP.md) para instrucciones paso a paso.

**Resumen**:
1. Compila: `go build -o mcp-server.exe ./cmd/mcp-server`
2. Edita `claude_desktop_config.json` (en `%APPDATA%\Claude\`)
3. Agrega:
```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\ruta\\a\\mcp-server.exe",
      "args": ["--stdio"]
    }
  }
}
```
4. Reinicia Claude Desktop
5. ¡Listo! Ahora puedes usar pdf2img desde Claude

**¿Problemas?** Lee [MCP_TROUBLESHOOTING.md](MCP_TROUBLESHOOTING.md) para solucionar problemas comunes.

#### Uso programático

Para usar el MCP Server en código Go:

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

	// Usar server.ExecuteTool() para ejecutar herramientas
	tools := server.GetTools()
	for _, tool := range tools {
		log.Printf("Disponible: %s", tool.Name)
	}
}
```

#### Herramientas disponibles en MCP

##### `pdf_to_images`

Convierte páginas PDF a imágenes PNG o JPG.

```json
{
  "pdf_path": "documento.pdf",
  "output_dir": "./output",
  "format": "png",
  "dpi": 150,
  "start_page": 0,
  "end_page": 0,
  "prefix": "page_"
}
```

**Nota**: `start_page` y `end_page` con valor 0 significan "todas las páginas".

##### `pdf_info`

Obtiene información del PDF (páginas, tamaño, dimensiones).

```json
{
  "pdf_path": "documento.pdf"
}
```

#### Modos de operación

El MCP Server soporta dos modos:

- **stdio** (por defecto para Claude Desktop): `mcp-server --stdio`
- **HTTP**: `mcp-server --port 8080` (para pruebas)

Más ejemplos en [EXAMPLES.md](EXAMPLES.md#mcp-server---ejemplos-de-integración).

## Estructura del proyecto

```
pdf2img/
├── cmd/
│   └── pdf2img/
│       └── main.go          # Aplicación CLI
├── mcp/
│   └── server.go            # Servidor MCP
├── pkg/
│   └── converter/
│       └── converter.go      # Lógica compartida de conversión
├── go.mod                   # Dependencias
└── README.md               # Este archivo
```

## Dependencias

- **go-pdfium v1.17.2** (Apache 2.0): Renderización de PDF usando PDFium WebAssembly
- **cobra v1.7.0** (Apache 2.0): Framework CLI
- **wazero v1.9.0** (Apache 2.0): Runtime WebAssembly para Go (incluido en go-pdfium)

## Ejemplos

### Ejemplo 1: Convertir documento completo a alta resolución

```bash
pdf2img -i libro.pdf -o ./images -f jpg -d 300
```

Crea archivos como: `page_0001.jpg`, `page_0002.jpg`, etc.

### Ejemplo 2: Extraer solo las primeras 5 páginas

```bash
pdf2img -i reporte.pdf -o ./thumbs --start 1 --end 5 --prefix thumb_
```

Crea: `thumb_0001.png`, `thumb_0002.png`, etc.

### Ejemplo 3: Información del PDF

```bash
pdf2img info formulario.pdf
```

### Ejemplo 4: Desde código Go

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
		InputPath: "documento.pdf",
		OutputDir: "./output",
		Format:    "png",
		DPI:       150,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Exitosas: %d, Fallidas: %d", result.Successful, result.Failed)
}
```

## Tecnología

### WebAssembly Implementation

Este proyecto utiliza **go-pdfium v1.17.2 con WebAssembly**, lo que proporciona:

- **Pure Go**: Sin dependencias de CGO, sin necesidad de compiladores C en el host
- **Single Binary**: Todo embebido en un ejecutable (~18 MB)
- **Multiplataforma**: El mismo binario funciona en Windows, Linux, macOS
- **Performance**: ~2x más rápido que la versión multi-threaded CGO
- **Seguridad**: Ejecución aislada en sandbox WebAssembly

### Por qué WebAssembly

Se eligió WebAssembly sobre CGO porque:
- ❌ CGO requiere compilador C en el host (problemático en diferentes sistemas)
- ❌ go-pdfium v1.12.0 CGO tenía bugs graves de incompatibilidad
- ✅ WebAssembly es más robusto y portable
- ✅ No requiere dependencias externas después de compilar
- ✅ Mejor aislamiento de recursos

## 📚 Documentación

### Documentación Pública
- **[QUICKSTART.md](QUICKSTART.md)** - Empieza en 5 minutos
- **[EXAMPLES.md](EXAMPLES.md)** - Casos prácticos de uso
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Guía de desarrollo
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Estructura del código
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Cómo contribuir

### Documentación Interna
- **[docs/README.md](docs/README.md)** - Índice de documentación interna
- **[docs/STATUS.md](docs/STATUS.md)** - Estado del proyecto
- **[docs/SECURITY_TEST_RESULTS.md](docs/SECURITY_TEST_RESULTS.md)** - Tests de seguridad
- **[docs/IMPLEMENTATION_NOTES.md](docs/IMPLEMENTATION_NOTES.md)** - Notas técnicas

### Más Información
- **[CHANGELOG.md](CHANGELOG.md)** - Historial de cambios
- **[SECURITY.md](SECURITY.md)** - Política de seguridad
- **[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)** - Código de conducta

---

## Solución de problemas

### Error: "PDF file not found"
Verifica que la ruta al PDF sea correcta:
```bash
ls -la documento.pdf
```

### Error: "failed to initialize PDFium"
Asegúrate de que Go puede descargar las dependencias:
```bash
go mod download
go mod verify
```

### Las imágenes salen en baja resolución
Aumenta el valor de DPI:
```bash
pdf2img -i documento.pdf -o ./output -d 300  # Más alto = más calidad
```

### El proceso es lento
Reduce DPI para procesamiento más rápido:
```bash
pdf2img -i documento.pdf -o ./output -d 96
```

## Licencia

Este proyecto usa go-pdfium bajo licencia Apache 2.0.

## Contribución

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Roadmap

- [ ] Soporte para marcas de agua
- [ ] OCR de imágenes renderizadas
- [ ] Procesamiento paralelo de páginas
- [ ] Caché de imágenes renderizadas
- [ ] API HTTP REST
- [ ] Soporte para PDF con contraseña

## Referencias

- [go-pdfium](https://github.com/klippa-app/go-pdfium)
- [PDFium](https://pdfium.googlesource.com/)
- [Cobra CLI](https://cobra.dev/)
