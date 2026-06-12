// Package validation contém regras de validação de entrada da API.
package validation

import (
	"errors"
	"unicode"
)

// ErrInvalidCEP é retornado quando o CEP não possui exatamente 8 dígitos numéricos.
var ErrInvalidCEP = errors.New("cep inválido: deve conter exatamente 8 dígitos numéricos")

// ValidateCEP verifica se o CEP informado está no formato aceito pela API.
// Regra alinhada ao ViaCEP: exatamente 8 caracteres, todos numéricos, sem hífen ou espaço.
func ValidateCEP(cep string) error {
	// Rejeita CEPs com tamanho diferente de 8 caracteres.
	if len(cep) != 8 {
		return ErrInvalidCEP
	}

	// Percorre cada caractere garantindo que seja um dígito (0-9).
	for _, ch := range cep {
		if !unicode.IsDigit(ch) {
			return ErrInvalidCEP
		}
	}

	// CEP válido: 8 dígitos numéricos consecutivos.
	return nil
}
