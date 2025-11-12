@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul

REM Color output
cls
echo.
echo ╔══════════════════════════════════════════════════════════════╗
echo ║           PDF2IMG - Conversor Interactivo de PDF             ║
echo ║                                                              ║
echo ║  Convierte páginas de PDF a imágenes PNG o JPG              ║
echo ╚══════════════════════════════════════════════════════════════╝
echo.

REM Get PDF information
:input_pdf_path
set "pdf_path="
echo.
echo [1] Ingresa la ruta del archivo PDF:
echo    Ejemplo: C:\documentos\archivo.pdf
echo    O deja en blanco para usar un PDF en la carpeta actual
set /p pdf_path="Ruta del PDF: "

if "%pdf_path%"=="" (
    echo.
    echo Buscando PDFs en la carpeta actual...
    dir /b *.pdf 2>nul
    if !errorlevel! equ 0 (
        echo.
        set /p pdf_path="Ingresa el nombre del PDF: "
    ) else (
        echo No se encontraron PDFs en la carpeta actual.
        goto input_pdf_path
    )
)

REM Validate PDF exists
if not exist "!pdf_path!" (
    echo.
    echo ❌ Error: El archivo no existe: !pdf_path!
    goto input_pdf_path
)

REM Get PDF info
echo.
echo 📋 Analizando PDF...
for %%F in ("!pdf_path!") do (
    set "pdf_name=%%~nF"
    set "pdf_size=%%~zF"
)

REM Display file info
echo ✓ Archivo encontrado: !pdf_name!

REM Get page count using pdf2img info command
echo.
echo [2] Obteniendo información del PDF...
for /f "tokens=*" %%A in ('pdf2img.exe info "!pdf_path!" 2^>nul ^| findstr /R "^Pages:"') do (
    set "line=%%A"
    for /f "tokens=2" %%B in ("!line!") do set "total_pages=%%B"
)

if "!total_pages!"=="" (
    echo ⚠️  No se pudo obtener el número de páginas automáticamente
    set /p total_pages="Ingresa el número de páginas del PDF: "
)

echo ✓ Total de páginas: !total_pages!

REM Recommend settings based on page count
echo.
echo [3] Recomendaciones automáticas según número de páginas:
echo.

set "recommended_dpi=150"
set "recommended_format=png"
set "recommended_pool=2"
set "recommended_refresh=50"

if !total_pages! geq 300 (
    set "recommended_dpi=150"
    set "recommended_format=jpg"
    set "recommended_pool=4"
    set "recommended_refresh=25"
    echo 📚 PDF MUY GRANDE (300+ páginas)
    echo    DPI recomendado: !recommended_dpi! ^(balance^)
    echo    Formato: !recommended_format! ^(menor tamaño^)
    echo    Pool size: !recommended_pool! ^(máximo rendimiento^)
    echo    Refresh: !recommended_refresh! ^(refresco frecuente^)
    echo.
) else if !total_pages! geq 200 (
    set "recommended_dpi=150"
    set "recommended_format=jpg"
    set "recommended_pool=3"
    set "recommended_refresh=50"
    echo 📚 PDF Grande (200-299 páginas)
    echo    DPI recomendado: !recommended_dpi! ^(balance^)
    echo    Formato: !recommended_format! ^(menor tamaño^)
    echo    Pool size: !recommended_pool! ^(buen rendimiento^)
    echo    Refresh: !recommended_refresh! ^(cada 50 págs^)
    echo.
) else if !total_pages! geq 100 (
    set "recommended_dpi=150"
    set "recommended_format=png"
    set "recommended_pool=2"
    set "recommended_refresh=50"
    echo 📖 PDF Mediano (100-199 páginas)
    echo    DPI recomendado: !recommended_dpi! ^(balance^)
    echo    Formato: !recommended_format! ^(buena calidad^)
    echo    Pool size: !recommended_pool! ^(estándar^)
    echo    Refresh: !recommended_refresh! ^(cada 50 págs^)
    echo.
) else (
    set "recommended_dpi=150"
    set "recommended_format=png"
    set "recommended_pool=2"
    set "recommended_refresh=0"
    echo 📄 PDF Pequeño (menos de 100 páginas)
    echo    DPI recomendado: !recommended_dpi! ^(balance^)
    echo    Formato: !recommended_format! ^(buena calidad^)
    echo    Pool size: !recommended_pool! ^(estándar^)
    echo    Refresh: desactivado ^(no es necesario^)
    echo.
)

REM Ask for DPI
echo [4] Selecciona la resolución (DPI):
echo    72-96   = Web/Pantalla (rápido, baja calidad)
echo    150     = Estándar ^(balance recomendado^)
echo    200-300 = Imprenta/Alta calidad
echo    600+    = Muy detallado (lento)
echo.
set /p user_dpi="DPI [!recommended_dpi!]: "
if "!user_dpi!"=="" set "user_dpi=!recommended_dpi!"

REM Ask for format
echo.
echo [5] Selecciona el formato:
echo    png = Lossless, mejor calidad, archivos más grandes
echo    jpg = Comprimido, archivos más pequeños, rápido
echo.
set /p user_format="Formato [!recommended_format!]: "
if "!user_format!"=="" set "user_format=!recommended_format!"

REM Convert to lowercase
for %%A in (png jpg jpeg) do (
    if /i "!user_format!"=="%%A" set "user_format=%%A"
)

REM Ask for output directory
echo.
echo [6] ¿Dónde exportar las imágenes?
echo    Deja en blanco para crear carpeta "output" en la carpeta actual
set /p output_dir="Carpeta de salida [.\output]: "
if "!output_dir!"=="" set "output_dir=.\output"

REM Create output directory if doesn't exist
if not exist "!output_dir!" (
    mkdir "!output_dir!"
    echo ✓ Carpeta creada: !output_dir!
)

REM Ask about page range
echo.
echo [7] ¿Convertir todas las páginas?
echo    1 = Todas las páginas (default)
echo    2 = Rango específico
set /p page_choice="Opción [1]: "
if "!page_choice!"=="" set "page_choice=1"

set "start_page=0"
set "end_page=0"

if "!page_choice!"=="2" (
    echo.
    set /p start_page="Página inicial [1]: "
    if "!start_page!"=="" set "start_page=1"

    set /p end_page="Página final [!total_pages!]: "
    if "!end_page!"=="" set "end_page=!total_pages!
)

REM Ask for retry on failure
echo.
echo [8] ¿Reintentar páginas fallidas?
echo    1 = No (default, más rápido)
echo    2 = Sí (más lento, pero reintenta con DPI reducido)
set /p retry_choice="Opción [1]: "
if "!retry_choice!"=="" set "retry_choice=1"

set "retry_flag="
if "!retry_choice!"=="2" set "retry_flag=--retry"

REM Summary
echo.
echo.
echo ╔══════════════════════════════════════════════════════════════╗
echo ║                   RESUMEN DE CONFIGURACIÓN                   ║
echo ╚══════════════════════════════════════════════════════════════╝
echo.
echo  📄 PDF:              !pdf_name!
echo  📄 Páginas:          !total_pages!
echo  📊 DPI:              !user_dpi!
echo  🎨 Formato:          !user_format!
echo  📁 Salida:           !output_dir!

if "!start_page!"=="0" (
    echo  📖 Rango:            Todas las páginas
) else (
    echo  📖 Rango:            Páginas !start_page! a !end_page!
)

if "!retry_flag!"=="" (
    echo  🔄 Reintentar:       No
) else (
    echo  🔄 Reintentar:       Sí
)

echo.
set /p confirm="¿Proceder con la conversión? [S/n]: "
if /i "!confirm!"=="n" goto end

REM Build command
set "command=pdf2img.exe -i "!pdf_path!" -o "!output_dir!" -f !user_format! -d !user_dpi!"

if "!total_pages!" geq 300 (
    set "command=!command! --pool-size 4 --refresh-every 25"
) else if "!total_pages!" geq 200 (
    set "command=!command! --pool-size 3 --refresh-every 50"
) else if "!total_pages!" geq 100 (
    set "command=!command! --pool-size 2 --refresh-every 50"
)

if "!start_page!" neq "0" (
    set "command=!command! --start !start_page! --end !end_page!"
)

if not "!retry_flag!"=="" (
    set "command=!command! !retry_flag!
)

REM Execute conversion
echo.
echo ═══════════════════════════════════════════════════════════════
echo 🚀 Iniciando conversión...
echo ═══════════════════════════════════════════════════════════════
echo.

%command%

REM Check result
if !errorlevel! equ 0 (
    echo.
    echo ═══════════════════════════════════════════════════════════════
    echo ✅ ¡Conversión completada exitosamente!
    echo ═══════════════════════════════════════════════════════════════
    echo.
    echo 📁 Archivos guardados en: !output_dir!
    echo.
    echo Abriendo carpeta de salida...
    start "" "!output_dir!"
) else (
    echo.
    echo ═══════════════════════════════════════════════════════════════
    echo ❌ Error durante la conversión
    echo ═══════════════════════════════════════════════════════════════
)

:end
echo.
pause
endlocal
