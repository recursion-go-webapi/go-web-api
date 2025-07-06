# マルチステージビルド用のNode.jsビルド環境
FROM node:20-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# package.jsonとpackage-lock.jsonをコピー
COPY package*.json ./

# 依存関係をインストール（開発依存関係も含める）
RUN npm ci

# ソースコードをコピー
COPY . .

# Viteアプリケーションをビルド
RUN npm run build

# 本番用のNginx環境
FROM nginx:alpine

# カスタムnginx設定をコピー
COPY nginx.conf /etc/nginx/nginx.conf

# ビルド済み静的ファイルをNginxのドキュメントルートにコピー
COPY --from=builder /app/dist /usr/share/nginx/html

# ポート80を公開
EXPOSE 80

# Nginxを起動
CMD ["nginx", "-g", "daemon off;"]