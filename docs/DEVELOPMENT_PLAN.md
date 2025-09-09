# 開発工程計画（フェーズ0–8）

本計画は、設計書の要件（GraphQL BFF + gRPC マイクロサービス + React Native + AWS 運用）を満たすための工程分解です。各フェーズは「作業」「成果物」「受け入れ基準（AC）」を明記します。

## フェーズ0 基盤の安定化
- 作業: ローカル起動、E2Eハッピーパス、最小構造化ログ（traceId）
- 成果物: 起動手順更新、最小ログ
- AC: E2E緑、ログにtraceId

## フェーズ1 契約の固定化（GraphQL）
- 作業: スキーマ唯一化、型生成パイプライン、互換性ルール
- 成果物: PRチェック項目、codegenスクリプト
- AC: スキーマ変更時にCIで型生成/影響検出

## フェーズ2 gRPC導入（最小）
- 作業: Protoツールチェーン、foods.v1.SearchFoods 実装、BFF差替
- 成果物: .proto / 生成 / サーバ / クライアント
- AC: ローカルで gRPC 経由 searchFood 成功、テスト追加

## フェーズ3 インフラ dev デプロイ
- 作業: Cloud Map + 複数ECS、ECR、CIマトリクス
- 成果物: dev環境に BFF/foodsvc 配置
- AC: ALB経由でGraphQL、BFF→foodsvc疎通

## フェーズ4 サービス拡張（logsvc/analytics）
- 作業: logs.v1 / analytics.v1 抽出、BFFをgRPC化、JWTメタデータ伝播
- 成果物: 2サービスの .proto/実装、BFFクライアント
- AC: 既存E2E緑、コントラクトテスト

## フェーズ5 観測性・信頼性
- 作業: OpenTelemetry Collector、タイムアウト/リトライ、CloudWatch/X-Ray
- 成果物: トレース/メトリクス可視化、基本アラート
- AC: p95可視、トレースで追跡可能

## フェーズ6 データ責務と移行
- 作業: テーブル所有権の明文化、スキーマ分離ロードマップ
- 成果物: 所有マッピング、移行手順
- AC: 参照/更新境界が明確

## フェーズ7 React Native 最小アプリ
- 作業: Expo + Amplify + Apollo、dailySummary 表示
- 成果物: mobile/ 最小構成
- AC: dev環境BFFからデータ取得/表示

## フェーズ8 ガバナンス/品質
- 作業: 契約互換性ゲート（GraphQL/Proto）、SAST、Secrets運用
- 成果物: 互換性ブロックCI、セキュリティ基本
- AC: 非互換PRはCIで検出

---

進め方と詳細手順は docs/LESSONS.md および各レッスン文書を参照。
