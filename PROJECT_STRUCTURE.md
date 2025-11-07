# Estructura del Proyecto

## Árbol de directorios

```
pdf2img/
├── cmd/                          # Aplicaciones ejecutables
│   ├── pdf2img/                  # CLI principal
│   │   └── main.go              # Aplicación de línea de comandos
│   └── mcp-server/               # Servidor MCP
│       └── main.go              # Servidor HTTP/stdio
├── mcp/                          # Implementación del protocolo MCP
│   ├── server.go                # Servidor MCP
│   └── example_server.go        # Ejemplo de uso
├── pkg/                          # Paquetes reutilizables
│   └── converter/               # Lógica de conversión
│       ├── converter.go         # Logica principal
│       └── converter_test.go    # Pruebas unitarias
├── go.mod                        # Definición de módulo Go
├── go.sum                        # Checksums de dependencias (generado)
├── Makefile                      # Automatización de build
├── install.sh                    # Script de instalación (Unix)
├── install.bat                   # Script de instalación (Windows)
├── .gitignore                    # Archivos a ignorar en Git
├── LICENSE                       # Licencia Apache 2.0
├── claude_config.json            # Configuración para Claude
├── README.md                     # Documentación principal
├── QUICKSTART.md                 # Guía de inicio rápido
├── EXAMPLES.md                   # Ejemplos de uso
├── DEVELOPMENT.md                # Guía para desarrolladores
└── PROJECT_STRUCTURE.md          # Este archivo
```

## Descripción de componentes

### `cmd/pdf2img/`
**Aplicación CLI principal**

- Interfaz de línea de comandos para convertir PDFs
- Usa Cobra para manejo de flags
- Soporta múltiples opciones: DPI, formato, rango de páginas, etc.
- Comando `info` para obtener metadatos del PDF

### `cmd/mcp-server/`
**Servidor MCP (Model Context Protocol)**

- Proporciona dos interfaces:
  - **HTTP**: API REST en puerto 8080
  - **stdio**: Para integración con Claude
- Expone herramientas MCP:
  - `pdf_to_images`: Convertir PDF a imágenes
  - `pdf_info`: Obtener información del PDF

### `mcp/`
**Lógica del servidor MCP**

- `server.go`: Implementación del servidor
  - `NewMCPServer()`: Crear instancia
  - `GetTools()`: Listar herramientas disponibles
  - `ExecuteTool()`: Ejecutar una herramienta
  - Handlers para cada herramienta

- `example_server.go`: Ejemplo de uso programático

### `pkg/converter/`
**Lógica compartida de conversión**

- `converter.go`: Módulo principal
  - `New()`: Inicializar PDFium
  - `Convert()`: Renderizar PDF a imágenes
  - `GetPDFInfo()`: Obtener metadatos
  - Funciones helper para validación y guardado

- `converter_test.go`: Tests unitarios
  - Validación de opciones
  - Determinación de rango de páginas
  - Formateo de tamaño de archivo

## Flujo de datos

### Conversión de PDF (CLI)

```
CLI (main.go)
  ├─→ Parsear argumentos (Cobra)
  ├─→ Validar opciones
  ├─→ Crear Converter
  ├─→ converter.Convert()
  │   ├─→ Validar opciones
  │   ├─→ PDFium.OpenFile()
  │   ├─→ Para cada página:
  │   │   ├─→ PDFium.RenderPageInDPI()
  │   │   └─→ saveImage() (PNG/JPG)
  │   └─→ Retornar ConvertResult
  └─→ Mostrar resultados
```

### Conversión vía MCP

```
Cliente (HTTP/stdio)
  ├─→ Enviar JSON {"tool": "pdf_to_images", "input": {...}}
  ├─→ MCPServer.ExecuteTool()
  │   ├─→ converter.Convert()
  │   └─→ Formatear resultado como JSON
  └─→ Retornar resultado JSON
```

## Dependencias

```
github.com/klippa-app/go-pdfium (v1.12.0)
  └─→ PDFium (C++)
      └─→ Renderización de PDF

github.com/spf13/cobra (v1.7.0)
  └─→ Framework CLI con flags
```

## Configuración

### Variables de entorno
Actualmente no requiere variables de entorno.

### Archivos de configuración
- `claude_config.json`: Configuración para integración con Claude

## Convenciones

### Nomenclatura
- Paquetes: `lowercase` sin underscore
- Funciones/métodos: `PascalCase` (públicas), `camelCase` (privadas)
- Constantes: `SCREAMING_SNAKE_CASE`
- Variables: `camelCase`

### Estructura de tests
- Nombre archivo: `{nombre}_test.go`
- Nombre función: `Test{FuncionAProbar}`
- Tabla de casos: `tests := []struct {...}`

### Error handling
- Propagar errores con contexto: `fmt.Errorf("acción fallida: %w", err)`
- No usar `panic()` excepto en inicialización
- Siempre defer para limpieza: `defer converter.Close()`

## Compilación

### Compilar CLI
```bash
go build -o pdf2img ./cmd/pdf2img
```

### Compilar MCP server
```bash
go build -o mcp-server ./cmd/mcp-server
```

### Compilación optimizada
```bash
go build -ldflags "-s -w" -o pdf2img ./cmd/pdf2img
```

## Testing

### Ejecutar todos los tests
```bash
go test -v ./...
```

### Test específico
```bash
go test -v ./pkg/converter -run TestValidateOptions
```

### Coverage
```bash
go test -cover ./...
```

## Performance

### Optimizaciones realizadas
1. **Lazy initialization**: PDFium se inicializa bajo demanda
2. **Resource cleanup**: `defer` para cerrar archivos y liberar memoria
3. **Error handling eficiente**: Sin recrear objetos innecesarios

### Posibles mejoras
1. Procesamiento paralelo de páginas
2. Caché de imágenes renderizadas
3. Pool de workers para procesamiento batch
4. Compresión de imágenes en salida

## Seguridad

### Validaciones
1. Verificación de archivo de entrada
2. Validación de formato de salida
3. Validación de rango de páginas
4. Creación segura de directorios

### Consideraciones
1. No aceptar rutas absolutas sin validación
2. Sanitizar nombres de salida
3. Limitar tamaño máximo de DPI
4. Proteger contra PDF malintencionados

## Internacionalización (i18n)
Actualmente sin soporte i18n. Todos los mensajes en inglés/español.

Para futura expansión:
- Usar `golang.org/x/text/language`
- Separar strings en archivos de mensajes
- Traducir automáticamente

## Versionado
Semántico (MAJOR.MINOR.PATCH):
- MAJOR: Cambios incompatibles
- MINOR: Nuevas características
- PATCH: Correcciones de bugs

Actualmente: v1.0.0

## Roadmap

### v1.1.0
- [ ] Soporte para PDF con contraseña
- [ ] Procesamiento paralelo de páginas
- [ ] Mejor reporte de errores

### v1.2.0
- [ ] OCR de imágenes
- [ ] Caché de renderizado
- [ ] API HTTP adicional

### v2.0.0
- [ ] Rediseño de arquitectura
- [ ] Soporte para más formatos (WebP, TIFF)
- [ ] Base de datos para metadatos
