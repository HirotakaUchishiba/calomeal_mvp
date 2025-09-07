# JWT認証とAWS Cognito連携実装結果

## 実装完了日時
2025年9月7日

## 実装内容

### ✅ 完了した機能

#### 1. JWT認証サービス (`backend/internal/service/auth/service.go`)
- **トークンペア生成機能**: アクセストークン + リフレッシュトークン
- **トークン検証機能**: JWT署名検証と有効期限チェック
- **トークンリフレッシュ機能**: リフレッシュトークンを使用した新しいトークンペア生成
- **トークン無効化機能**: トークンの無効化（ブラックリスト対応準備済み）
- **環境変数による設定管理**: JWT秘密鍵と有効期限の設定

#### 2. AWS Cognito連携サービス (`backend/internal/service/auth/cognito.go`)
- **ユーザー登録機能**: Cognito User Poolへの新規ユーザー登録
- **登録確認機能**: メール確認コードによる登録完了
- **ログイン機能**: Cognito認証によるユーザー認証
- **ログアウト機能**: グローバルサインアウト
- **パスワードリセット機能**: パスワードリセット開始と確認
- **ユーザー情報取得機能**: アクセストークンによるユーザー情報取得

#### 3. 認証ミドルウェア (`backend/internal/bff/middleware/auth.go`)
- **JWT認証機能**: Bearer tokenの抽出・検証
- **ユーザー情報のcontext注入**: 認証されたユーザー情報をGraphQL contextに注入
- **開発環境と本番環境の切り替え**: 環境変数による認証スキップ機能
- **ヘルパー関数**: GetUserIDFromContext, GetEmailFromContext, GetClaimsFromContext

#### 4. GraphQL API統合
- **認証関連mutation**: signUp, signIn, signOut, refreshToken, resetPassword, confirmResetPassword
- **認証結果型**: AuthResult (アクセストークン、リフレッシュトークン、ユーザー情報)
- **ユーザー登録結果型**: SignUpResult (ユーザーID、確認要否、メッセージ)
- **GraphQLリゾルバー**: 全認証機能のGraphQL API実装

#### 5. サーバー統合 (`backend/cmd/server/main.go`)
- **認証サービスの初期化**: アプリケーション起動時の認証サービス設定
- **認証ミドルウェアの初期化**: 依存性注入による認証ミドルウェア設定
- **CORS設定**: Authorization headerの許可

### 📦 追加された依存関係

```go
// JWT認証
github.com/golang-jwt/jwt/v5 v5.3.0

// AWS SDK v2
github.com/aws/aws-sdk-go-v2 v1.38.3
github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.57.3
github.com/aws/aws-sdk-go-v2/config v1.31.6
```

### 🔧 環境変数設定

```bash
# JWT認証設定
JWT_ACCESS_SECRET=your-access-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key

# AWS Cognito設定
COGNITO_CLIENT_ID=your-cognito-client-id
COGNITO_USER_POOL_ID=your-cognito-user-pool-id

# 環境設定
ENVIRONMENT=development  # 開発環境では認証をスキップ
```

### 🚀 GraphQL API仕様

#### 認証関連Mutation

```graphql
# ユーザー登録
mutation SignUp($email: String!, $password: String!) {
  signUp(email: $email, password: $password) {
    userId
    email
    confirmationRequired
    message
  }
}

# 登録確認
mutation ConfirmSignUp($email: String!, $confirmationCode: String!) {
  confirmSignUp(email: $email, confirmationCode: $confirmationCode)
}

# ログイン
mutation SignIn($email: String!, $password: String!) {
  signIn(email: $email, password: $password) {
    accessToken
    refreshToken
    expiresIn
    tokenType
    user {
      id
      email
    }
  }
}

# ログアウト
mutation SignOut {
  signOut
}

# トークンリフレッシュ
mutation RefreshToken($refreshToken: String!) {
  refreshToken(refreshToken: $refreshToken) {
    accessToken
    refreshToken
    expiresIn
    tokenType
    user {
      id
      email
    }
  }
}

# パスワードリセット開始
mutation ResetPassword($email: String!) {
  resetPassword(email: $email)
}

# パスワードリセット確認
mutation ConfirmResetPassword($email: String!, $confirmationCode: String!, $newPassword: String!) {
  confirmResetPassword(email: $email, confirmationCode: $confirmationCode, newPassword: $newPassword)
}
```

### 🔒 セキュリティ機能

1. **JWT署名検証**: HMAC-SHA256による署名検証
2. **トークン有効期限**: アクセストークン1時間、リフレッシュトークン7日
3. **Bearer token認証**: Authorization headerによる認証
4. **CORS設定**: 適切なオリジン制限
5. **環境分離**: 開発環境と本番環境の認証切り替え

### 📊 実装統計

- **新規ファイル**: 2個
- **修正ファイル**: 4個
- **追加されたコード行数**: 約500行
- **実装された機能**: 7個の認証関連機能
- **GraphQL mutation**: 7個
- **GraphQL type**: 2個

### 🎯 次のステップ

1. **AWS Cognito User Poolの設定**: 実際のCognito環境の構築
2. **環境変数の設定**: 本番環境での適切な設定値の設定
3. **フロントエンド統合**: Reactアプリケーションでの認証機能統合
4. **E2Eテスト**: 認証フローの自動テスト実装
5. **セキュリティ強化**: トークンブラックリスト、レート制限等

### ✅ 実装完了確認

- [x] JWT認証サービスの実装
- [x] AWS Cognito連携サービスの実装
- [x] 認証ミドルウェアの実装
- [x] GraphQL API統合
- [x] サーバー統合
- [x] コンパイルエラーの修正
- [x] 依存関係の追加
- [x] ドキュメント化

**認証システムの実装が完了しました！** 🎉

これにより、本格的な認証機能を持つMVPアプリケーションの基盤が完成しました。
