# Docker Deployment Guide

Go Movie ExplorerアプリケーションのDocker化とデプロイメント手順書です。

## 📋 前提条件

- Docker Desktop または Docker Engine がインストールされていること
- Google Cloud CLI がインストールされていること（Cloud Run使用時）
- 必要な環境変数が設定されていること

## 🚀 ローカルでのDocker実行

### バックエンドのビルドと実行

```bash
# バックエンドディレクトリに移動
cd backend

# Dockerイメージをビルド
docker build -t go-movie-explorer-backend .

# コンテナを実行（環境変数を設定）
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e TMDB_API_KEY=your_tmdb_api_key \
  -e GO_ENV=production \
  go-movie-explorer-backend
```

### フロントエンドのビルドと実行

```bash
# フロントエンドディレクトリに移動
cd frontend

# Dockerイメージをビルド
docker build -t go-movie-explorer-frontend .

# コンテナを実行
docker run -p 80:80 go-movie-explorer-frontend
```

### Docker Composeでの実行（推奨）

プロジェクトルートに以下の`docker-compose.yml`を作成：

```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - TMDB_API_KEY=${TMDB_API_KEY}
      - GO_ENV=production
    networks:
      - app-network

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

実行コマンド：
```bash
# 環境変数を設定
export TMDB_API_KEY=your_tmdb_api_key

# アプリケーションを起動
docker-compose up --build
```

## ☁️ Google Cloud Runでのデプロイメント

### 1. Google Cloud Projectの設定

```bash
# プロジェクトを設定
gcloud config set project YOUR_PROJECT_ID

# Artifact Registryを有効化
gcloud services enable artifactregistry.googleapis.com

# レポジトリを作成
gcloud artifacts repositories create go-movie-explorer \
  --repository-format=docker \
  --location=asia-northeast1
```

### 2. バックエンドのデプロイ

```bash
# バックエンドディレクトリに移動
cd backend

# イメージをビルドしてプッシュ
gcloud builds submit --tag asia-northeast1-docker.pkg.dev/YOUR_PROJECT_ID/go-movie-explorer/backend

# Cloud Runにデプロイ
gcloud run deploy go-movie-explorer-backend \
  --image asia-northeast1-docker.pkg.dev/YOUR_PROJECT_ID/go-movie-explorer/backend \
  --platform managed \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --set-env-vars PORT=8080,GO_ENV=production,TMDB_API_KEY=your_tmdb_api_key \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10
```

### 3. フロントエンドのデプロイ

```bash
# フロントエンドディレクトリに移動
cd frontend

# イメージをビルドしてプッシュ
gcloud builds submit --tag asia-northeast1-docker.pkg.dev/YOUR_PROJECT_ID/go-movie-explorer/frontend

# Cloud Runにデプロイ
gcloud run deploy go-movie-explorer-frontend \
  --image asia-northeast1-docker.pkg.dev/YOUR_PROJECT_ID/go-movie-explorer/frontend \
  --platform managed \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --port 80 \
  --memory 256Mi \
  --cpu 1 \
  --max-instances 5
```

## 🔧 環境変数の設定

### バックエンド必須環境変数
- `PORT`: サーバーのポート番号（デフォルト: 8080）
- `TMDB_API_KEY`: TMDB APIキー
- `GO_ENV`: 環境（production推奨）
- `FRONTEND_URL`: フロントエンドのURL（CORS設定用）

### セキュリティ設定
```bash
# Secret Managerを使用（推奨）
gcloud secrets create tmdb-api-key --data-file=-

# Cloud Runで環境変数として設定
gcloud run services update go-movie-explorer-backend \
  --update-secrets TMDB_API_KEY=tmdb-api-key:latest
```

## 📊 モニタリングとログ

### ヘルスチェック
- バックエンド: `http://your-backend-url/healthz`
- フロントエンド: `http://your-frontend-url/health`

### ログの確認
```bash
# Cloud Runのログを確認
gcloud run logs tail go-movie-explorer-backend --region asia-northeast1
gcloud run logs tail go-movie-explorer-frontend --region asia-northeast1
```

## 🐛 トラブルシューティング

### よくある問題

1. **TMDB_API_KEYが設定されていない**
   - 環境変数が正しく設定されているか確認
   - Secret Managerの設定を確認

2. **CORS エラー**
   - `FRONTEND_URL`環境変数が正しく設定されているか確認
   - バックエンドのセキュリティ設定を確認

3. **ビルドエラー**
   - `.dockerignore`ファイルの設定を確認
   - 依存関係の問題がないか確認

### パフォーマンス最適化

- **メモリ使用量の調整**: Cloud Runのメモリ設定を最適化
- **CPU使用量の調整**: 必要に応じてCPU設定を調整
- **自動スケーリング**: トラフィックに応じてインスタンス数を調整

## 📝 備考

- 本番環境では必ずHTTPSを使用してください
- セキュリティヘッダーの設定を確認してください
- 定期的にコンテナイメージを更新してください
- ログとメトリクスを監視してください