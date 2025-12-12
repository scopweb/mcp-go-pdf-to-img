# Ejemplos de Uso

## CLI - Ejemplos prácticos

### 1. Conversión básica - Todas las páginas a PNG

```bash
pdf2img -i documento.pdf -o ./images
```

Resultado:
```
✓ Conversion Complete
Total pages: 42
Successful: 42
Failed: 0
```

Archivos generados:
- `images/page_0001.png`
- `images/page_0002.png`
- ...
- `images/page_0042.png`

### 2. Convertir a JPG con alta resolución

```bash
pdf2img -i documento.pdf -o ./output -f jpg -d 300
```

Para imprimir, usa DPI 300. Para pantalla, 72-150 es suficiente.

### 3. Extraer rango específico de páginas

```bash
pdf2img -i libro.pdf -o ./thumbs --start 1 --end 10
```

Solo convierte las primeras 10 páginas.

### 4. Usar prefijo personalizado

```bash
pdf2img -i reporte.pdf -o ./output --prefix img_
```

Genera: `img_0001.png`, `img_0002.png`, etc.

### 5. Conversión con salida detallada

```bash
pdf2img -i documento.pdf -o ./output -v
```

Muestra cada archivo generado y errores detallados.

### 6. Obtener información del PDF

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

## MCP Server - Ejemplos de integración

### 1. Configuración con Claude Desktop (Recomendado)

#### Requisitos
- Claude Desktop instalado
- `mcp-server.exe` compilado en el proyecto

#### Paso 1: Ubicar el archivo de configuración

El archivo `claude_desktop_config.json` está ubicado en:

**Windows**:
```
%APPDATA%\Claude\claude_desktop_config.json
```

O simplemente abre Claude Desktop → Configuración (⚙️) → Developer settings → Edit config

**macOS**:
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Linux**:
```
~/.config/Claude/claude_desktop_config.json
```

#### Paso 2: Agregar la configuración de pdf2img

Edita el archivo `claude_desktop_config.json` y agrega la siguiente sección en `mcpServers`:

```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\ruta\\al\\proyecto\\mcp-server.exe",
      "args": ["--stdio"],
      "env": {}
    }
  }
}
```

**Importante**: Usa la ruta completa donde compilaste `mcp-server.exe`. Por ejemplo:
```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\MCPs\\clone_PROYECTOS\\mcp-go-pdf-to-img-2\\mcp-server.exe",
      "args": ["--stdio"],
      "env": {}
    }
  }
}
```

#### Paso 3: Reiniciar Claude Desktop

Cierra y abre Claude Desktop nuevamente. Deberías ver el MCP server "pdf2img" conectado.

#### Paso 4: Usar con Claude

Ahora puedes usar pdf2img directamente en Claude:

**Ejemplo 1**: "Convierte el documento.pdf a imágenes PNG con 300 DPI"

Claude automáticamente utilizará la herramienta `pdf_to_images` con los parámetros correctos.

**Ejemplo 2**: "¿Cuántas páginas tiene el archivo ejemplo.pdf?"

Claude utilizará `pdf_info` para obtener la información.

### 2. Herramientas Disponibles

#### Herramienta: `pdf_to_images`

Convierte páginas de un PDF a imágenes PNG o JPG.

**Parámetros**:
```json
{
  "pdf_path": "documento.pdf",           // Ruta del PDF (requerido)
  "output_dir": "./output",              // Directorio de salida (requerido)
  "format": "png",                       // png o jpg (opcional, default: png)
  "dpi": 150,                            // Resolución (opcional, default: 150)
  "start_page": 1,                       // Primera página (opcional, default: todas)
  "end_page": 10,                        // Última página (opcional, default: todas)
  "prefix": "page_"                      // Prefijo de archivos (opcional, default: page_)
}
```

**Respuesta**:
```json
{
  "total_pages": 42,
  "successful": 42,
  "failed": 0,
  "files": [
    "./output/page_0001.png",
    "./output/page_0002.png",
    ...
  ]
}
```

#### Herramienta: `pdf_info`

Obtiene información sobre un PDF.

**Parámetros**:
```json
{
  "pdf_path": "documento.pdf"  // Ruta del PDF (requerido)
}
```

**Respuesta**:
```json
{
  "file": "documento.pdf",
  "pages": 25,
  "file_size": "2.50 MB",
  "width": 612.00,
  "height": 792.00
}
```

#### Herramienta: `pdf_compress`

Comprime un PDF para reducir su tamaño optimizando imágenes y removiendo elementos innecesarios.

**Parámetros**:
```json
{
  "pdf_path": "documento.pdf",        // Ruta del PDF a comprimir (requerido)
  "output_path": "documento_comprimido.pdf"  // Ruta del PDF comprimido (requerido)
}
```

**Respuesta**:
```
PDF compressed successfully. Output saved to: documento_comprimido.pdf
```

### 3. Modo HTTP (Para pruebas)

Si prefieres usar HTTP en lugar de Claude Desktop:

```bash
mcp-server --port 8080
```

#### Obtener herramientas disponibles

```bash
curl http://localhost:8080/tools
```

Respuesta:
```json
[
  {
    "name": "pdf_to_images",
    "description": "Convert PDF pages to PNG or JPG images",
    "inputSchema": { ... }
  },
  {
    "name": "pdf_info",
    "description": "Get information about a PDF file",
    "inputSchema": { ... }
  }
]
```

#### Convertir PDF a imágenes

```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "pdf_to_images",
    "input": {
      "pdf_path": "documento.pdf",
      "output_dir": "./output",
      "format": "png",
      "dpi": 150
    }
  }'
```

Respuesta:
```json
{
  "type": "text",
  "content": "{\n  \"total_pages\": 42,\n  \"successful\": 42,\n  \"failed\": 0,\n  \"files\": [\n    \"./output/page_0001.png\",\n    \"./output/page_0002.png\",\n    ...\n  ]\n}\n"
}
```

#### Obtener información del PDF

```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "pdf_info",
    "input": {
      "pdf_path": "documento.pdf"
    }
  }'
```

### 4. Modo stdio (Para integración en scripts)

```bash
mcp-server --stdio
```

El servidor lee comandos JSON de stdin y escribe resultados en stdout:

```bash
echo '{"tool": "pdf_info", "input": {"pdf_path": "documento.pdf"}}' | mcp-server --stdio
```

## Desde código Go

### Básico

```go
package main

import (
	"log"
	"github.com/tu-usuario/pdf2img/pkg/converter"
)

func main() {
	// Crear convertidor
	conv, err := converter.New()
	if err != nil {
		log.Fatal(err)
	}
	defer conv.Close()

	// Convertir
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

### Con MCP Server

```go
package main

import (
	"encoding/json"
	"log"
	"github.com/tu-usuario/pdf2img/mcp"
)

func main() {
	server, err := mcp.NewMCPServer()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	// Convertir PDF usando MCP
	input := json.RawMessage(`{
		"pdf_path": "documento.pdf",
		"output_dir": "./output",
		"format": "png",
		"dpi": 150
	}`)

	result, err := server.ExecuteTool("pdf_to_images", input)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result.Content)
}
```

## Casos de uso reales

### 1. Generar miniaturas de PDFs en lote

```bash
#!/bin/bash
for pdf in *.pdf; do
    output_dir="thumbnails/${pdf%.pdf}"
    mkdir -p "$output_dir"
    pdf2img -i "$pdf" -o "$output_dir" -f jpg -d 96 --start 1 --end 1
    echo "✓ Miniatura de $pdf generada"
done
```

### 2. Convertir PDF a galería de imágenes web

```bash
pdf2img -i articulo.pdf -o ./web/images -f jpg -d 150
# Luego crear HTML que liste las imágenes
for img in web/images/*.jpg; do
    echo "<img src='$(basename $img)' />"
done > web/gallery.html
```

### 3. Integración con servicio web

```go
// En tu aplicación web
func handlePDFUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Upload failed", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Guardar PDF temporalmente
	tmpFile, _ := ioutil.TempFile("", "*.pdf")
	io.Copy(tmpFile, file)
	tmpFile.Close()

	// Convertir
	conv, _ := converter.New()
	defer conv.Close()
	result, _ := conv.Convert(&converter.ConvertOptions{
		InputPath: tmpFile.Name(),
		OutputDir: "./static/uploads",
		Format:    "jpg",
		DPI:       150,
	})

	// Retornar URLs de imágenes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.OutputFiles)
}
```

### 4. Procesamiento en background con workers

```go
type PDFJob struct {
	Path   string
	Output string
	DPI    float64
}

func worker(jobs <-chan PDFJob, results chan<- PDFJob) {
	conv, _ := converter.New()
	defer conv.Close()

	for job := range jobs {
		conv.Convert(&converter.ConvertOptions{
			InputPath: job.Path,
			OutputDir: job.Output,
			DPI:       job.DPI,
		})
		results <- job
	}
}

// Crear pool de workers
jobs := make(chan PDFJob, 100)
results := make(chan PDFJob)

for i := 0; i < 4; i++ {
	go worker(jobs, results)
}

// Encolar trabajos
for _, pdf := range pdfFiles {
	jobs <- PDFJob{Path: pdf, Output: "./output", DPI: 150}
}
```

## Solución de problemas comunes

### Las imágenes se ven borrosas
Aumenta DPI:
```bash
pdf2img -i documento.pdf -o ./output -d 300
```

### El proceso es muy lento
Reduce DPI:
```bash
pdf2img -i documento.pdf -o ./output -d 96
```

### Falta memoria (Out of Memory)
Procesa por rangos:
```bash
# Páginas 1-100
pdf2img -i documento.pdf -o ./out1 --start 1 --end 100

# Páginas 101-200
pdf2img -i documento.pdf -o ./out2 --start 101 --end 200
```

### Error: "PDF file not found"
Verifica la ruta:
```bash
ls -la documento.pdf
# Use ruta absoluta
pdf2img -i /path/completa/documento.pdf -o ./output
```

### MCP server no responde
Verifica que el puerto está disponible:
```bash
netstat -tuln | grep 8080
# O usa otro puerto
mcp-server -port 9090
```
