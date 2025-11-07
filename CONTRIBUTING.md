# Contribuir a pdf2img

¡Gracias por tu interés en contribuir! Este documento explica cómo hacerlo.

## Código de Conducta

- Sé respetuoso con otros colaboradores
- Usa lenguaje inclusivo
- Acepta críticas constructivas
- Enfócate en el código, no en las personas

## Cómo contribuir

### Reportar bugs

1. Verifica que el bug no está reportado
2. Crea un issue con:
   - Descripción clara del problema
   - Pasos para reproducir
   - Comportamiento esperado vs actual
   - Sistema operativo y versión de Go
   - Ejemplo mínimo reproducible

```
**Descripción:**
El conversor falla con PDFs que contienen imágenes JPG.

**Pasos para reproducir:**
1. Descargar documento.pdf
2. Ejecutar: pdf2img -i documento.pdf -o ./output
3. Ver error

**Error:**
```
failed to render page: invalid image format
```

**Sistema:**
- OS: Windows 11
- Go: 1.21
```

### Sugerir mejoras

1. Abre un issue con etiqueta `enhancement`
2. Describe la mejora detalladamente
3. Explica el caso de uso
4. Sugiere una posible solución

### Enviar un Pull Request

#### Preparación

```bash
# Fork el repositorio
git clone https://github.com/tu-usuario/pdf2img.git
cd pdf2img

# Crear rama para tu feature
git checkout -b feature/tu-feature-name
```

#### Desarrollo

1. **Mantener código limpio**
   ```bash
   # Formatear código
   make fmt

   # Lint
   make lint
   ```

2. **Escribir tests**
   ```bash
   # Si agregaste nuevas funciones
   go test -v ./...

   # Verificar coverage
   go test -cover ./...
   ```

3. **Commits significativos**
   ```bash
   # Commits claros y atómicos
   git commit -m "Add feature: description"

   # Evitar
   git commit -m "Fix stuff"
   ```

#### Enviar PR

```bash
# Push a tu fork
git push origin feature/tu-feature-name

# Crear Pull Request en GitHub
# - Título claro
# - Descripción detallada
# - Referencia a issues relacionados
# - Screenshots si es UI
```

## Estructura de un buen PR

```markdown
## Descripción
Qué hace tu PR (2-3 líneas)

## Tipo de cambio
- [ ] Bug fix
- [ ] Nueva feature
- [ ] Mejora existente
- [ ] Documentación
- [ ] Refactoring

## Cambios realizados
- Cambio 1
- Cambio 2
- Cambio 3

## Testing
Cómo probar:
1. Paso 1
2. Paso 2

## Checklist
- [ ] Código formateado (make fmt)
- [ ] Tests pasando (make test)
- [ ] Lint pasando (make lint)
- [ ] Documentación actualizada
- [ ] Commits atómicos y claros
```

## Guías de estilo

### Go

```go
// ✓ Bueno
func convertPDF(path string) error {
    if err := validatePath(path); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // ...
    return nil
}

// ✗ Evitar
func ConvertPDF(path string) (err error) {
    err = validatePath(path)
    if err != nil {
        panic(err)
    }
    // ...
}
```

**Reglas:**
- Nombres en inglés
- 80 caracteres max por línea (flexible)
- Error handling explícito
- Defer para limpieza
- Documentar funciones públicas

### Commit Messages

```
Formato:
<tipo>(<scope>): <subject>

<body>

<footer>

Ejemplos:
feat(converter): add support for encrypted PDFs
fix(cli): correct DPI calculation
docs(readme): update examples
refactor(mcp): simplify server initialization
test(converter): add edge case tests

Tipos:
- feat: Nueva feature
- fix: Corrección de bug
- docs: Cambios en documentación
- style: Cambios de formato (no lógica)
- refactor: Refactoring de código
- test: Agregar o actualizar tests
- chore: Cambios de dependencias, tooling, etc.
```

### Documentación

- README: Público, claro y conciso
- Comentarios: Explica el "por qué", no el "qué"
- Nombres: Auto-explicativos
- Ejemplos: Código funcional y actualizado

```go
// ✓ Bueno
// renderPage renders a PDF page to an image.
// If DPI is 0, defaults to 150.
func (c *Converter) renderPage(page int, dpi float64) error {
    // ...
}

// ✗ Evitar
// render page
func render(p int, d float64) error {
    // ...
}
```

## Proceso de review

1. Verificamos:
   - Código limpio y bien estructurado
   - Tests pasando
   - Documentación actualizada
   - Commits significativos

2. Podemos pedir:
   - Cambios en el código
   - Más tests
   - Actualizar documentación
   - Refactorizar

3. Una vez aprobado:
   - Merge automático
   - Agradecimiento en CHANGELOG (próximas releases)

## Desarrollo local

### Setup

```bash
git clone https://github.com/tu-usuario/pdf2img.git
cd pdf2img
go mod download
```

### Build

```bash
make build          # Compilar CLI
make build-dev      # Con símbolos de debug
make test          # Ejecutar tests
make fmt           # Formatear código
make lint          # Lint
make clean         # Limpiar
```

### Debug

```bash
# Con logs
./pdf2img -i documento.pdf -o ./output -v

# Profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

## Áreas donde podemos ayuda

### Alto impacto
- Performance improvements
- Nuevos formatos de salida
- Mejor manejo de errores
- Documentación

### Mediano impacto
- Tests adicionales
- Ejemplos mejorados
- Refactoring
- CI/CD improvements

### Baja prioridad
- Cambios cosméticos
- Comentarios sin valor
- Cambios de nombres sin contexto

## Preguntas frecuentes

**¿Puedo trabajar en una feature que está en el roadmap?**
Abre un issue primero para coordinar.

**¿Cuánto tiempo tarda un review?**
Usualmente 3-5 días, máximo 2 semanas.

**¿Mi PR fue rechazada. ¿Qué hago?**
Lée los comentarios, aprende de ellos. Podés ajustar y resubmitir.

**¿Necesito firmar algo?**
No. Solo mantén la licencia Apache 2.0.

## Agradecimientos

Valoramos cada contribución, grande o pequeña. Si tu PR fue mergeado, aparecerás en:
- Historial de commits
- CONTRIBUTORS.md (próximamente)
- Release notes

---

**¿Preguntas?** Abre un issue o contacta a los maintainers.

¡Gracias por contribuir a pdf2img!
