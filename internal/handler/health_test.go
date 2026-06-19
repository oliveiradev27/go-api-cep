package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/handler"
)

// TestHealthCheck_ReturnsOK valida que GET /health responde 200 com status "ok".
func TestHealthCheck_ReturnsOK(t *testing.T) {
	h := handler.NewCEPHandler(&mockCEPService{})
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: esperava application/json, obteve %q", ct)
	}

	var result map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("status: esperava %q, obteve %q", "ok", result["status"])
	}
}
