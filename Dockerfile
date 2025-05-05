FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o cms-api ./cmd/api/main.go

# 実行用の小さなイメージを作成
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドイメージからバイナリをコピー
COPY --from=builder /app/cms-api .

# アプリケーションの実行
CMD ["./cms-api"]