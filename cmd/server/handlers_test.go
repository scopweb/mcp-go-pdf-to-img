package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scopweb/mcp-go-pdf-to-img/internal/config"
	"github.com/scopweb/mcp-go-pdf-to-img/internal/logging"
	"github.com/scopweb/mcp-go-pdf-to-img/pkg/converter"
)

type convStub struct {
	conv *converter.Converter
}

func (c *convStub) Convert(opts *converter.ConvertOptions) (*converter.ConvertResult, error) {
	return &converter.ConvertResult{}, nil
}

func (c *convStub) GetPDFInfo(pdfPath string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (c *convStub) Close() error {
	return nil
}

func TestHealth(t *testing.T) {
	handlers := NewHandlers(&convStub{}, logging.New("error"), config.NewServerConfig())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handlers.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got '%s'", w.Body.String())
	}
}

func TestConvertMissingFile(t *testing.T) {
	handlers := NewHandlers(&convStub{}, logging.New("error"), config.NewServerConfig())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pdf/convert", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	handlers.Convert(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestInfoMissingFile(t *testing.T) {
	handlers := NewHandlers(&convStub{}, logging.New("error"), config.NewServerConfig())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pdf/info", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	handlers.Info(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestConvertMethodNotAllowed(t *testing.T) {
	handlers := NewHandlers(&convStub{}, logging.New("error"), config.NewServerConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pdf/convert", nil)
	w := httptest.NewRecorder()
	handlers.Convert(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestInfoMethodNotAllowed(t *testing.T) {
	handlers := NewHandlers(&convStub{}, logging.New("error"), config.NewServerConfig())

	req := httptest.NewRequest(http.MethodPut, "/api/v1/pdf/info", nil)
	w := httptest.NewRecorder()
	handlers.Info(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"document.pdf", "document"},
		{"../../etc/passwd", "passwd"},
		{"path/with spaces/file.pdf", "file"},
	}

	for _, tt := range tests {
		result := sanitizeFilename(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFilename(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestParseConvertParams(t *testing.T) {
	h := &Handlers{
		config: &config.ServerConfig{
			Convert: config.ConvertConfig{
				DefaultDPI:      150,
				DefaultFormat:   "png",
				DefaultPoolSize: 2,
				RefreshEvery:    50,
				RetryOnFailure:  false,
			},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	cr := &ConvertRequest{}
	h.parseConvertParams(req, cr)

	if cr.DPI != 150 {
		t.Errorf("expected DPI 150, got %d", cr.DPI)
	}
	if cr.Format != "png" {
		t.Errorf("expected format png, got %s", cr.Format)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/?dpi=300&format=jpg", nil)
	cr2 := &ConvertRequest{}
	h.parseConvertParams(req2, cr2)

	if cr2.DPI != 300 {
		t.Errorf("expected DPI 300, got %d", cr2.DPI)
	}
	if cr2.Format != "jpg" {
		t.Errorf("expected format jpg, got %s", cr2.Format)
	}
}

func urlencodedForm(values map[string]string) io.Reader {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	for k, v := range values {
		w.WriteField(k, v)
	}
	w.Close()
	return buf
}