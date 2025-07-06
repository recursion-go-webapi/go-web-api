# Docker Compose 使用ガイド

Go Movie ExplorerをDocker Composeで起動する方法です。
基本的にはrootで`docker compose up -d`とコマンド打つとDockerが起動します。

## 🚀 クイックスタート

### 1. 環境変数の設定

```bash
# .envファイルを作成
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
# .envファイルを編集してTMDB_API_KEYを設定
```

### 2. アプリケーションの起動

```bash
# アプリケーションを起動
docker compose up --build

# バックグラウンドで起動
docker compose up -d --build

# ログを確認
docker compose logs -f
```

## 📍 アクセス先

- **フロントエンド**: http://localhost:3003
- **フロントエンドNginxヘルスチェック**: http://localhost:3003/health
- **バックエンドAPI**: http://localhost:8080/api
- **ヘルスチェック**: http://localhost:8080/healthz

## 🛠️ 便利なコマンド

### サービス管理

```bash
# 特定のサービスのみ起動
docker compose up backend
docker compose up frontend

# サービス停止
docker compose stop

# サービス削除（ボリューム保持）
docker compose down

# サービス削除（ボリューム含む）
docker compose down -v

# 再ビルド
docker compose build --no-cache
```

### ログとデバッグ

```bash
# 全サービスのログ
docker compose logs

# 特定のサービスのログ
docker compose logs backend
docker compose logs frontend

# リアルタイムログ
docker compose logs -f backend

# コンテナに接続
docker compose exec backend sh
docker compose exec frontend sh
```

## ⚙️ 設定詳細

### ポート設定

| サービス | 内部ポート | 外部ポート | 説明                  |
| -------- | ---------- | ---------- | --------------------- |
| backend  | 8080       | 8080       | Go APIサーバー        |
| frontend | 80         | 3003       | Nginx静的ファイル配信 |

### ヘルスチェック

- **バックエンド**: `/healthz` エンドポイントで TMDB API 接続確認
- **フロントエンド**: `/health` エンドポイントで Nginx 稼働確認

## 🐛 トラブルシューティング

### よくある問題

1. **TMDB_API_KEYが設定されていない**
   ```bash
   # 環境変数を確認
   echo $TMDB_API_KEY
   
   # 設定
   export TMDB_API_KEY=your_api_key
   ```

2. **ポートが既に使用されている**
   ```bash
   # ポート使用状況を確認
   lsof -i :8080
   lsof -i :3003
   
   # docker compose.ymlでポートを変更
   # "3004:80" など

   # または `kill -9 {PID}` で使用中ポートを削除
   ```

3. **ビルドエラー**
   ```bash
   # キャッシュを無効化して再ビルド
   docker compose build --no-cache
   
   # イメージを削除して再作成
   docker compose down --rmi all
   docker compose up --build
   ```

4. **CORS エラー**
   ```bash
   # FRONTEND_URLが正しく設定されているか確認
   docker compose logs backend | grep CORS
   ```

### パフォーマンス最適化

```bash
# 不要なイメージ・コンテナを削除
docker system prune -a

# ビルドキャッシュを活用
docker compose build

# リソース使用量を確認
docker compose top
docker stats
```

### 環境設定

**開発環境**:
```bash
# .envファイルを作成・編集
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
# .envファイル内のTMDB_API_KEYを設定
docker compose up --build
```

**本番環境**:
```bash
# 環境変数を直接設定
export TMDB_API_KEY=your_production_api_key
export FRONTEND_URL=https://your-domain.com
export GO_ENV=production
docker compose up -d --build
```
