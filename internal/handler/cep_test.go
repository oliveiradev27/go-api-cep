package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
	"github.com/leandrodeoliveiranovais/api-cep/internal/handler"
	"github.com/leandrodeoliveiranovais/api-cep/internal/validation"
)

// mockCEPService é um double de teste para o serviço de CEP.
type mockCEPService struct {
	getFunc func(ctx context.Context, cep string) (domain.Address, error)
}

func (m *mockCEPService) GetAddress(ctx context.Context, cep string) (domain.Address, error) {
	return m.getFunc(ctx, cep)
}

// TestCEPHandler_GetCEP_Success valida resposta HTTP 200 com JSON do contrato ViaCEP.
func TestCEPHandler_GetCEP_Success(t *testing.T) {
	mockSvc := &mockCEPService{
		getFunc: func(_ context.Context, cep string) (domain.Address, error) {
			return domain.Address{
				CEP:         "01001-000",
				Logradouro:  "Praça da Sé",
				Complemento: "lado ímpar",
				Unidade:     "",
				Bairro:      "Sé",
				Localidade:  "São Paulo",
				UF:          "SP",
				Estado:      "São Paulo",
				Regiao:      "Sudeste",
				IBGE:        "3550308",
				GIA:         "1004",
				DDD:         "11",
				SIAFI:       "7107",
			}, nil
		},
	}

	h := handler.NewCEPHandler(mockSvc)
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cep/01001000", nil)
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

	if result["cep"] != "01001-000" {
		t.Errorf("cep: esperava %q, obteve %q", "01001-000", result["cep"])
	}
	if result["localidade"] != "São Paulo" {
		t.Errorf("localidade: esperava %q, obteve %q", "São Paulo", result["localidade"])
	}
}

// TestCEPHandler_GetCEP_InvalidCEP valida resposta 400 para CEP inválido.
func TestCEPHandler_GetCEP_InvalidCEP(t *testing.T) {
	mockSvc := &mockCEPService{
		getFunc: func(_ context.Context, _ string) (domain.Address, error) {
			t.Error("serviço não deveria ser chamado para CEP inválido no path")
			return domain.Address{}, nil
		},
	}

	h := handler.NewCEPHandler(mockSvc)
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cep/95010A10", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusBadRequest, rec.Code)
	}
}

// TestCEPHandler_GetCEP_NotFound valida resposta 404 para CEP inexistente.
func TestCEPHandler_GetCEP_NotFound(t *testing.T) {
	mockSvc := &mockCEPService{
		getFunc: func(_ context.Context, _ string) (domain.Address, error) {
			return domain.Address{}, domain.ErrNotFound
		},
	}

	h := handler.NewCEPHandler(mockSvc)
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cep/99999999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusNotFound, rec.Code)
	}
}

// TestCEPHandler_GetCEP_ServiceError valida resposta 502 para falha no upstream.
func TestCEPHandler_GetCEP_ServiceError(t *testing.T) {
	mockSvc := &mockCEPService{
		getFunc: func(_ context.Context, _ string) (domain.Address, error) {
			return domain.Address{}, errors.New("upstream indisponível")
		},
	}

	h := handler.NewCEPHandler(mockSvc)
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cep/01001000", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusBadGateway, rec.Code)
	}
}

// TestCEPHandler_GetCEP_InvalidCEPFromService valida 400 quando serviço retorna ErrInvalidCEP.
func TestCEPHandler_GetCEP_InvalidCEPFromService(t *testing.T) {
	mockSvc := &mockCEPService{
		getFunc: func(_ context.Context, _ string) (domain.Address, error) {
			return domain.Address{}, validation.ErrInvalidCEP
		},
	}

	h := handler.NewCEPHandler(mockSvc)
	router := handler.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cep/01001000", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: esperava %d, obteve %d", http.StatusBadRequest, rec.Code)
	}
}
