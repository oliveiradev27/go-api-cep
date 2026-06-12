// Package main é o ponto de entrada da API de consulta de CEP.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leandrodeoliveiranovais/api-cep/internal/client"
	"github.com/leandrodeoliveiranovais/api-cep/internal/handler"
	"github.com/leandrodeoliveiranovais/api-cep/internal/service"
)

// main inicializa as dependências e inicia o servidor HTTP.
func main() {
	// Lê a porta do servidor a partir da variável de ambiente (padrão: 8080).
	port := envOrDefault("PORT", "8080")

	// Lê a URL base do ViaCEP (padrão: produção; use mock local via VIACEP_BASE_URL).
	viacepBaseURL := envOrDefault("VIACEP_BASE_URL", "https://viacep.com.br")

	// Cria o client HTTP com timeout para evitar requisições penduradas.
	httpClient := &http.Client{Timeout: 10 * time.Second}

	// Monta a cadeia de dependências: client → service → handler.
	viacepClient := client.NewViaCEPClient(httpClient, viacepBaseURL)
	cepService := service.NewCEPService(viacepClient)
	cepHandler := handler.NewCEPHandler(cepService)

	// Configura o roteador com as rotas da API.
	router := handler.NewRouter(cepHandler)

	// Inicia o servidor HTTP na porta configurada.
	log.Printf("servidor iniciado em :%s (viacep: %s)", port, viacepBaseURL)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
	}
}

// envOrDefault retorna o valor da variável de ambiente ou o padrão informado.
func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
