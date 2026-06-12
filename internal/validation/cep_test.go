package validation_test

import (
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/validation"
)

// TestValidateCEP_Valid verifica que um CEP com exatamente 8 dígitos numéricos é aceito.
func TestValidateCEP_Valid(t *testing.T) {
	// Arrange: CEP no formato esperado pela API (sem hífen).
	cep := "01001000"

	// Act: executa a validação.
	err := validation.ValidateCEP(cep)

	// Assert: não deve retornar erro.
	if err != nil {
		t.Fatalf("esperava CEP válido %q, obteve erro: %v", cep, err)
	}
}

// TestValidateCEP_InvalidLength verifica rejeição de CEPs com tamanho diferente de 8.
func TestValidateCEP_InvalidLength(t *testing.T) {
	cases := []struct {
		name string
		cep  string
	}{
		{name: "menos de 8 dígitos", cep: "0100100"},
		{name: "mais de 8 dígitos", cep: "010010000"},
		{name: "vazio", cep: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validation.ValidateCEP(tc.cep)
			if err == nil {
				t.Fatalf("esperava erro para CEP %q", tc.cep)
			}
		})
	}
}

// TestValidateCEP_InvalidCharacters verifica rejeição de CEPs com caracteres não numéricos.
func TestValidateCEP_InvalidCharacters(t *testing.T) {
	cases := []struct {
		name string
		cep  string
	}{
		{name: "alfanumérico", cep: "95010A10"},
		{name: "com hífen", cep: "01001-000"},
		{name: "com espaço", cep: "01001 000"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validation.ValidateCEP(tc.cep)
			if err == nil {
				t.Fatalf("esperava erro para CEP %q", tc.cep)
			}
		})
	}
}
