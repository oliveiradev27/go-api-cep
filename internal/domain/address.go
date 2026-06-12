// Package domain contém os tipos centrais do negócio e erros de domínio.
package domain

import "errors"

// ErrNotFound é retornado quando o CEP tem formato válido mas não existe na base do ViaCEP.
var ErrNotFound = errors.New("cep não encontrado")

// Address representa o endereço retornado pelo webservice ViaCEP.
// Cada campo possui uma tag json que mapeia exatamente o contrato da API externa.
type Address struct {
	CEP         string `json:"cep"`         // CEP formatado com hífen (ex: "01001-000")
	Logradouro  string `json:"logradouro"`  // Nome da rua, avenida ou logradouro
	Complemento string `json:"complemento"` // Informação complementar do endereço
	Unidade     string `json:"unidade"`     // Unidade específica, quando aplicável
	Bairro      string `json:"bairro"`      // Bairro do endereço
	Localidade  string `json:"localidade"`  // Cidade
	UF          string `json:"uf"`          // Sigla do estado (ex: "SP")
	Estado      string `json:"estado"`      // Nome completo do estado
	Regiao      string `json:"regiao"`      // Região geográfica (ex: "Sudeste")
	IBGE        string `json:"ibge"`        // Código IBGE do município
	GIA         string `json:"gia"`         // Código GIA/ICMS (somente SP)
	DDD         string `json:"ddd"`         // Código de área telefônica
	SIAFI       string `json:"siafi"`       // Código SIAFI do município
	Erro        string `json:"erro,omitempty"` // Flag de erro do ViaCEP; omitido na resposta de sucesso
}

// IsNotFound verifica se o ViaCEP sinalizou que o CEP não foi encontrado.
// O webservice retorna {"erro": "true"} para CEPs válidos porém inexistentes.
func (a Address) IsNotFound() bool {
	return a.Erro == "true"
}

// IsNotFoundError verifica se o erro informado é um ErrNotFound.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}
