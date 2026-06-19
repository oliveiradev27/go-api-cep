// Package handler expõe os endpoints HTTP da API.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/leandrodeoliveiranovais/api-cep/internal/domain"
	"github.com/leandrodeoliveiranovais/api-cep/internal/validation"
)

// CEPService define o contrato do serviço de CEP (permite injeção de mocks em testes).
type CEPService interface {
	GetAddress(ctx context.Context, cep string) (domain.Address, error)
}

// CEPHandler trata requisições HTTP relacionadas a consulta de CEP.
type CEPHandler struct {
	service CEPService // Serviço de negócio injetado
}

// NewCEPHandler cria o handler com o serviço informado.
func NewCEPHandler(service CEPService) *CEPHandler {
	return &CEPHandler{service: service}
}

// NewRouter registra as rotas da API e retorna um http.Handler pronto para uso.
func NewRouter(h *CEPHandler) http.Handler {
	mux := http.NewServeMux()
	// Health check: confirma que o processo HTTP está ativo (liveness).
	mux.HandleFunc("GET /health", HealthCheck)
	// Rota principal: GET /cep/{cep}
	mux.HandleFunc("GET /cep/{cep}", h.GetCEP)
	return mux
}

// GetCEP é o handler do endpoint GET /cep/{cep}.
func (h *CEPHandler) GetCEP(w http.ResponseWriter, r *http.Request) {
	// Extrai o CEP do path parameter definido na rota.
	cep := r.PathValue("cep")

	// Valida o formato do CEP antes de chamar o serviço (resposta rápida para entradas inválidas).
	if err := validation.ValidateCEP(cep); err != nil {
		writeError(w, http.StatusBadRequest, "CEP inválido: deve conter exatamente 8 dígitos numéricos")
		return
	}

	// Delega a busca ao serviço de negócio.
	addr, err := h.service.GetAddress(r.Context(), cep)
	if err != nil {
		// Mapeia erros de domínio para status HTTP semânticos.
		switch {
		case errors.Is(err, validation.ErrInvalidCEP):
			writeError(w, http.StatusBadRequest, err.Error())
		case domain.IsNotFoundError(err):
			writeError(w, http.StatusNotFound, "CEP não encontrado")
		default:
			writeError(w, http.StatusBadGateway, "falha ao consultar serviço de CEP")
		}
		return
	}

	// Serializa a struct de domínio para JSON e envia como resposta.
	writeJSON(w, http.StatusOK, addr)
}

// writeJSON serializa qualquer valor para JSON e escreve na resposta HTTP.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	// Define o content-type como JSON.
	w.Header().Set("Content-Type", "application/json")
	// Define o status HTTP da resposta.
	w.WriteHeader(status)
	// Serializa o payload e escreve no body da resposta.
	_ = json.NewEncoder(w).Encode(payload)
}

// writeError escreve uma resposta de erro padronizada em JSON.
func writeError(w http.ResponseWriter, status int, message string) {
	// Monta o objeto de erro com a mensagem informada.
	payload := map[string]string{"error": strings.TrimSpace(message)}
	writeJSON(w, status, payload)
}
