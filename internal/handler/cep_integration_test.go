package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/client"
	"github.com/leandrodeoliveiranovais/api-cep/internal/handler"
	"github.com/leandrodeoliveiranovais/api-cep/internal/service"
)

// TestIntegration_CEPFlow_Success valida o fluxo completo:
// Handler → Service → Client → Mock ViaCEP → resposta JSON.
func TestIntegration_CEPFlow_Success(t *testing.T) {
	// Arrange: mock server simulando o ViaCEP.
	mockViaCEP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"cep": "01001-000",
			"logradouro": "Praça da Sé",
			"complemento": "lado ímpar",
			"unidade": "",
			"bairro": "Sé",
			"localidade": "São Paulo",
			"uf": "SP",
			"estado": "São Paulo",
			"regiao": "Sudeste",
			"ibge": "3550308",
			"gia": "1004",
			"ddd": "11",
			"siafi": "7107"
		}`))
	}))
	defer mockViaCEP.Close()

	// Monta a cadeia completa de dependências (como em produção).
	viaClient := client.NewViaCEPClient(mockViaCEP.Client(), mockViaCEP.URL)
	svc := service.NewCEPService(viaClient)
	h := handler.NewCEPHandler(svc)
	router := handler.NewRouter(h)

	// Act: requisição HTTP real ao handler.
	req := httptest.NewRequest(http.MethodGet, "/cep/01001000", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert: resposta 200 com contrato ViaCEP completo.
	if rec.Code != http.StatusOK {
		t.Fatalf("status: esperava %d, obteve %d, body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var result map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("resposta inválida: %v", err)
	}

	if result["cep"] != "01001-000" {
		t.Errorf("cep: esperava %q, obteve %q", "01001-000", result["cep"])
	}
	if result["uf"] != "SP" {
		t.Errorf("uf: esperava %q, obteve %q", "SP", result["uf"])
	}
}

// TestIntegration_CEPFlow_NotFound valida fluxo completo para CEP inexistente.
func TestIntegration_CEPFlow_NotFound(t *testing.T) {
	mockViaCEP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"erro": "true"}`))
	}))
	defer mockViaCEP.Close()

	viaClient := client.NewViaCEPClient(mockViaCEP.Client(), mockViaCEP.URL)
	svc := service.NewCEPService(viaClient)
	h := handler.NewCEPHandler(svc)
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cep/99999999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusNotFound, rec.Code)
	}
}
