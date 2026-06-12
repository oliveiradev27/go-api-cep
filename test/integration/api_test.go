//go:build integration

package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// apiBaseURL é a URL base da API em execução (configurável via env).
var apiBaseURL = envOrDefault("API_BASE_URL", "http://localhost:8080")

// TestAPI_GetCEP_Valid consulta um CEP válido na API em execução.
func TestAPI_GetCEP_Valid(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(apiBaseURL + "/cep/01001000")
	if err != nil {
		t.Fatalf("falha na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status: esperava %d, obteve %d, body: %s", http.StatusOK, resp.StatusCode, string(body))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("resposta inválida: %v", err)
	}

	if result["cep"] != "01001-000" {
		t.Errorf("cep: esperava %q, obteve %q", "01001-000", result["cep"])
	}
}

// TestAPI_GetCEP_NotFound consulta um CEP inexistente na API em execução.
func TestAPI_GetCEP_NotFound(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(apiBaseURL + "/cep/99999999")
	if err != nil {
		t.Fatalf("falha na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusNotFound, resp.StatusCode)
	}
}

// TestAPI_GetCEP_InvalidCEP consulta um CEP com formato inválido.
func TestAPI_GetCEP_InvalidCEP(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(apiBaseURL + "/cep/95010A10")
	if err != nil {
		t.Fatalf("falha na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
