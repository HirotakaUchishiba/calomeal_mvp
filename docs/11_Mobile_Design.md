このドキュメントは、React Native（Expo）アプリの設計・ビルド/配布戦略・API 共有方針を定義します。

## 1. 範囲と目標

- 画面: 認証、オンボーディング、ダッシュボード、記録（食事/運動/体重）
- 目標: Web と同一の GraphQL 契約を使用し、最小機能で早期バリデーション

## 2. 技術選定

- RN: Expo（TypeScript）
- 認証: Amazon Cognito（Amplify Auth）
- API: Apollo Client（GraphQL）、型は共有 codegen
- ナビゲーション: React Navigation（Stack + Tab）

## 3. プロジェクト構成（例）

```
mobile/
  app/
    screens/
    components/
    graphql/        # 共有生成物を symlink or パッケージ参照
  app.json
  package.json
```

## 4. 設定と環境切替

- `amplify-config` を `dev`/`prd` で分離し、ビルド時に環境変数で切替
- GraphQL エンドポイントは BFF の URL を参照

## 5. ビルド/配布

- Expo EAS を利用
- 内部配布: EAS build + EAS update（OTA）
- 本番: ストア配信（後続フェーズ）

## 6. テスト

- スモーク: 主要フロー（認証→記録→ダッシュボード反映）
- 将来: Detox/E2E 導入検討


