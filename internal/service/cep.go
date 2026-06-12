// Package service contém a lógica de negócio da aplicação.
package service

import (
	"context"

	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
	"github.com/leandrodeoliveiranovais/api-cep/internal/validation"
)

// ViaCEPClient define o contrato do client externo (permite injeção de mocks em testes).
type ViaCEPClient interface {
	FetchByCEP(ctx context.Context, cep string) (domain.Address, error)
}

// CEPService orquestra validação e consulta de endereços por CEP.
type CEPService struct {
	client ViaCEPClient // Dependência injetada do client ViaCEP
}

// NewCEPService cria o serviço com o client ViaCEP informado.
func NewCEPService(client ViaCEPClient) *CEPService {
	return &CEPService{client: client}
}

// GetAddress valida o CEP e delega a consulta ao client ViaCEP.
func (s *CEPService) GetAddress(ctx context.Context, cep string) (domain.Address, error) {
	// Valida o formato do CEP antes de chamar o serviço externo (fail fast).
	if err := validation.ValidateCEP(cep); err != nil {
		return domain.Address{}, err
	}

	// Delega a busca ao client HTTP do ViaCEP.
	addr, err := s.client.FetchByCEP(ctx, cep)
	if err != nil {
		return domain.Address{}, err
	}

	// Retorna o endereço obtido.
	return addr, nil
}
