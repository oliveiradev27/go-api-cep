package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/client"
	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
)

// TestViaCEPClient_FetchByCEP_Success valida a integração do client HTTP
// com um servidor mock que simula o ViaCEP.
func TestViaCEPClient_FetchByCEP_Success(t *testing.T) {
	// Arrange: servidor HTTP mock retornando JSON do ViaCEP.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Errorf("path inesperado: %s", r.URL.Path)
		}
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
	defer server.Close()

	viaClient := client.NewViaCEPClient(server.Client(), server.URL)

	// Act: consulta o CEP.
	addr, err := viaClient.FetchByCEP(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("esperava sucesso, obteve erro: %v", err)
	}

	// Assert: dados mapeados corretamente na struct.
	if addr.CEP != "01001-000" {
		t.Errorf("CEP: esperava %q, obteve %q", "01001-000", addr.CEP)
	}
	if addr.Localidade != "São Paulo" {
		t.Errorf("Localidade: esperava %q, obteve %q", "São Paulo", addr.Localidade)
	}
}

// TestViaCEPClient_FetchByCEP_NotFound valida tratamento de CEP inexistente.
func TestViaCEPClient_FetchByCEP_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"erro": "true"}`))
	}))
	defer server.Close()

	viaClient := client.NewViaCEPClient(server.Client(), server.URL)

	_, err := viaClient.FetchByCEP(context.Background(), "99999999")
	if err == nil {
		t.Fatal("esperava erro para CEP inexistente")
	}

	if !domain.IsNotFoundError(err) {
		t.Errorf("esperava ErrNotFound, obteve: %v", err)
	}
}

// TestViaCEPClient_FetchByCEP_BadRequest valida propagação de erro 400 do ViaCEP.
func TestViaCEPClient_FetchByCEP_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	viaClient := client.NewViaCEPClient(server.Client(), server.URL)

	_, err := viaClient.FetchByCEP(context.Background(), "invalido")
	if err == nil {
		t.Fatal("esperava erro para requisição inválida")
	}
}

// TestViaCEPClient_FetchByCEP_UpstreamError valida tratamento de falha no upstream.
func TestViaCEPClient_FetchByCEP_UpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	viaClient := client.NewViaCEPClient(server.Client(), server.URL)

	_, err := viaClient.FetchByCEP(context.Background(), "01001000")
	if err == nil {
		t.Fatal("esperava erro para falha no upstream")
	}
}
