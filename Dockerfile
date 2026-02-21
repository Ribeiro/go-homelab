# --- Estágio 1: Build ---
# Usamos a imagem do SDK do Go no Wolfi para compilar
FROM cgr.dev/chainguard/go:latest AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro (otimiza cache de camadas)
COPY go.mod ./
# Se você tiver um go.sum, descomente a linha abaixo:
# COPY go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário:
# -ldflags="-w -s" remove tabelas de símbolos (binário menor)
# CGO_ENABLED=0 garante que o binário seja estático (roda em qualquer lugar)
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app-homelab main.go

# --- Estágio 2: Runtime ---
# Usamos a imagem base estática do Wolfi (mínima e segura)
FROM cgr.dev/chainguard/static:latest

# Define um usuário não-privilegiado (padrão nas imagens Chainguard/Wolfi)
USER nonroot

WORKDIR /home/nonroot

# Copia apenas o binário compilado do estágio anterior
COPY --from=builder /app/app-homelab .

# Metadata e porta (se seu app abrir uma porta no futuro)
LABEL org.opencontainers.image.source="https://github.com/seu-user/go-homelab"
LABEL org.opencontainers.image.description="Go App seguro rodando no MacBook 2011"

# Executa o binário
ENTRYPOINT ["./app-homelab"]
