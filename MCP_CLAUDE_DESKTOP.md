# 🚀 Guía de Configuración: PDF2IMG con Claude Desktop

Esta guía te explica paso a paso cómo configurar y usar el MCP Server de pdf2img con Claude Desktop.

---

## 📋 Requisitos

- **Claude Desktop** instalado (descárgalo desde https://claude.ai/download)
- **mcp-server.exe** compilado en tu proyecto
- Archivo PDF para convertir (para pruebas)

---

## 🔧 Configuración (5 pasos)

### Paso 1: Compilar el MCP Server

Primero, asegúrate de que el MCP server esté compilado:

```bash
# Desde la raíz del proyecto
go build -o mcp-server.exe ./cmd/mcp-server
```

Verificar que se creó el archivo:
```bash
ls -la mcp-server.exe
```

Debería tener ~18 MB.

### Paso 2: Ubicar el archivo de configuración

Claude Desktop guarda su configuración en:

**Windows**:
```
%APPDATA%\Claude\claude_desktop_config.json
```

**Acceso rápido**:
- Abre Claude Desktop
- Click en el menú (⚙️) en la esquina inferior izquierda
- Selecciona "Developer" o "Settings"
- Click en "Edit Config" para abrir el archivo

**macOS**:
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Linux**:
```
~/.config/Claude/claude_desktop_config.json
```

### Paso 3: Editar la configuración

Abre el archivo `claude_desktop_config.json` con tu editor de texto favorito.

**Estructura actual** (si no existe, créalo):

```json
{
  "mcpServers": {
    "other-server": {
      "command": "...",
      "args": [...]
    }
  }
}
```

**Agregar pdf2img**:

Copia y pega la siguiente configuración en la sección `mcpServers`:

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

**⚠️ IMPORTANTE**: Cambia `C:\\MCPs\\clone_PROYECTOS\\mcp-go-pdf-to-img-2` por la ruta **REAL** donde tienes tu proyecto.

**Ejemplos de rutas válidas**:

Windows:
```json
"command": "C:\\Users\\tu-usuario\\Projects\\pdf2img\\mcp-server.exe"
"command": "D:\\codigo\\mcp-go-pdf-to-img-2\\mcp-server.exe"
```

macOS/Linux:
```json
"command": "/Users/tu-usuario/projects/pdf2img/mcp-server.exe"
"command": "/home/usuario/projects/pdf2img/mcp-server"
```

**Archivo completo de ejemplo**:

```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\MCPs\\clone_PROYECTOS\\mcp-go-pdf-to-img-2\\mcp-server.exe",
      "args": ["--stdio"],
      "env": {}
    },
    "other-server": {
      "command": "other-command",
      "args": []
    }
  }
}
```

### Paso 4: Guardar y reiniciar Claude Desktop

1. Guarda el archivo `claude_desktop_config.json`
2. **Cierra completamente** Claude Desktop
3. **Abre de nuevo** Claude Desktop

### Paso 5: Verificar que está conectado

Abre una conversación en Claude y busca en la esquina inferior derecha. Deberías ver:

```
Connected to: pdf2img
✓ pdf_to_images
✓ pdf_info
```

Si ves esto, ¡está funcionando! 🎉

Si no lo ves, revisa:
- ¿La ruta en la configuración es correcta?
- ¿El archivo `mcp-server.exe` existe en esa ruta?
- ¿Guardaste y reiniciaste Claude Desktop?

---

## 📚 Usando pdf2img en Claude Desktop

Una vez configurado, puedes usar pdf2img de forma natural:

### Ejemplo 1: Convertir un PDF

**Tú**:
> Convierte el archivo documento.pdf a imágenes PNG con 300 DPI. Guarda las imágenes en la carpeta output.

**Claude**:
Automáticamente ejecutará la herramienta `pdf_to_images` con los parámetros correctos.

### Ejemplo 2: Obtener información de un PDF

**Tú**:
> ¿Cuántas páginas tiene el archivo reporte.pdf? ¿Cuáles son sus dimensiones?

**Claude**:
Automáticamente ejecutará la herramienta `pdf_info` para obtener la información.

### Ejemplo 3: Convertir solo algunas páginas

**Tú**:
> Quiero convertir solo las primeras 10 páginas del documento.pdf a JPG con 150 DPI

**Claude**:
Ejecutará `pdf_to_images` con `start_page: 1` y `end_page: 10`.

---

## 🛠️ Herramientas Disponibles

### Herramienta 1: `pdf_to_images`

**Qué hace**: Convierte páginas de un PDF a imágenes PNG o JPG.

**Parámetros**:
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `pdf_path` | string | ✅ | Ruta al archivo PDF | `documento.pdf` |
| `output_dir` | string | ✅ | Directorio de salida | `./output` |
| `format` | string | ❌ | Formato: `png` o `jpg` | `png` (default) |
| `dpi` | number | ❌ | Resolución en DPI | `150` (default) |
| `start_page` | integer | ❌ | Primera página (1-indexed) | `1` |
| `end_page` | integer | ❌ | Última página | `50` |
| `prefix` | string | ❌ | Prefijo de archivos | `page_` (default) |

**Ejemplo de respuesta**:
```json
{
  "total_pages": 25,
  "successful": 25,
  "failed": 0,
  "files": [
    "./output/page_0001.png",
    "./output/page_0002.png",
    "./output/page_0003.png",
    ...
  ]
}
```

### Herramienta 2: `pdf_info`

**Qué hace**: Obtiene información sobre un PDF (número de páginas, tamaño, dimensiones).

**Parámetros**:
| Parámetro | Tipo | Requerido | Descripción |
|-----------|------|-----------|-------------|
| `pdf_path` | string | ✅ | Ruta al archivo PDF |

**Ejemplo de respuesta**:
```json
{
  "file": "documento.pdf",
  "pages": 25,
  "file_size": "2.50 MB",
  "width": 612.00,
  "height": 792.00
}
```

---

## 💡 Casos de Uso Comunes

### 1. Generar miniaturas de un PDF

**Tú**:
> Genera miniaturas (primera página) de documento.pdf en formato JPG con 96 DPI. Guárdalo en output/thumbnail.jpg

**Claude** ejecuta:
```
pdf_to_images {
  pdf_path: "documento.pdf",
  output_dir: "./output",
  format: "jpg",
  dpi: 96,
  start_page: 1,
  end_page: 1,
  prefix: "thumbnail_"
}
```

### 2. Convertir documento completo en alta resolución

**Tú**:
> Necesito el documento.pdf convertido a PNG en alta resolución (300 DPI) para imprimir. Guarda en images/.

**Claude** ejecuta:
```
pdf_to_images {
  pdf_path: "documento.pdf",
  output_dir: "./images",
  format: "png",
  dpi: 300
}
```

### 3. Procesar solo parte de un PDF

**Tú**:
> Del documento.pdf, quiero solo las páginas 5 a 15 convertidas a JPG con nombre img_XXX.jpg

**Claude** ejecuta:
```
pdf_to_images {
  pdf_path: "documento.pdf",
  output_dir: "./output",
  format: "jpg",
  start_page: 5,
  end_page: 15,
  prefix: "img_"
}
```

### 4. Analizar múltiples PDFs

**Tú**:
> Tengo 3 PDFs: reporte1.pdf, reporte2.pdf y reporte3.pdf. ¿Cuántas páginas tiene cada uno?

**Claude** ejecuta `pdf_info` tres veces:
```
pdf_info { pdf_path: "reporte1.pdf" }
pdf_info { pdf_path: "reporte2.pdf" }
pdf_info { pdf_path: "reporte3.pdf" }
```

---

## 🔍 Solución de Problemas

### El MCP no aparece conectado en Claude

**Problema**: No ves "pdf2img" en la esquina inferior de Claude.

**Soluciones**:
1. Verifica la ruta del archivo en `claude_desktop_config.json`
2. Asegúrate de que `mcp-server.exe` existe
3. Abre el archivo con `cmd` para verificar:
   ```bash
   C:\MCPs\clone_PROYECTOS\mcp-go-pdf-to-img-2\mcp-server.exe --help
   ```
4. Si sale un error, compila de nuevo:
   ```bash
   go build -o mcp-server.exe ./cmd/mcp-server
   ```

### La conversión falla

**Problema**: Claude intenta convertir pero retorna un error.

**Causas comunes**:
- ❌ Ruta de PDF incorrecta → Usa rutas absolutas
- ❌ Directorio de salida no existe → Claude lo crea automáticamente
- ❌ Archivo PDF corrupto → Verifica el PDF con `pdf2img info`
- ❌ Permisos de escritura → Asegúrate de tener permisos en el directorio

### Las imágenes se ven borrosas

**Solución**: Aumenta el DPI

```
De: dpi: 72
A: dpi: 300
```

### El proceso es lento

**Solución**: Disminuye el DPI o usa rangos de páginas

```
De: dpi: 300, todas las páginas
A: dpi: 150, start_page: 1, end_page: 50
```

### Error: "File not found"

**Solución**: Usa rutas absolutas

```
❌ Incorrecto: documento.pdf
✅ Correcto: C:\Users\tu-usuario\Documents\documento.pdf
✅ Correcto: /home/usuario/documents/documento.pdf
```

---

## 📝 Configuración Avanzada

### Usar en otra máquina

Si quieres usar pdf2img en otra máquina:

1. Copia todo el directorio del proyecto
2. Compila el MCP server en esa máquina:
   ```bash
   go build -o mcp-server.exe ./cmd/mcp-server
   ```
3. Actualiza la ruta en `claude_desktop_config.json`

### Usar con variables de entorno

Si necesitas configurar variables de entorno (futuro):

```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\..\\mcp-server.exe",
      "args": ["--stdio"],
      "env": {
        "LOG_LEVEL": "debug",
        "MAX_DPI": "600"
      }
    }
  }
}
```

### Usar puerto HTTP en lugar de stdio

Si prefieres HTTP (no recomendado para Claude Desktop):

1. Modifica `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\..\\mcp-server.exe",
      "args": ["--port", "8080"]
    }
  }
}
```

2. Luego conecta manualmente a `http://localhost:8080`

---

## 📚 Documentación Relacionada

- [EXAMPLES.md](EXAMPLES.md) - Ejemplos completos de uso
- [README.md](README.md) - Información general del proyecto
- [DEVELOPMENT.md](DEVELOPMENT.md) - Cómo desarrollar pdf2img

---

## ✅ Checklist de Configuración

- [ ] Claude Desktop instalado
- [ ] mcp-server.exe compilado
- [ ] Ruta correcta en `claude_desktop_config.json`
- [ ] Archivo guardado
- [ ] Claude Desktop reiniciado
- [ ] "pdf2img" aparece conectado en Claude
- [ ] Probé un ejemplo simple (pdf_info)
- [ ] Probé una conversión (pdf_to_images)

---

## 🆘 ¿Necesitas ayuda?

1. Revisa [EXAMPLES.md](EXAMPLES.md) para más ejemplos
2. Revisa [SECURITY.md](SECURITY.md) para problemas de seguridad
3. Crea un issue en GitHub: https://github.com/tu-usuario/pdf2img/issues

---

**Última actualización**: 2025-11-07
**Versión**: 1.0.0
