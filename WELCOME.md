# 🎉 ¡Bienvenido a pdf2img!

## ¿Qué es pdf2img?

**pdf2img** es una herramienta completa en Go para convertir páginas de PDF a imágenes PNG o JPG. Incluye tanto una **interfaz CLI** como un **servidor MCP** para integración con Claude y otros sistemas.

- ✅ **Renderizado de alta calidad** usando go-pdfium (basado en PDFium de Google)
- 💻 **CLI completa** con múltiples opciones
- 🔌 **Servidor MCP** para integración con IA
- 📚 **Documentación exhaustiva**
- 🚀 **Listo para producción**

## 🚀 Empezar en 5 minutos

### Opción 1: Windows
```cmd
install.bat
pdf2img.exe -i documento.pdf -o .\output
```

### Opción 2: Linux/macOS
```bash
chmod +x install.sh
./install.sh
./pdf2img -i documento.pdf -o ./output
```

### Opción 3: Con Go
```bash
go build -o pdf2img ./cmd/pdf2img
./pdf2img -i documento.pdf -o ./output
```

**¡Listo!** Tus imágenes están en `./output/`

## 📚 Documentación

### Inicio rápido
- **[QUICKSTART.md](QUICKSTART.md)** - Primeros pasos en 5 minutos

### Aprender a usar
- **[README.md](README.md)** - Documentación completa
- **[EXAMPLES.md](EXAMPLES.md)** - Ejemplos prácticos
- **[INDEX.md](INDEX.md)** - Índice de navegación

### Para desarrolladores
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Arquitectura del proyecto
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Guía de desarrollo
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Cómo contribuir

### Referencia
- **[SUMMARY.md](SUMMARY.md)** - Resumen del proyecto
- **[LICENSE](LICENSE)** - Apache 2.0

## 💡 Ejemplos rápidos

### Convertir todas las páginas
```bash
pdf2img -i documento.pdf -o ./output
```

### Alta resolución (300 DPI)
```bash
pdf2img -i documento.pdf -o ./output -d 300
```

### Formato JPG
```bash
pdf2img -i documento.pdf -o ./output -f jpg
```

### Solo primeras 5 páginas
```bash
pdf2img -i documento.pdf -o ./output --start 1 --end 5
```

### Información del PDF
```bash
pdf2img info documento.pdf
```

## 🔌 Usar como servidor MCP

```bash
# Modo HTTP
mcp-server

# Modo stdio (Claude)
mcp-server --stdio
```

## 🎯 Características principales

| Característica | Detalles |
|---|---|
| **Formatos** | PNG, JPG |
| **DPI** | 72 a 600+ configurable |
| **Rango de páginas** | Seleccionar rango específico |
| **CLI** | 10+ flags y opciones |
| **MCP Server** | HTTP y stdio |
| **Tests** | Unitarios incluidos |
| **Licencia** | Apache 2.0 |
| **Plataformas** | Windows, Linux, macOS |

## 🛠️ Tecnologías

- **Go 1.21+** - Lenguaje de programación
- **go-pdfium** - Renderización de PDF
- **Cobra** - Framework CLI
- **PDFium** - Motor de renderización

## 📞 ¿Necesitas ayuda?

### Por donde empezar
1. Lee [QUICKSTART.md](QUICKSTART.md) si es tu primer uso
2. Instala con los scripts incluidos
3. Prueba con tu primer PDF
4. Explora [EXAMPLES.md](EXAMPLES.md) para más casos

### Problemas comunes
- Ver [README.md#solución-de-problemas](README.md#solución-de-problemas)
- Ver [EXAMPLES.md#solución-de-problemas-comunes](EXAMPLES.md#solución-de-problemas-comunes)

### Desarrollo
- Leer [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)
- Leer [DEVELOPMENT.md](DEVELOPMENT.md)
- Ejecutar `make test` para verificar

### Contribución
- Ver [CONTRIBUTING.md](CONTRIBUTING.md)
- Seguir las guías de estilo
- Crear un PR

## ✨ Lo que hace pdf2img

### Antes (PDF)
```
documento.pdf (25 páginas)
```

### Después (Imágenes)
```
page_0001.png
page_0002.png
...
page_0025.png
```

## 🚀 Casos de uso

✅ **Convertir documentos** a imágenes para almacenamiento
✅ **Generar miniaturas** para vistas previas
✅ **Extraer páginas específicas** de un PDF
✅ **Alta resolución para impresión** con DPI 300+
✅ **Integración con Claude** via MCP
✅ **Automatizar procesos** con scripts

## 📊 Estructura del proyecto

```
pdf2img/
├── cmd/pdf2img/           → CLI
├── cmd/mcp-server/        → Servidor MCP
├── pkg/converter/         → Lógica compartida
├── Documentación (8 .md)
├── Scripts (install.sh/.bat)
└── Configuración
```

## 🎓 Aprenderás

- Cómo usar go-pdfium para renderizar PDFs
- Cómo crear CLI profesionales con Cobra
- Cómo implementar servidores MCP
- Best practices en Go
- Arquitectura de software limpio

## 🏆 Destacados

- 📚 **Documentación completa** (8 guías)
- 🧪 **Tests unitarios** incluidos
- 🔧 **Scripts de instalación** (Windows + Unix)
- ✨ **Código limpio** y profesional
- 🚀 **Listo para producción**
- 📜 **Licencia libre** (Apache 2.0)

## 📈 Roadmap

- [ ] Soporte para PDF con contraseña
- [ ] Procesamiento paralelo
- [ ] OCR de imágenes
- [ ] Caché de renderizado
- [ ] Más formatos (WebP, TIFF)

## 🎯 Próximo paso

**¿Listo para empezar?**

→ [QUICKSTART.md](QUICKSTART.md) (5 minutos)

**¿Quieres entender todo?**

→ [INDEX.md](INDEX.md) (Navegación completa)

**¿Quieres el resumen?**

→ [SUMMARY.md](SUMMARY.md) (Resumen del proyecto)

---

## 📝 Información del proyecto

| Aspecto | Detalles |
|---|---|
| **Versión** | 1.0.0 |
| **Licencia** | Apache 2.0 |
| **Estado** | Completo y listo |
| **Plataformas** | Windows, Linux, macOS |
| **Go version** | 1.21+ |
| **Dependencias** | go-pdfium, Cobra |

---

## 🎉 ¡Comienza ahora!

1. **Lee** [QUICKSTART.md](QUICKSTART.md) (5 min)
2. **Instala** con `./install.sh` o `install.bat`
3. **Prueba** con `pdf2img -i documento.pdf -o ./output`
4. **Explora** [EXAMPLES.md](EXAMPLES.md) para más

---

**Creado con ❤️ usando Go y PDFium**

¡Disfruta convertiendo tus PDFs a imágenes! 📸

Preguntas? Ver [INDEX.md](INDEX.md) para navegación completa.
