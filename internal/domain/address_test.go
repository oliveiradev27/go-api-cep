package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
)

// TestAddress_JSONSerialization garante que a struct mapeia e serializa
// exatamente o contrato JSON retornado pelo ViaCEP.
func TestAddress_JSONSerialization(t *testing.T) {
	// Arrange: endereço preenchido conforme resposta real do ViaCEP.
	addr := domain.Address{
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
	}

	// Act: serializa para JSON.
	data, err := json.Marshal(addr)
	if err != nil {
		t.Fatalf("falha ao serializar: %v", err)
	}

	// Assert: campos obrigatórios do contrato estão presentes.
	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("falha ao deserializar resultado: %v", err)
	}

	expected := map[string]string{
		"cep":         "01001-000",
		"logradouro":  "Praça da Sé",
		"complemento": "lado ímpar",
		"unidade":     "",
		"bairro":      "Sé",
		"localidade":  "São Paulo",
		"uf":          "SP",
		"estado":      "São Paulo",
		"regiao":      "Sudeste",
		"ibge":        "3550308",
		"gia":         "1004",
		"ddd":         "11",
		"siafi":       "7107",
	}

	for key, want := range expected {
		got, ok := result[key]
		if !ok {
			t.Errorf("campo %q ausente no JSON", key)
			continue
		}
		if got != want {
			t.Errorf("campo %q: esperava %q, obteve %q", key, want, got)
		}
	}
}

// TestAddress_JSONDeserialization garante que o JSON do ViaCEP é corretamente
// deserializado para a struct de domínio.
func TestAddress_JSONDeserialization(t *testing.T) {
	// Arrange: payload JSON idêntico ao retorno do ViaCEP.
	raw := `{
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
	}`

	// Act: deserializa para struct.
	var addr domain.Address
	if err := json.Unmarshal([]byte(raw), &addr); err != nil {
		t.Fatalf("falha ao deserializar: %v", err)
	}

	// Assert: campos mapeados corretamente.
	if addr.CEP != "01001-000" {
		t.Errorf("CEP: esperava %q, obteve %q", "01001-000", addr.CEP)
	}
	if addr.Logradouro != "Praça da Sé" {
		t.Errorf("Logradouro: esperava %q, obteve %q", "Praça da Sé", addr.Logradouro)
	}
	if addr.UF != "SP" {
		t.Errorf("UF: esperava %q, obteve %q", "SP", addr.UF)
	}
}

// TestAddress_ErrorFlag detecta o campo "erro": true retornado pelo ViaCEP
// quando o CEP é válido no formato mas inexistente na base.
func TestAddress_ErrorFlag(t *testing.T) {
	raw := `{"erro": "true"}`

	var addr domain.Address
	if err := json.Unmarshal([]byte(raw), &addr); err != nil {
		t.Fatalf("falha ao deserializar: %v", err)
	}

	if !addr.IsNotFound() {
		t.Error("esperava IsNotFound() == true para resposta com erro do ViaCEP")
	}
}
