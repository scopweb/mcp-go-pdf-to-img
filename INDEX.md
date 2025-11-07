# 📑 Índice del Proyecto pdf2img

## 🚀 Inicio rápido

### ¿Quiero empezar ya?
→ [QUICKSTART.md](QUICKSTART.md) 

### ¿Quiero entender el proyecto?
→ [SUMMARY.md](SUMMARY.md)

### ¿Quiero la documentación completa?
→ [README.md](README.md)

---

## 📚 Documentación por tema

### Para usuarios

| Documento | Contenido |
|-----------|----------|
| [QUICKSTART.md](QUICKSTART.md) | Instalación y primeros pasos |
| [README.md](README.md) | Documentación completa y referencia |
| [EXAMPLES.md](EXAMPLES.md) | Ejemplos prácticos de uso |

**Recomendado leer en orden:** QUICKSTART → README → EXAMPLES

### Para desarrolladores

| Documento | Contenido |
|-----------|----------|
| [DEVELOPMENT.md](DEVELOPMENT.md) | Setup, build, debugging |
| [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) | Arquitectura interna |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Cómo contribuir |

**Recomendado leer en orden:** PROJECT_STRUCTURE → DEVELOPMENT → CONTRIBUTING

### Referencia

| Documento | Contenido |
|-----------|----------|
| [LICENSE](LICENSE) | Licencia Apache 2.0 |
| [SUMMARY.md](SUMMARY.md) | Resumen del proyecto |

---

## 🎯 Por caso de uso

### "Necesito convertir un PDF a imágenes"
1. [QUICKSTART.md](QUICKSTART.md) - Instalación
2. Ejecutar: `pdf2img -i documento.pdf -o ./output`

### "Quiero entender cómo funciona"
1. [SUMMARY.md](SUMMARY.md) - Resumen general
2. [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Arquitectura
3. Leer el código en `pkg/converter/converter.go`

### "Quiero usar el servidor MCP con Claude"
1. [README.md](README.md#servidor-mcp) - Sección MCP
2. [EXAMPLES.md](EXAMPLES.md#mcp-server---ejemplos-de-integración) - Ejemplos MCP

### "Quiero contribuir al proyecto"
1. [DEVELOPMENT.md](DEVELOPMENT.md) - Setup local
2. [CONTRIBUTING.md](CONTRIBUTING.md) - Guía de contribución
3. [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Entender el código

### "Necesito información específica"
1. [README.md](README.md#solución-de-problemas) - Troubleshooting
2. [EXAMPLES.md](EXAMPLES.md#solución-de-problemas-comunes) - Problemas comunes

---

## 📂 Estructura de directorios

```
pdf2img/
├── cmd/                          # Aplicaciones
│   ├── pdf2img/main.go          # CLI
│   └── mcp-server/main.go       # Servidor MCP
├── mcp/                          # Lógica MCP
├── pkg/converter/                # Lógica de conversión
├── [Documentación]               # 8 archivos .md
├── go.mod, Makefile, etc.
└── LICENSE
```

Ver [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) para detalles completos.

---

## 🔍 Búsqueda rápida

### Instalación
- [QUICKSTART.md](QUICKSTART.md)
- [README.md#instalación](README.md#instalación)

### Uso del CLI
- [README.md#uso](README.md#uso)
- [EXAMPLES.md#cli---ejemplos-prácticos](EXAMPLES.md#cli---ejemplos-prácticos)

### Servidor MCP
- [README.md#servidor-mcp](README.md#servidor-mcp)
- [EXAMPLES.md#mcp-server](EXAMPLES.md#mcp-server---ejemplos-de-integración)

### Troubleshooting
- [README.md#solución-de-problemas](README.md#solución-de-problemas)
- [EXAMPLES.md#solución-de-problemas-comunes](EXAMPLES.md#solución-de-problemas-comunes)

### Desarrollo
- [DEVELOPMENT.md](DEVELOPMENT.md)
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)

### Contribución
- [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📊 Contenido por archivo

### README.md (5.6 KB)
- Características principales
- Instalación (3 métodos)
- Uso del CLI
- Servidor MCP
- Dependencias
- Roadmap
- Troubleshooting

### QUICKSTART.md (3.2 KB)
- Instalación rápida (3 opciones)
- Primeros pasos
- Casos comunes
- Problemas rápidos

### EXAMPLES.md (6.1 KB)
- CLI: 6 ejemplos prácticos
- MCP HTTP: Ejemplos completos
- Desde código Go
- Casos de uso reales
- Troubleshooting avanzado

### DEVELOPMENT.md (4.4 KB)
- Setup del entorno
- Estructura del código
- Compilación y pruebas
- Debugging
- Extensiones futuras

### PROJECT_STRUCTURE.md (6.6 KB)
- Árbol de directorios
- Descripción de componentes
- Flujo de datos
- Dependencias
- Convenciones de código

### CONTRIBUTING.md (5.7 KB)
- Código de conducta
- Cómo reportar bugs
- Enviar PRs
- Guía de estilo
- Proceso de review

### SUMMARY.md (7.4 KB)
- Resumen del proyecto
- Características completadas
- Estadísticas
- Documentación disponible
- Próximos pasos

---

## 🎯 Flujo de navegación recomendado

```
¿Primer uso?
    ├─→ QUICKSTART.md
    ├─→ Instalar
    ├─→ Probar con un PDF
    └─→ README.md para más detalles

¿Usar como desarrollador?
    ├─→ SUMMARY.md
    ├─→ PROJECT_STRUCTURE.md
    ├─→ DEVELOPMENT.md
    └─→ Leer el código

¿Contribuir?
    ├─→ PROJECT_STRUCTURE.md
    ├─→ DEVELOPMENT.md
    ├─→ CONTRIBUTING.md
    └─→ Crear un PR

¿Usar el MCP server?
    ├─→ README.md (sección MCP)
    ├─→ EXAMPLES.md (sección MCP)
    └─→ Integrar con Claude

¿Necesito ayuda?
    ├─→ README.md (Troubleshooting)
    ├─→ EXAMPLES.md (Problemas comunes)
    └─→ Abrir un issue
```

---

## 📞 Navegación rápida

### 🔧 Instalación
- [Unix/Linux/macOS](QUICKSTART.md#opción-1-linux--macos)
- [Windows](QUICKSTART.md#opción-2-windows)
- [Go directo](QUICKSTART.md#opción-3-go-directo-cualquier-so)

### 💻 Comandos
- [CLI básico](README.md#cli---conversión-básica)
- [CLI avanzado](EXAMPLES.md)
- [MCP server](README.md#servidor-mcp)

### 🚀 Casos de uso
- [Miniaturas](EXAMPLES.md#ejemplo-1-generar-miniaturas-de-pdfs-en-lote)
- [Galería web](EXAMPLES.md#ejemplo-2-convertir-pdf-a-galería-de-imágenes-web)
- [Integración web](EXAMPLES.md#ejemplo-3-integración-con-servicio-web)

### 📚 Conceptos
- [Arquitectura](PROJECT_STRUCTURE.md#flujo-de-datos)
- [Convenciones](PROJECT_STRUCTURE.md#convenciones)
- [Performance](PROJECT_STRUCTURE.md#performance)

---

## 🆘 SOS - Necesito ayuda

| Problema | Solución |
|----------|----------|
| No sé por dónde empezar | [QUICKSTART.md](QUICKSTART.md) |
| Tengo un error | [Troubleshooting](README.md#solución-de-problemas) |
| Quiero entender el código | [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) |
| Tengo una pregunta | [EXAMPLES.md](EXAMPLES.md) o [README.md](README.md) |
| Quiero contribuir | [CONTRIBUTING.md](CONTRIBUTING.md) |
| El servidor no funciona | [DEVELOPMENT.md](DEVELOPMENT.md#debugging) |

---

## 📈 Estadísticas

- **21** archivos en total
- **6** archivos Go
- **7** archivos de documentación
- **~600+** líneas de código
- **25+** KB de documentación
- **Licencia**: Apache 2.0

---

## ✨ Características destacadas

- ✅ CLI completa
- ✅ Servidor MCP
- ✅ Documentación exhaustiva
- ✅ Código limpio y testeable
- ✅ Instalación fácil
- ✅ Soporte Windows/Unix
- ✅ Control de DPI
- ✅ Rango de páginas

---

**Última actualización**: 2025-11-06

Ver [SUMMARY.md](SUMMARY.md) para un resumen completo del proyecto.
