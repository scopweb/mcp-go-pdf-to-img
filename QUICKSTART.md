# Quick Start Guide

## 5 minutos para tener pdf2img funcionando

### Opción 1: Linux / macOS

```bash
# 1. Clona o descarga el proyecto
cd pdf2img

# 2. Ejecuta el script de instalación
chmod +x install.sh
./install.sh

# 3. ¡Listo! Usa:
./pdf2img -i tu-documento.pdf -o ./output
```

### Opción 2: Windows

```bash
# 1. Abre PowerShell como administrador
# 2. Navega a la carpeta del proyecto
cd pdf2img

# 3. Ejecuta el script de instalación
install.bat

# 4. ¡Listo! Usa:
pdf2img.exe -i tu-documento.pdf -o .\output
```

### Opción 3: Go directo (cualquier SO)

```bash
# Requiere Go 1.21+
go mod download
go build -o pdf2img ./cmd/pdf2img

# Usar:
./pdf2img -i documento.pdf -o ./output
```

## Primeros pasos

### Convertir un PDF completo

```bash
pdf2img -i documento.pdf -o ./salida
```

**Resultado:**
- ✓ Crea directorio `./salida`
- ✓ Genera `page_0001.png`, `page_0002.png`, etc.

### Obtener información del PDF

```bash
pdf2img info documento.pdf
```

**Resultado:**
```
PDF Information
===============
File: documento.pdf
Pages: 42
Size: 2.50 MB
```

### Convertir solo las primeras 5 páginas

```bash
pdf2img -i documento.pdf -o ./salida --start 1 --end 5
```

### Usar JPG en lugar de PNG

```bash
pdf2img -i documento.pdf -o ./salida -f jpg
```

### Aumentar calidad (más DPI)

```bash
pdf2img -i documento.pdf -o ./salida -d 300
```

| DPI | Uso |
|-----|-----|
| 72-96 | Web, pantalla |
| 150 | Balance calidad/velocidad (default) |
| 300 | Imprenta, documentos |
| 600+ | Archivos muy detallados |

## Casos de uso comunes

### Miniatura de primera página
```bash
pdf2img -i documento.pdf -o ./thumbs --start 1 --end 1 -d 96 --prefix thumb_
```

### Crear galería web
```bash
pdf2img -i documento.pdf -o ./web/images -f jpg -d 150
```

### Extraer rango específico
```bash
pdf2img -i libro.pdf -o ./capitulo5 --start 100 --end 150
```

## Usar como servidor MCP

### Opción HTTP (puerto 8080)

```bash
mcp-server
# Servidor en: http://localhost:8080

# En otra terminal:
curl http://localhost:8080/tools
```

### Opción stdio (integración con Claude)

```bash
mcp-server --stdio
```

## Documentación completa

- [README.md](README.md) - Documentación completa
- [EXAMPLES.md](EXAMPLES.md) - Ejemplos detallados
- [DEVELOPMENT.md](DEVELOPMENT.md) - Guía para desarrolladores

## Problemas comunes

| Problema | Solución |
|----------|----------|
| `command not found: pdf2img` | Asegúrate de estar en el directorio del proyecto o instalar globalmente con `--global` |
| PDFium no se descarga | Ejecuta `go mod download` |
| Las imágenes se ven borrosas | Aumenta DPI: `-d 300` |
| El proceso es lento | Reduce DPI: `-d 96` o especifica un rango: `--start 1 --end 10` |
| `file not found` | Verifica que el PDF existe: `ls documento.pdf` |

## Próximos pasos

1. ✅ **Primeros pasos**: Convertir tu primer PDF
2. 📖 **Aprender más**: Lee [EXAMPLES.md](EXAMPLES.md)
3. 🔧 **Personalizar**: Ajusta DPI, formato y opciones
4. 🌐 **Integrar**: Usa como servidor MCP o en tu código Go
5. 🚀 **Contribuir**: ¡Las PRs son bienvenidas!

## Ayuda

```bash
pdf2img --help
pdf2img info --help
mcp-server --help
```

---

**¿Preguntas?** Ver [README.md](README.md) para documentación completa.
