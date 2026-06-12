// Package client encapsula a comunicação HTTP com o webservice ViaCEP.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
)

// ViaCEPClient consulta endereços no webservice ViaCEP via HTTP.
type ViaCEPClient struct {
	httpClient *http.Client // Cliente HTTP reutilizável (permite timeout e injeção em testes)
	baseURL    string       // URL base do ViaCEP (ex: "https://viacep.com.br")
}

// NewViaCEPClient cria um client configurado com o HTTP client e base URL informados.
func NewViaCEPClient(httpClient *http.Client, baseURL string) *ViaCEPClient {
	// Remove barra final da base URL para evitar paths duplicados.
	return &ViaCEPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

// FetchByCEP consulta o ViaCEP e retorna o endereço mapeado em domain.Address.
func (c *ViaCEPClient) FetchByCEP(ctx context.Context, cep string) (domain.Address, error) {
	// Monta a URL no formato oficial: /ws/{cep}/json/
	url := fmt.Sprintf("%s/ws/%s/json/", c.baseURL, cep)

	// Cria requisição HTTP com contexto para suportar cancelamento e timeout.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Address{}, fmt.Errorf("criar requisição: %w", err)
	}

	// Define o header Accept para indicar que esperamos JSON.
	req.Header.Set("Accept", "application/json")

	// Executa a requisição HTTP contra o ViaCEP.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Address{}, fmt.Errorf("executar requisição: %w", err)
	}
	// Garante que o body será fechado ao final da função.
	defer resp.Body.Close()

	// ViaCEP retorna 400 para CEPs com formato inválido.
	if resp.StatusCode == http.StatusBadRequest {
		return domain.Address{}, fmt.Errorf("viacep retornou status %d", resp.StatusCode)
	}

	// Qualquer outro status fora do 2xx é tratado como falha no upstream.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.Address{}, fmt.Errorf("viacep retornou status %d", resp.StatusCode)
	}

	// Lê o corpo da resposta para deserialização JSON.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Address{}, fmt.Errorf("ler corpo da resposta: %w", err)
	}

	// Deserializa o JSON do ViaCEP para a struct de domínio.
	var addr domain.Address
	if err := json.Unmarshal(body, &addr); err != nil {
		return domain.Address{}, fmt.Errorf("deserializar resposta: %w", err)
	}

	// ViaCEP sinaliza CEP inexistente com {"erro": "true"} e status 200.
	if addr.IsNotFound() {
		return domain.Address{}, domain.ErrNotFound
	}

	// Retorna o endereço mapeado com sucesso.
	return addr, nil
}
