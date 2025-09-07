このドキュメントは、栄養管理アプリ「Lean MVP」の技術的な設計と構成を定義します。

## 2.1. ハイレベルアーキテクチャ

アーキテクチャは、React SPA、Go GraphQL BFF、Goサービス群から成る三層構造を維持します。この「関心の分離」は、将来の拡張性を見据えた上で有効です。

- **React SPA (Single Page Application)**: ユーザーが直接操作するフロントエンド。
- **Go GraphQL BFF (Backend for Frontend)**: フロントエンドからのデータ取得要求に特化したAPIレイヤー。
- **Goサービス群**: ビジネスロジックを担うバックエンドサービス。

## 2.2. 「構造化モノリス」戦略

MVPのデプロイメント戦略として、複数の論理サービス（`fooddata`, `log`, `user`）を単一のコンテナにパッケージングする「**構造化モノリス**」アプローチを採用します。

これにより、MVPフェーズにおける運用上の複雑性（サービス間通信、分散トレーシングなど）を大幅に削減できます。将来的に特定サービスの負荷が高まった場合、その部分だけを独立したマイクロサービスとして切り出すことが容易な、進化的な設計です。

## 2.3. 技術スタックの最適化

高レベルな技術選定は維持しつつ、MVPの実装速度を最大化するために内部の通信方式とデータ層の実装を簡素化します。

| カテゴリ | 技術 | 根拠 |
| :--- | :--- | :--- |
| **フロントエンド** | React, TypeScript | 変更なし。モダンなSPA開発の標準。 |
| **バックエンド言語** | Go | 変更なし。マイクロサービスに適した言語特性。 |
| **APIレイヤー** | GraphQL (BFF) | 変更なし。フロントエンドのデータ取得を最適化。 |
| **サービス間通信** | **Goインターフェースによる直接呼び出し** | **変更。**構造化モノリス内でのgRPC通信は不要なオーバーヘッド。直接的な関数呼び出しで実装を大幅に簡素化。 |
| **データベース** | PostgreSQL (Amazon RDS) | 変更なし。信頼性の高いRDB。 |
| **インフラストラクチャ** | AWS (ECS, Fargate, S3, ALB, VPC) | 変更なし。スケーラブルなクラウド基盤。 |
| **認証基盤** | Amazon Cognito | 変更なし。セキュアな認証機能のマネージドサービス。 |
| **Infrastructure as Code** | Terraform | 変更なし。インフラの再現性と一貫性を確保。 |
| **CI/CD** | GitHub Actions | 変更なし。開発サイクルの自動化。 |

## 2.4. プロジェクトディレクトリ構成

以下は、本プロジェクトのルートディレクトリからの構成です。

```plaintext
.
├── .github/
│   └── workflows/
│       └── deploy.yml      # CI/CDパイプライン定義
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go     # アプリケーション起動エントリーポイント
│   ├── database/
│   │   ├── migrations/     # DBマイグレーションファイル
│   │   └── seeds/
│   │       └── foods.sql   # 初期投入用の食品データ
│   ├── internal/
│   │   ├── bff/
│   │   │   └── resolvers/  # GraphQLリゾルバ
│   │   └── service/        # ビジネスロジック層 (fooddata/, log/, user/)
│   ├── go.mod
│   └── schema.graphqls     # GraphQLスキーマ定義
├── frontend/
│   ├── public/
│   └── src/
│       ├── components/     # 共通UIコンポーネント
│       └── features/       # 機能ごとのコンポーネント群 (auth/, dashboard/, records/)
├── terraform/
│   ├── environments/       # 環境ごとの設定 (dev/, prd/)
│   └── modules/            # 再利用可能なインフラ構成要素 (ecs/, rds/, vpc/, etc.)
├── .env.example            # ローカル開発用の環境変数テンプレート
├── .gitignore
└── docker-compose.yml      # ローカル開発環境定義
```