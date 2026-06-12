// Package main implementa um servidor mock do ViaCEP para testes locais.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// cepPattern valida o formato do CEP na URL (8 dígitos numéricos).
var cepPattern = regexp.MustCompile(`^\d{8}$`)

// main inicia o mock server que simula o comportamento do ViaCEP.
func main() {
	// Porta do mock server (padrão: 8081).
	port := envOrDefault("MOCK_PORT", "8081")

	// Diretório com os arquivos JSON de resposta mockados.
	responsesDir := envOrDefault("MOCK_RESPONSES_DIR", "mocks/responses")

	// Registra o handler que intercepta requisições no formato /ws/{cep}/json/.
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		handleViaCEP(w, r, responsesDir)
	})

	// Inicia o servidor mock.
	addr := ":" + port
	log.Printf("mock ViaCEP iniciado em http://localhost%s", addr)
	log.Printf("diretório de respostas: %s", responsesDir)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("falha ao iniciar mock server: %v", err)
	}
}

// handleViaCEP processa requisições no formato do ViaCEP: /ws/{cep}/json/.
func handleViaCEP(w http.ResponseWriter, r *http.Request, responsesDir string) {
	// Extrai o CEP do path: /ws/01001000/json/ → 01001000
	cep := extractCEP(r.URL.Path)

	// CEP com formato inválido retorna 400 (comportamento do ViaCEP real).
	if !cepPattern.MatchString(cep) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Tenta carregar resposta mockada do arquivo JSON correspondente.
	data, err := os.ReadFile(filepath.Join(responsesDir, cep+".json"))
	if err != nil {
		// CEP sem mock cadastrado: retorna erro de não encontrado.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"erro": "true"}`))
		return
	}

	// Valida que o arquivo é JSON válido antes de servir.
	if !json.Valid(data) {
		http.Error(w, "arquivo de mock inválido", http.StatusInternalServerError)
		return
	}

	// Retorna o JSON mockado com status 200.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// extractCEP extrai o CEP do path no formato /ws/{cep}/json/.
func extractCEP(path string) string {
	// Remove prefixo /ws/ e sufixo /json/ do path.
	path = strings.TrimPrefix(path, "/ws/")
	path = strings.TrimSuffix(path, "/json/")
	path = strings.TrimSuffix(path, "/")
	return path
}

// envOrDefault retorna o valor da variável de ambiente ou o padrão.
func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
