# CaloMeal MVP

栄養管理アプリケーション「CaloMeal」のMVP（Minimum Viable Product）です。ユーザーが日々の食事摂取量や身体活動を記録し、健康目標の達成を支援するWebベースのアプリケーションです。

## 🎯 プロジェクト概要

CaloMeal MVPは「Lean MVP」の思想に基づき、最小限のコストと時間で「**ユーザーは栄養を追跡するために食事を記録する**」という中核的なユーザー行動仮説を検証することを目的としています。

### 主要機能

- **ユーザー認証**: AWS Cognitoによる安全な認証システム
- **食事記録**: 食品データベースからのキーワード検索による食事記録
- **運動記録**: 運動の種類と時間、消費カロリーの記録
- **体重記録**: 日々の体重と体脂肪率の記録
- **ダッシュボード**: 日次サマリーとPFCバランスの可視化
- **オンボーディング**: 初回プロフィールと目標設定

## 🏗️ アーキテクチャ

### 技術スタック

- **フロントエンド**: React + TypeScript + Vite
- **バックエンド**: Go + GraphQL (gqlgen)
- **データベース**: PostgreSQL
- **認証**: AWS Cognito
- **インフラ**: AWS (ECS Fargate, RDS, ALB, VPC)
- **コンテナ**: Docker + Docker Compose

### システム構成

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React SPA     │    │   GraphQL BFF   │    │  Microservices  │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│   (Go + gRPC)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │   (Database)    │
                       └─────────────────┘
```

## 🚀 クイックスタート

### 前提条件

- Docker & Docker Compose
- Go 1.25.0+
- Node.js 18+
- AWS CLI (本番環境用)

### ローカル開発環境の起動

1. **リポジトリのクローン**
   ```bash
   git clone https://github.com/HirotakaUchishiba/calomeal_mvp.git
   cd calomeal_mvp
   ```

2. **環境変数の設定**
   ```bash
   cp .env.example .env
   # .envファイルを編集して必要な環境変数を設定
   ```

3. **Docker Composeでサービス起動**
   ```bash
   docker-compose up -d
   ```

4. **サービスアクセス**
   - フロントエンド: http://localhost:3000
   - バックエンド: http://localhost:8080
   - GraphQL Playground: http://localhost:8080/query

### 個別開発

#### フロントエンド開発
```bash
cd frontend
npm install
npm run dev
```

#### バックエンド開発
```bash
cd backend
go mod download
go run cmd/server/main.go
```

## 📁 プロジェクト構造

```
calomeal_mvp/
├── frontend/                 # React フロントエンド
│   ├── src/
│   │   ├── components/      # 再利用可能なコンポーネント
│   │   ├── features/        # 機能別ページ・モーダル
│   │   ├── contexts/        # React Context
│   │   └── graphql/         # GraphQL クエリ
│   └── package.json
├── backend/                  # Go バックエンド
│   ├── cmd/server/          # アプリケーションエントリーポイント
│   ├── internal/
│   │   ├── bff/            # GraphQL BFF層
│   │   └── service/        # ビジネスロジック層
│   ├── database/           # マイグレーション・シード
│   └── schema.graphql      # GraphQL スキーマ
├── terraform/              # インフラストラクチャ
├── tests/                  # E2Eテスト
└── docs/                   # プロジェクトドキュメント
```

## 🧪 テスト

### E2Eテストの実行
```bash
# 全テスト実行
npm run test:e2e

# UI付きでテスト実行
npm run test:e2e:ui

# ヘッド付きブラウザでテスト実行
npm run test:e2e:headed
```

## 📚 ドキュメント

詳細なドキュメントは `docs/` ディレクトリにあります：

- [プロジェクト憲章とスコープ](docs/01_Project_Charter_and_Scope.md)
- [システムアーキテクチャ](docs/02_System_Architecture.md)
- [バックエンド設計](docs/03_Backend_Design.md)
- [フロントエンド設計](docs/04_Frontend_Design.md)
- [インフラストラクチャとDevOps](docs/05_Infrastructure_and_DevOps.md)

## 🔧 開発ガイドライン

### GraphQL スキーマの更新
```bash
cd backend
go run github.com/99designs/gqlgen generate
```

### データベースマイグレーション
```bash
# 新しいマイグレーションファイルを作成
# backend/database/migrations/ にSQLファイルを追加
```

## 🚀 デプロイメント

### 本番環境へのデプロイ
```bash
# Terraformでインフラ構築
cd terraform/environments/prd
terraform init
terraform plan
terraform apply

# フロントエンドデプロイ
./scripts/deploy-frontend.sh
```

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成
