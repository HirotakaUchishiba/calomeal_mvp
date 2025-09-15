# CaloMeal MVP - Production Deployment Guide

## 概要

このガイドでは、CaloMeal MVPアプリケーションを本番環境にデプロイする手順を説明します。

## 前提条件

### システム要件
- Docker 20.10+
- Docker Compose 2.0+
- 最低4GB RAM
- 最低20GB ディスク容量
- Linux/macOS/Windows

### 必要なツール
- `curl`
- `openssl`
- `nc` (netcat)
- `psql` (PostgreSQL client) - オプション

## デプロイ手順

### 1. 初期セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/HirotakaUchishiba/calomeal_mvp.git
cd calomeal_mvp

# 本番環境セットアップスクリプトを実行
./scripts/setup-production.sh
```

このスクリプトは以下を実行します：
- 本番環境設定ファイルの生成
- SSL証明書の生成
- 必要なディレクトリの作成
- モニタリングスクリプトの設定
- バックアップスクリプトの設定

### 2. 環境設定のカスタマイズ

`config/production.env` ファイルを編集して、本番環境に合わせて設定を調整します：

```bash
# 重要な設定項目
DB_PASSWORD=your_secure_database_password
JWT_SECRET=your_jwt_secret_key
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### 3. SSL証明書の設定

本番環境では、自己署名証明書ではなく、信頼できるCAから発行された証明書を使用してください：

```bash
# 証明書を配置
cp your-cert.pem nginx/ssl/cert.pem
cp your-key.pem nginx/ssl/key.pem

# 適切な権限を設定
chmod 644 nginx/ssl/cert.pem
chmod 600 nginx/ssl/key.pem
```

### 4. デプロイの実行

```bash
# 本番環境にデプロイ
./scripts/deploy-production.sh
```

このスクリプトは以下を実行します：
- 既存サービスの停止
- バックアップの作成
- Dockerイメージのビルド
- サービスの起動
- データベースマイグレーション
- ヘルスチェック
- モニタリングの設定

### 5. デプロイの確認

```bash
# 本番環境テストを実行
./scripts/test-production.sh

# または個別のテスト
./scripts/test-production.sh docker
./scripts/test-production.sh network
./scripts/test-production.sh database
./scripts/test-production.sh http
```

## サービス構成

### デプロイされるサービス

1. **PostgreSQL Database** (ポート: 5432)
   - メインデータベース
   - 永続化ボリューム使用

2. **Foods gRPC Service** (ポート: 50051)
   - 食品データ管理サービス

3. **Logs gRPC Service** (ポート: 50052)
   - ログ記録サービス

4. **Analytics gRPC Service** (ポート: 50053)
   - 分析・レポートサービス

5. **Backend BFF Service** (ポート: 8080)
   - GraphQL API
   - 認証・認可

6. **Frontend React App** (ポート: 80)
   - ユーザーインターフェース

7. **Nginx Reverse Proxy** (ポート: 443)
   - SSL終端
   - ロードバランシング
   - セキュリティヘッダー

## アクセス情報

### サービスURL
- **フロントエンド**: https://yourdomain.com
- **GraphQL API**: https://yourdomain.com/api/
- **ヘルスチェック**: https://yourdomain.com/health
- **GraphQL Playground**: https://yourdomain.com/playground (開発時のみ)

### 管理コマンド

```bash
# サービスの状態確認
docker-compose -f docker-compose.prod.yml ps

# ログの確認
docker-compose -f docker-compose.prod.yml logs -f

# 特定サービスのログ
docker-compose -f docker-compose.prod.yml logs -f backend

# サービスの再起動
docker-compose -f docker-compose.prod.yml restart backend

# サービスの停止
docker-compose -f docker-compose.prod.yml down
```

## モニタリング

### ヘルスチェック

```bash
# 包括的なヘルスチェック
./scripts/health-check.sh check

# 継続的なモニタリング
./scripts/health-check.sh monitor

# サービスの準備待ち
./scripts/health-check.sh wait
```

### ログ監視

```bash
# モニタリングスクリプトの実行
./monitoring/monitor.sh

# ログファイルの確認
tail -f /var/log/calomeal/monitor.log
```

### パフォーマンス監視

```bash
# リソース使用量の確認
docker stats

# ディスク使用量の確認
df -h

# メモリ使用量の確認
free -h
```

## バックアップ

### 自動バックアップ

```bash
# バックアップスクリプトの実行
./backups/backup.sh

# バックアップのリスト
ls -la backups/
```

### 手動バックアップ

```bash
# データベースのバックアップ
docker exec calomeal-db-prod pg_dump -U postgres calomeal > backup_$(date +%Y%m%d_%H%M%S).sql

# 設定ファイルのバックアップ
tar -czf config_backup_$(date +%Y%m%d_%H%M%S).tar.gz config/
```

## セキュリティ

### ファイアウォール設定

```bash
# UFWの設定 (Ubuntu/Debian)
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### セキュリティヘッダー

Nginx設定により以下のセキュリティヘッダーが自動的に設定されます：
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security`
- `Content-Security-Policy`

### レート制限

- APIエンドポイント: 10リクエスト/秒
- 認証エンドポイント: 5リクエスト/分

## トラブルシューティング

### よくある問題

1. **サービスが起動しない**
   ```bash
   # ログを確認
   docker-compose -f docker-compose.prod.yml logs
   
   # 環境変数を確認
   cat config/production.env
   ```

2. **データベース接続エラー**
   ```bash
   # データベースの状態確認
   docker exec calomeal-db-prod pg_isready -U postgres -d calomeal
   
   # データベースログの確認
   docker logs calomeal-db-prod
   ```

3. **SSL証明書エラー**
   ```bash
   # 証明書の確認
   openssl x509 -in nginx/ssl/cert.pem -text -noout
   
   # 証明書の権限確認
   ls -la nginx/ssl/
   ```

### ロールバック

```bash
# 前のバージョンにロールバック
./scripts/deploy-production.sh rollback
```

### 緊急時の対応

```bash
# 全サービスの停止
docker-compose -f docker-compose.prod.yml down

# データベースのみ起動
docker-compose -f docker-compose.prod.yml up -d db

# バックアップからの復元
docker exec -i calomeal-db-prod psql -U postgres -d calomeal < backup_file.sql
```

## メンテナンス

### 定期メンテナンス

1. **ログローテーション**
   ```bash
   # ログローテーション設定
   sudo cp monitoring/logrotate.conf /etc/logrotate.d/calomeal
   ```

2. **セキュリティアップデート**
   ```bash
   # システムアップデート
   sudo apt update && sudo apt upgrade
   
   # Dockerイメージの更新
   docker-compose -f docker-compose.prod.yml pull
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **バックアップの確認**
   ```bash
   # バックアップのテスト
   ./backups/backup.sh
   
   # 古いバックアップの削除
   find backups/ -name "*.sql.gz" -mtime +30 -delete
   ```

## スケーリング

### 水平スケーリング

```bash
# 特定サービスのスケール
docker-compose -f docker-compose.prod.yml up -d --scale backend=3
```

### 垂直スケーリング

`docker-compose.prod.yml` でリソース制限を調整：

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 1G
          cpus: '0.5'
```

## サポート

### ログの収集

```bash
# デバッグ情報の収集
./scripts/collect-debug-info.sh
```

### パフォーマンス分析

```bash
# パフォーマンステスト
./scripts/test-production.sh performance
```

## 更新履歴

- 2025-01-14: 初版作成
- 本番環境デプロイ準備完了

---

**注意**: 本番環境での運用前に、必ずステージング環境でテストを実施してください。
