package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
	"github.com/leandrodeoliveiranovais/api-cep/internal/service"
	"github.com/leandrodeoliveiranovais/api-cep/internal/validation"
)

// mockViaCEPClient é um double de teste para o client ViaCEP.
type mockViaCEPClient struct {
	fetchFunc func(ctx context.Context, cep string) (domain.Address, error)
}

func (m *mockViaCEPClient) FetchByCEP(ctx context.Context, cep string) (domain.Address, error) {
	return m.fetchFunc(ctx, cep)
}

// TestCEPService_GetAddress_Success valida o fluxo feliz do serviço.
func TestCEPService_GetAddress_Success(t *testing.T) {
	// Arrange: mock retorna endereço válido.
	mockClient := &mockViaCEPClient{
		fetchFunc: func(_ context.Context, cep string) (domain.Address, error) {
			if cep != "01001000" {
				t.Errorf("CEP enviado ao client: esperava %q, obteve %q", "01001000", cep)
			}
			return domain.Address{
				CEP:        "01001-000",
				Logradouro: "Praça da Sé",
				Localidade: "São Paulo",
				UF:         "SP",
			}, nil
		},
	}

	svc := service.NewCEPService(mockClient)

	// Act: busca endereço.
	addr, err := svc.GetAddress(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("esperava sucesso, obteve erro: %v", err)
	}

	// Assert: retorno correto.
	if addr.CEP != "01001-000" {
		t.Errorf("CEP: esperava %q, obteve %q", "01001-000", addr.CEP)
	}
}

// TestCEPService_GetAddress_InvalidCEP valida que CEP inválido é rejeitado
// antes de chamar o client externo.
func TestCEPService_GetAddress_InvalidCEP(t *testing.T) {
	called := false
	mockClient := &mockViaCEPClient{
		fetchFunc: func(_ context.Context, _ string) (domain.Address, error) {
			called = true
			return domain.Address{}, nil
		},
	}

	svc := service.NewCEPService(mockClient)

	_, err := svc.GetAddress(context.Background(), "123")
	if err == nil {
		t.Fatal("esperava erro para CEP inválido")
	}

	if !errors.Is(err, validation.ErrInvalidCEP) {
		t.Errorf("esperava ErrInvalidCEP, obteve: %v", err)
	}

	if called {
		t.Error("client externo não deveria ser chamado para CEP inválido")
	}
}

// TestCEPService_GetAddress_NotFound propaga erro de CEP inexistente.
func TestCEPService_GetAddress_NotFound(t *testing.T) {
	mockClient := &mockViaCEPClient{
		fetchFunc: func(_ context.Context, _ string) (domain.Address, error) {
			return domain.Address{}, domain.ErrNotFound
		},
	}

	svc := service.NewCEPService(mockClient)

	_, err := svc.GetAddress(context.Background(), "99999999")
	if err == nil {
		t.Fatal("esperava erro para CEP inexistente")
	}

	if !domain.IsNotFoundError(err) {
		t.Errorf("esperava ErrNotFound, obteve: %v", err)
	}
}

// TestCEPService_GetAddress_UpstreamFailure propaga erro do client.
func TestCEPService_GetAddress_UpstreamFailure(t *testing.T) {
	upstreamErr := errors.New("falha de rede")
	mockClient := &mockViaCEPClient{
		fetchFunc: func(_ context.Context, _ string) (domain.Address, error) {
			return domain.Address{}, upstreamErr
		},
	}

	svc := service.NewCEPService(mockClient)

	_, err := svc.GetAddress(context.Background(), "01001000")
	if !errors.Is(err, upstreamErr) {
		t.Errorf("esperava erro upstream, obteve: %v", err)
	}
}
