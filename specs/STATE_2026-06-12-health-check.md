# STATE — Health Check Endpoint

**Data:** 2026-06-12  
**Status:** 🟢 Implementado  
**Escopo:** Endpoint de health check para monitoramento e liveness probe

---

## 1. Objetivo

Expor um endpoint leve que confirme que o processo HTTP da API está ativo, sem depender do ViaCEP ou de outras integrações externas.

---

## 2. Contrato

### Request

```
GET /health
```

### Response — Sucesso (200)

```json
{
  "status": "ok"
}
```

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `status` | string | Sempre `"ok"` quando o servidor está respondendo |

---

## 3. Decisões tomadas

| # | Decisão | Alternativa considerada | Motivo |
|---|---------|------------------------|--------|
| 1 | Path `/health` | `/healthz`, `/ready`, `/live` | Convenção mais comum em APIs REST e load balancers |
| 2 | Liveness only (sem checar ViaCEP) | Readiness com ping ao ViaCEP | Health check deve ser rápido e estável; falha do ViaCEP não significa que a API está down |
| 3 | Resposta JSON `{"status":"ok"}` | Resposta vazia ou plain text `OK` | Consistência com o restante da API (JSON) e extensível no futuro |
| 4 | Handler em `internal/handler/health.go` | Endpoint inline no `main.go` | Mantém roteamento centralizado em `NewRouter` |
| 5 | Sem camada service | Service dedicado `HealthService` | Operação trivial; adicionar camada seria over-engineering |
| 6 | Reutilizar `writeJSON` existente | Handler com `json.Marshal` próprio | DRY — mesmo content-type e serialização dos demais endpoints |
| 7 | Registro em `NewRouter` | Router separado para health | Um único ponto de registro de rotas facilita testes e manutenção |

---

## 4. Arquitetura

```mermaid
flowchart LR
    LB["Load Balancer / K8s"] -->|GET /health| Router
    Router --> HealthCheck["HealthCheck handler"]
    HealthCheck -->|200 JSON| LB

    Client -->|GET /cep/{cep}| Router
    Router --> CEPHandler["GetCEP handler"]
    CEPHandler --> Service --> ViaCEP
```

O health check **não** passa por service, client ou ViaCEP.

---

## 5. Arquivos alterados / criados

| Arquivo | Ação |
|---------|------|
| `internal/handler/health.go` | Criado — handler e struct de resposta |
| `internal/handler/health_test.go` | Criado — teste unitário |
| `internal/handler/cep.go` | Alterado — registro da rota `GET /health` em `NewRouter` |

---

## 6. Testes

| Arquivo | Caso |
|---------|------|
| `internal/handler/health_test.go` | `GET /health` retorna 200, `Content-Type: application/json`, body `{"status":"ok"}` |

Executar:

```bash
go test ./internal/handler/... -v -run Health
```

---

## 7. Uso operacional

### Verificação manual

```bash
curl -s http://localhost:8080/health
# {"status":"ok"}
```

### Kubernetes liveness probe (exemplo)

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 10
```

---

## 8. Evolução futura (fora do escopo atual)

- **`GET /ready`** — readiness probe consultando disponibilidade do ViaCEP
- **Campo `version`** — expor versão da build no JSON de health
- **Campo `uptime`** — tempo desde o start do processo

---

## 9. Referências

- Implementação inicial da API: [`STATE_2026-06-12.md`](STATE_2026-06-12.md) (seção API CEP)
- Changelog deste incremento: [`../changelogs/CHANGE_LOG_2026-06-12.md`](../changelogs/CHANGE_LOG_2026-06-12.md)
