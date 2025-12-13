# Etapa 1: build
FROM golang:1.25.3-alpine AS builder

# Instala dependências básicas
RUN apk add --no-cache git

# Define diretório de trabalho
WORKDIR /app

# Copia arquivos do módulo e baixa dependências
COPY go.mod ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN go build -o server ./cmd/server

# Etapa 2: imagem final
FROM alpine:3.20

# Instala certificados SSL (necessários p/ conexões HTTPS)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copia o binário gerado
COPY --from=builder /app/server .

COPY private.pem .
COPY public.pem .

# Define variáveis de ambiente padrão
ENV PORT=8080
EXPOSE 8080

# Comando de execução
CMD ["./server"]

