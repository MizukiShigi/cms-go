FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o cms-api ./cmd/api/main.go

# 開発用イメージ 
FROM golang:1.24-alpine AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY --from=builder /app .
# 開発モードでの実行コマンド
CMD ["air", "-c", ".air.toml"]

# 本番用イメージ
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# ビルドイメージからバイナリをコピー
COPY --from=builder /app/cms-api .
# 本番実行コマンド
CMD ["./cms-api"]