# E2Eテスト

このディレクトリには、CaloMeal MVPアプリケーションのEnd-to-Endテストが含まれています。

## テストツール

- **Playwright**: ブラウザ自動化テストフレームワーク
- **TypeScript**: テストコードの記述言語

## テストケース

### 1. ハッピーパス (happy-path.spec.ts)
- **E2E-HP-001**: 新規ユーザー登録から食事記録までの完全なフロー
- 設計書で定義された主要なユーザージャーニーをテスト

### 2. 認証機能 (auth.spec.ts)
- ログイン機能
- サインアップ機能
- パスワードリセット機能

### 3. 記録機能 (logging.spec.ts)
- 食事記録機能
- 運動記録機能
- 体重記録機能
- 日付ナビゲーター機能

## テストの実行

### 前提条件
1. アプリケーションが起動していること
   ```bash
   docker-compose up -d
   ```

2. フロントエンドが起動していること
   ```bash
   cd frontend && npm run dev
   ```

### 実行コマンド

```bash
# 全テストを実行
npm run test:e2e

# UIモードでテストを実行（推奨）
npm run test:e2e:ui

# ヘッドモードでテストを実行（ブラウザが表示される）
npm run test:e2e:headed

# デバッグモードでテストを実行
npm run test:e2e:debug

# 特定のテストファイルを実行
npx playwright test tests/e2e/happy-path.spec.ts

# 特定のブラウザでテストを実行
npx playwright test --project=chromium
```

## テストデータ

テストでは以下のテストユーザーを使用します：
- **メール**: `testuser@example.com`
- **パスワード**: `TestPassword123!`

## CI/CD統合

これらのテストはCI/CDパイプラインに統合され、以下のタイミングで実行されます：
- プルリクエスト作成時
- developブランチへのマージ時
- 本番デプロイ前

## トラブルシューティング

### テストが失敗する場合
1. アプリケーションが正常に起動しているか確認
2. データベースの状態を確認
3. テストデータが正しく設定されているか確認

### タイムアウトエラー
- アプリケーションの起動に時間がかかる場合は、`playwright.config.ts`の`timeout`設定を調整

### セレクターエラー
- フロントエンドのUIが変更された場合は、テストのセレクターを更新
