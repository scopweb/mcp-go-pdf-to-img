# 🔧 Troubleshooting - MCP Server con Claude Desktop

Guía para resolver problemas comunes cuando configuras pdf2img con Claude Desktop.

---

## ✅ Verificación Previa (Antes de Continuar)

Antes de intentar resolver problemas, verifica estos puntos:

### 1. MCP Server compilado
```bash
# Desde la raíz del proyecto
ls -la mcp-server.exe
# Debería tener ~18 MB
```

### 2. Archivo de configuración correcto
```bash
# Abre el archivo (Windows)
%APPDATA%\Claude\claude_desktop_config.json

# Verifica que veas algo como:
{
  "mcpServers": {
    "pdf2img": {
      "command": "C:\\ruta\\a\\mcp-server.exe",
      "args": ["--stdio"]
    }
  }
}
```

### 3. Ruta correcta
- ✅ Usa la ruta **COMPLETA** a `mcp-server.exe`
- ✅ Reemplaza `C:\ruta\a` con tu ruta real
- ❌ NO uses rutas relativas como `./mcp-server.exe`
- ❌ NO uses `~` (tilde) en Windows

### 4. Archivo guardado
Después de editar `claude_desktop_config.json`:
- Presiona Ctrl+S
- Cierra el editor
- Verifica que el archivo tenga los cambios

---

## 🐛 Problemas Comunes y Soluciones

### Problema 1: "pdf2img no aparece conectado"

**Síntomas**:
- No ves "pdf2img" en la esquina inferior de Claude Desktop
- O ves "pdf2img" pero con estado "desconectado"

**Soluciones** (en orden):

1. **Reinicia Claude Desktop completamente**
   ```bash
   # Cierra todas las ventanas de Claude Desktop
   taskkill /F /IM claude.exe
   # Espera 2 segundos
   # Abre Claude Desktop nuevamente
   ```

2. **Verifica la configuración nuevamente**
   ```json
   // ✅ Correcto
   {
     "command": "C:\\Users\\DAVID\\Projects\\pdf2img\\mcp-server.exe",
     "args": ["--stdio"]
   }

   // ❌ Incorrecto
   {
     "command": "mcp-server.exe",  // Falta ruta completa
     "args": ["--stdio"]
   }
   ```

3. **Prueba el servidor manualmente**
   ```bash
   cd C:\ruta\a\tu\proyecto
   mcp-server.exe --stdio
   # Debería decir: "MCP Server running in stdio mode"
   # Presiona Ctrl+C para salir
   ```

4. **Busca errores en el log**
   - En Claude Desktop, ve a Settings (⚙️) → Developer
   - Busca mensajes de error en rojo
   - Copia los mensajes de error exactos

### Problema 2: "MCP Server connection failed"

**Síntomas**:
- Ves un error de conexión
- El servidor no responde a Claude

**Soluciones**:

1. **Verifica que el archivo existe**
   ```bash
   dir "C:\ruta\exacta\mcp-server.exe"
   # Si no existe, compila nuevamente:
   cd C:\tu\proyecto
   go build -o mcp-server.exe ./cmd/mcp-server
   ```

2. **Verifica permisos**
   - Click derecho en `mcp-server.exe`
   - Propiedades → General
   - ¿Hay algún botón "Desbloquear"?
   - Si sí, haz click en Desbloquear
   - Haz click en OK
   - Reinicia Claude Desktop

3. **Verifica que Go está instalado**
   ```bash
   go version
   # Debería mostrar algo como: go version go1.21.5 windows/amd64
   ```

### Problema 3: "JSON-RPC errors" en los logs

**Síntomas**:
- En los logs ves errores de "invalid_literal", "Unrecognized key(s)"
- El servidor responde pero Claude no acepta las respuestas

**✅ SOLUCIONADO EN v2.0.0**:
El servidor ahora usa `github.com/mark3labs/mcp-go` (SDK oficial) que implementa correctamente JSON-RPC 2.0.

**Qué cambió**:
```
❌ Antes: Implementación manual de JSON-RPC 2.0
✅ Ahora: SDK oficial que maneja todo el protocolo automáticamente
```

**Para obtener la versión corregida**:
```bash
cd C:\tu\proyecto
go build -o mcp-server.exe ./cmd/mcp-server
taskkill /F /IM claude.exe
# Abre Claude Desktop nuevamente
```

**Si aún ves errores**:
- Verifica que compilaste la última versión (timestamp reciente)
- Ejecuta `go mod tidy` si hay problemas de compilación
- Reporta el error exacto del log con los detalles completos

### Problema 4: "Tool execution error" en Claude

**Síntomas**:
- Claude conecta al MCP
- Ves las herramientas disponibles
- Pero al usar una herramienta, sale error

**Causas y soluciones**:

1. **Ruta del PDF incorrecta**
   ```
   ❌ "documento.pdf"  (ruta relativa)
   ✅ "C:\Users\DAVID\Documents\documento.pdf"  (ruta absoluta)
   ```

2. **Directorio de salida no existe**
   - El servidor debería crearlo automáticamente
   - Pero verifica que tienes permisos de escritura en ese directorio

3. **Archivo PDF corrupto**
   ```bash
   # Prueba con el PDF de ejemplo
   # Desde Claude: "¿Cuántas páginas tiene example.pdf?"
   ```

4. **DPI demasiado alto**
   - Si usas DPI > 300, puede haber problemas de memoria
   - Intenta con DPI más bajo (150 o 200)

### Problema 5: "File not found"

**Síntomas**:
- Error: "input file not found"
- O: "PDF file not found"

**Soluciones**:

1. **Verifica que el archivo existe**
   ```bash
   ls "C:\Users\DAVID\Documents\documento.pdf"
   ```

2. **Usa ruta absoluta**
   - Pregunta a Claude: "Usa la ruta completa: C:\Users\DAVID\Documents\documento.pdf"

3. **Verifica permisos**
   - El archivo debe ser legible
   - No debe estar abierto en otro programa

### Problema 6: "Timeout" o "Server not responding"

**Síntomas**:
- Claude inicia la operación pero nunca responde
- Después de esperar, sale timeout

**Soluciones**:

1. **Verifica que el servidor está corriendo**
   ```bash
   tasklist | findstr mcp-server
   # Debería mostrar: mcp-server.exe ...
   ```

2. **Reinicia Claude Desktop**
   ```bash
   taskkill /F /IM claude.exe
   # Espera 3 segundos
   # Abre Claude Desktop nuevamente
   ```

3. **Verifica que el PDF no es muy grande**
   - PDFs > 50 MB pueden tomar tiempo
   - Intenta con un PDF más pequeño

---

## 📊 Checklist de Diagnóstico

Si nada funciona, ve por este checklist:

- [ ] ¿Compilé el `mcp-server.exe` correctamente?
  ```bash
  go build -o mcp-server.exe ./cmd/mcp-server
  ```

- [ ] ¿El archivo `mcp-server.exe` existe?
  ```bash
  ls -lh mcp-server.exe
  ```

- [ ] ¿La configuración en Claude Desktop es correcta?
  - Ruta completa (no relativa)
  - Sin caracteres especiales
  - `"args": ["--stdio"]` exacto

- [ ] ¿Reinicié Claude Desktop después de cambiar config?
  ```bash
  taskkill /F /IM claude.exe
  ```

- [ ] ¿El servidor funciona manualmente?
  ```bash
  mcp-server.exe --stdio
  ```

- [ ] ¿Veo "pdf2img" en la esquina de Claude Desktop?
  - Esquina inferior derecha
  - O Settings → Developer → MCP Servers

- [ ] ¿Puedo listar herramientas en Claude?
  - Pregunta: "¿Qué herramientas tengo disponibles?"
  - Debería mencionar `pdf_info` y `pdf_to_images`

- [ ] ¿Puedo usar `pdf_info`?
  - Pregunta: "¿Cuántas páginas tiene example.pdf?"
  - Debería retornar: "1 page"

---

## 🆘 Cuando Nada Funciona

Si ya probaste todo, aquí está qué reportar:

1. **Abre una terminal**
   ```bash
   cd C:\tu\proyecto
   mcp-server.exe --stdio
   ```

2. **Toma screenshot de:**
   - La salida del servidor
   - El log de Claude Desktop (Settings → Developer)

3. **Incluye información de:**
   - Sistema operativo: `ver` (en Windows)
   - Versión de Go: `go version`
   - Ruta del proyecto: `cd` (mostrar en terminal)

4. **Crea un issue en GitHub**
   - Incluye todo lo anterior
   - Describe qué probaste
   - Incluye mensajes de error exactos

---

## 📝 Logs Útiles

### Ver logs del servidor en tiempo real
```bash
# Terminal 1: Inicia el servidor
mcp-server.exe --stdio

# Terminal 2: Intenta usar desde Claude
# (Ve los logs en Terminal 1)
```

### Guardar logs para análisis
```bash
# Redirige stderr a un archivo
mcp-server.exe --stdio 2> mcp-errors.log

# Luego abre el archivo
cat mcp-errors.log
```

---

## ✅ Verificación Exitosa

Cuando funciona, verás:

1. **En Claude Desktop**
   - ✓ "pdf2img" conectado en la esquina
   - ✓ Sin errores rojos en Settings → Developer

2. **Al usar una herramienta**
   - ✓ Claude inicia la operación
   - ✓ La herramienta se ejecuta
   - ✓ Ves el resultado en menos de 10 segundos

3. **En el log del servidor**
   ```
   ✓ Server started and connected successfully
   ✓ Received method: initialize
   ✓ Received method: tools/list
   ✓ Received method: tools/call
   ```

---

**Última actualización**: 2025-11-07
**Versión**: 1.0.0
