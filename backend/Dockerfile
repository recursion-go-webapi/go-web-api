# マルチステージビルド用のGoビルド環境
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /app

# go.modとgo.sumをコピーしてdependenciesをダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# 静的リンクされたバイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 本番用の軽量イメージ
FROM alpine:latest

# セキュリティアップデートとCA証明書をインストール
RUN apk --no-cache add ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /root/

# ビルド済みバイナリをコピー
COPY --from=builder /app/main .

# logsディレクトリを作成
RUN mkdir -p logs

# ポート8080を公開
EXPOSE 8080

# 本番環境用の環境変数を設定
ENV GO_ENV=production

# アプリケーションを実行
CMD ["./main"]
