// Package handler expõe os endpoints HTTP da API.
package handler

import "net/http"

// healthResponse representa o payload JSON retornado pelo endpoint de health check.
type healthResponse struct {
	Status string `json:"status"` // Indica se o processo da API está ativo ("ok")
}

// HealthCheck responde se o servidor HTTP está ativo e pronto para receber requisições.
// Não consulta dependências externas (ViaCEP) — adequado para liveness probe.
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}
