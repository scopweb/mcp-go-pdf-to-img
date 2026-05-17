package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/scopweb/mcp-go-pdf-to-img/internal/config"
	"github.com/scopweb/mcp-go-pdf-to-img/internal/logging"
	"github.com/scopweb/mcp-go-pdf-to-img/pkg/converter"
)

// Converter interface for PDF operations
type Converter interface {
	Convert(opts *converter.ConvertOptions) (*converter.ConvertResult, error)
	GetPDFInfo(pdfPath string) (map[string]interface{}, error)
	Close() error
}

// Handlers encapsula todos los handlers HTTP con sus dependencias.
type Handlers struct {
	converter Converter
	logger    logging.Logger
	config    *config.ServerConfig
}

// NewHandlers crea un nuevo set de handlers.
func NewHandlers(conv Converter, logger logging.Logger, config *config.ServerConfig) *Handlers {
	return &Handlers{
		converter: conv,
		logger:    logger,
		config:    config,
	}
}

// Health responde con status OK.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// ConvertRequest contiene los parámetros para conversión.
type ConvertRequest struct {
	DPI         int    `json:"dpi,omitempty"`
	Format      string `json:"format,omitempty"`
	StartPage   int    `json:"start_page,omitempty"`
	EndPage     int    `json:"end_page,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	RetryFailed bool   `json:"retry_failed,omitempty"`
	PoolSize    int    `json:"pool_size,omitempty"`
	RefreshEvery int   `json:"refresh_every,omitempty"`
}

// ConvertResponse contiene el resultado de la conversión.
type ConvertResponse struct {
	TotalPages   int      `json:"total_pages"`
	Successful   int      `json:"successful"`
	Failed        int      `json:"failed"`
	OutputFiles   []string `json:"output_files,omitempty"`
	Errors        []string `json:"errors,omitempty"`
	WarningPages []int    `json:"warning_pages,omitempty"`
}

// Convert convierte un PDF a imágenes y devuelve un ZIP con los resultados.
func (h *Handlers) Convert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(h.config.MaxUploadSize); err != nil {
		h.logger.Error("failed to parse multipart form", err)
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Warn("missing file field", slog.Any("error", err))
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crear archivo temporal para upload
	tmpFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		h.logger.Error("failed to create temp file", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Guardar archivo
	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		h.logger.Error("failed to save uploaded file", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	tmpFile.Close()

	// Parsear parámetros de conversión
	var req ConvertRequest
	if err := h.parseConvertParams(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Crear directorio temporal de salida
	tmpOutputDir, err := os.MkdirTemp("", "pdf2img-*")
	if err != nil {
		h.logger.Error("failed to create temp output dir", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpOutputDir)

	// Configurar opciones de conversión
	opts := &converter.ConvertOptions{
		InputPath:    tmpPath,
		OutputDir:    tmpOutputDir,
		Format:       req.Format,
		DPI:          float64(req.DPI),
		StartPage:    req.StartPage,
		EndPage:      req.EndPage,
		Prefix:       req.Prefix,
		RetryFailed:  req.RetryFailed,
		MaxPoolSize:  req.PoolSize,
		RefreshEvery: req.RefreshEvery,
	}

	// Convertir PDF
	result, err := h.converter.Convert(opts)
	if err != nil {
		h.logger.Error("PDF conversion failed", err)
		http.Error(w, "failed to convert PDF: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Preparar respuesta
	response := ConvertResponse{
		TotalPages:   result.TotalPages,
		Successful:   result.Successful,
		Failed:       result.Failed,
		OutputFiles:  result.OutputFiles,
		Errors:       result.Errors,
		WarningPages: result.WarningPages,
	}

	// Si no hay archivos de salida, devolver JSON con el resultado
	if len(result.OutputFiles) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Devolver ZIP con las imágenes
	h.serveZip(w, header.Filename, result.OutputFiles)
}

func (h *Handlers) parseConvertParams(r *http.Request, req *ConvertRequest) error {
	// Defaults
	req.Format = h.config.Convert.DefaultFormat
	if req.Format == "" {
		req.Format = "png"
	}
	req.DPI = h.config.Convert.DefaultDPI
	if req.DPI <= 0 {
		req.DPI = 150
	}
	req.PoolSize = h.config.Convert.DefaultPoolSize
	req.RefreshEvery = h.config.Convert.RefreshEvery
	req.RetryFailed = h.config.Convert.RetryOnFailure

	// Parsear DPI
	if dpiStr := r.FormValue("dpi"); dpiStr != "" {
		if dpi, err := strconv.Atoi(dpiStr); err == nil && dpi > 0 {
			req.DPI = dpi
		}
	}

	// Parsear formato
	if format := r.FormValue("format"); format != "" {
		req.Format = format
	}

	// Parsear páginas
	if startStr := r.FormValue("start_page"); startStr != "" {
		if start, err := strconv.Atoi(startStr); err == nil {
			req.StartPage = start
		}
	}
	if endStr := r.FormValue("end_page"); endStr != "" {
		if end, err := strconv.Atoi(endStr); err == nil {
			req.EndPage = end
		}
	}

	// Parsear prefijo
	req.Prefix = r.FormValue("prefix")
	if req.Prefix == "" {
		req.Prefix = "page_"
	}

	// Parsear retry
	if retryStr := r.FormValue("retry"); retryStr != "" {
		req.RetryFailed, _ = strconv.ParseBool(retryStr)
	}

	// Parsear pool size
	if poolStr := r.FormValue("pool_size"); poolStr != "" {
		if pool, err := strconv.Atoi(poolStr); err == nil && pool > 0 {
			req.PoolSize = pool
		}
	}

	// Parsear refresh every
	if refreshStr := r.FormValue("refresh_every"); refreshStr != "" {
		if refresh, err := strconv.Atoi(refreshStr); err == nil && refresh >= 0 {
			req.RefreshEvery = refresh
		}
	}

	return nil
}

func (h *Handlers) serveZip(w http.ResponseWriter, filename string, files []string) {
	// Crear archivo ZIP temporal
	zipFile, err := os.CreateTemp("", "result-*.zip")
	if err != nil {
		h.logger.Error("failed to create zip file", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	zipPath := zipFile.Name()
	zipFile.Close()
	defer os.Remove(zipPath)

	// Crear ZIP
	if err := h.createZip(zipPath, files); err != nil {
		h.logger.Error("failed to create zip", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Abrir archivo ZIP
	zipRead, err := os.Open(zipPath)
	if err != nil {
		h.logger.Error("failed to open zip", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer zipRead.Close()

	// Enviar respuesta
	w.Header().Set("Content-Type", "application/zip")
	zipName := sanitizeFilename(filename) + "-images.zip"
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipName))

	if _, err := io.Copy(w, zipRead); err != nil {
		h.logger.Error("error writing zip response", err)
	}
}

func (h *Handlers) createZip(zipPath string, files []string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	for _, filePath := range files {
		fw, err := zw.Create(filepath.Base(filePath))
		if err != nil {
			continue
		}
		fr, err := os.Open(filePath)
		if err != nil {
			continue
		}
		io.Copy(fw, fr)
		fr.Close()
	}
	return nil
}

// Info devuelve información sobre un PDF.
func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tmpPath string

	if r.Method == http.MethodPost {
		// Parse multipart form
		if err := r.ParseMultipartForm(h.config.MaxUploadSize); err != nil {
			h.logger.Error("failed to parse multipart form", err)
			http.Error(w, "invalid request format", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			h.logger.Warn("missing file field", slog.Any("error", err))
			http.Error(w, "missing file field", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Crear archivo temporal para upload
		tmpFile, err := os.CreateTemp("", "upload-*.pdf")
		if err != nil {
			h.logger.Error("failed to create temp file", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		tmpPath = tmpFile.Name()
		defer os.Remove(tmpPath)

		if _, err := io.Copy(tmpFile, file); err != nil {
			tmpFile.Close()
			h.logger.Error("failed to save uploaded file", err)
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}
		tmpFile.Close()
	} else {
		// GET: obtener URL del archivo como parámetro
		fileURL := r.FormValue("url")
		if fileURL == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		// Descargar archivo
		tmpFile, err := os.CreateTemp("", "upload-*.pdf")
		if err != nil {
			h.logger.Error("failed to create temp file", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		tmpPath = tmpFile.Name()
		defer os.Remove(tmpPath)

		resp, err := http.Get(fileURL)
		if err != nil {
			tmpFile.Close()
			h.logger.Error("failed to download file", err)
			http.Error(w, "failed to download file", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			tmpFile.Close()
			h.logger.Error("failed to save downloaded file", err)
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}
		tmpFile.Close()
	}

	// Obtener información del PDF
	info, err := h.converter.GetPDFInfo(tmpPath)
	if err != nil {
		h.logger.Error("failed to get PDF info", err)
		http.Error(w, "failed to get PDF info: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// sanitizeFilename limpia un nombre de archivo para evitar caracteres problemáticos.
func sanitizeFilename(filename string) string {
	if ext := filepath.Ext(filename); ext != "" {
		filename = filename[:len(filename)-len(ext)]
	}
	return filepath.Base(filename)
}