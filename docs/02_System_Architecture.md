このドキュメントは、栄養管理アプリ「Lean MVP」の技術的な設計と構成を定義します。

## 2.1. ハイレベルアーキテクチャ

アーキテクチャは、React/React Native クライアント、Go GraphQL BFF、ドメイン別マイクロサービス（gRPC）から成る構造です。GraphQL はクライアントとの公開契約、gRPC はサービス間契約として機能します。

- **React SPA / React Native**: ユーザーが操作するクライアント。GraphQL 経由で BFF と通信。
- **Go GraphQL BFF (Backend for Frontend)**: クライアントに最適化された API レイヤー。Resolver は gRPC クライアントとして各サービスを呼び出す。
- **ドメイン別マイクロサービス（gRPC）**: `usersvc`, `foodsvc`, `logsvc`, `analytics` 等のサービスが gRPC で連携。

## 2.2. マイクロサービス + gRPC 戦略

ドメイン境界を意識したサービス分割を行い、サービス間通信は gRPC を採用します。BFF は公開契約（GraphQL）を維持しつつ、内部実装は gRPC 経由で各サービスをオーケストレーションします。サービスディスカバリに AWS Cloud Map を用い、プライベートネットワーク内での安全な通信を確保します。

## 2.3. 技術スタックの最適化

サービス間契約を gRPC（Protocol Buffers）に統一し、スケーラブルで相互運用性の高い構成とします。

| カテゴリ | 技術 | 根拠 |
| :--- | :--- | :--- |
| **フロントエンド** | React, TypeScript | モダンなSPA開発の標準。 |
| **バックエンド言語** | Go | マイクロサービスに適した言語特性。 |
| **APIレイヤー** | GraphQL (BFF) | フロントエンドのデータ取得を最適化。 |
| **サービス間通信** | **gRPC (Protocol Buffers)** | 明確な契約と多言語互換性、パフォーマンス。 |
| **データベース** | PostgreSQL (Amazon RDS) | 信頼性の高いRDB。 |
| **インフラストラクチャ** | AWS (ECS Fargate, S3, ALB, VPC, Cloud Map) | サービスディスカバリに Cloud Map を使用。 |
| **認証基盤** | Amazon Cognito | セキュアな認証機能のマネージドサービス。 |
| **Infrastructure as Code** | Terraform | インフラの再現性と一貫性を確保。 |
| **CI/CD** | GitHub Actions | 開発サイクルの自動化。 |
| **観測性** | OpenTelemetry, CloudWatch Logs, X-Ray | 分散トレーシング/メトリクス/ログの統合。 |

## 2.4. プロジェクトディレクトリ構成

以下は、本プロジェクトのルートディレクトリからの構成例です（要点のみ）。

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
│   └── schema.graphql      # GraphQLスキーマ定義
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
├── proto/                  # gRPCの.proto定義（サービス間契約）
├── services/               # マイクロサービス（将来的に分割配置）
└── mobile/                 # React Native アプリ

## 2.5. リクエストフロー（例）

1. クライアント（Web/RN）が GraphQL にリクエスト
2. BFF の Resolver が Cognito JWT を検証（`@auth` ディレクティブ）
3. Resolver が対応する gRPC サービスにリクエスト（JWT 由来のユーザー情報をメタデータとして伝播）
4. サービスが DB にアクセスし結果を返却
5. BFF が GraphQL レスポンスとしてクライアントに返却
```