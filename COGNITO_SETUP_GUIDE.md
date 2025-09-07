# AWS Cognito User Pool セットアップガイド

## 概要
このドキュメントでは、CaloMeal MVPアプリケーション用のAWS Cognito User Poolの設定とデプロイ方法について説明します。

## アーキテクチャ

### Cognito User Pool
- **認証方式**: メールアドレス + パスワード
- **パスワードポリシー**: 8文字以上、大文字・小文字・数字・記号必須
- **自動確認**: メールアドレス自動確認
- **アカウント回復**: メール経由でのパスワードリセット

### Cognito User Pool Client
- **認証フロー**: USER_PASSWORD_AUTH, USER_SRP_AUTH, REFRESH_TOKEN_AUTH
- **トークン有効期限**: アクセストークン1時間、リフレッシュトークン7日
- **OAuth設定**: Authorization Code, Implicit Grant対応
- **コールバックURL**: 環境別に設定

## 環境別設定

### 開発環境 (dev)
```hcl
# コールバックURL
callback_urls = [
  "http://localhost:5173",
  "https://localhost:5173",
  "http://localhost:3000",
  "https://localhost:3000"
]

# Identity Pool: 無効
enable_identity_pool = false
```

### 本番環境 (prd)
```hcl
# コールバックURL
callback_urls = [
  "https://calomeal.com",
  "https://www.calomeal.com",
  "https://app.calomeal.com"
]

# Identity Pool: 有効
enable_identity_pool = true
```

## デプロイ手順

### 1. 前提条件
- AWS CLI設定済み
- Terraform 1.0以上インストール済み
- 適切なAWS権限（Cognito, IAM, VPC等）

### 2. 開発環境のデプロイ
```bash
# 開発環境ディレクトリに移動
cd terraform/environments/dev

# Terraform初期化
terraform init

# プラン確認
terraform plan

# デプロイ実行
terraform apply
```

### 3. 本番環境のデプロイ
```bash
# 本番環境ディレクトリに移動
cd terraform/environments/prd

# Terraform初期化
terraform init

# 変数ファイル作成
cat > terraform.tfvars << EOF
db_password = "your-secure-db-password"
jwt_access_secret = "your-jwt-access-secret"
jwt_refresh_secret = "your-jwt-refresh-secret"
cognito_custom_domain = "auth.calomeal.com"  # オプション
EOF

# プラン確認
terraform plan

# デプロイ実行
terraform apply
```

## 環境変数設定

### バックエンドアプリケーション用
```bash
# Cognito設定
COGNITO_USER_POOL_ID=ap-northeast-1_xxxxxxxxx
COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxx
COGNITO_USER_POOL_DOMAIN=calomeal-dev-auth
COGNITO_REGION=ap-northeast-1

# JWT設定
JWT_ACCESS_SECRET=your-jwt-access-secret
JWT_REFRESH_SECRET=your-jwt-refresh-secret

# 環境設定
ENVIRONMENT=development  # または production
AWS_REGION=ap-northeast-1
```

### フロントエンドアプリケーション用
```javascript
// Cognito設定
const cognitoConfig = {
  UserPoolId: 'ap-northeast-1_xxxxxxxxx',
  ClientId: 'xxxxxxxxxxxxxxxxxxxxxxxxxx',
  region: 'ap-northeast-1'
};
```

## セキュリティ設定

### パスワードポリシー
- 最小長: 8文字
- 大文字必須: あり
- 小文字必須: あり
- 数字必須: あり
- 記号必須: あり

### トークン設定
- アクセストークン: 1時間
- IDトークン: 1時間
- リフレッシュトークン: 7日

### セキュリティ機能
- トークン取り消し: 有効
- ユーザー存在エラー防止: 有効
- デバイス記憶: 有効

## ユーザー管理

### ユーザー登録フロー
1. ユーザーがメールアドレスとパスワードで登録
2. Cognitoが確認メールを送信
3. ユーザーがメール内のリンクをクリック
4. アカウントが有効化される

### 管理者によるユーザー作成
```bash
# AWS CLIを使用
aws cognito-idp admin-create-user \
  --user-pool-id ap-northeast-1_xxxxxxxxx \
  --username user@example.com \
  --user-attributes Name=email,Value=user@example.com \
  --temporary-password TempPassword123! \
  --message-action SUPPRESS
```

## 監視とログ

### CloudWatch Logs
- ユーザープールのログを有効化
- 認証イベントの監視
- エラーログの追跡

### メトリクス
- サインイン成功/失敗回数
- パスワードリセット回数
- アクティブユーザー数

## トラブルシューティング

### よくある問題

#### 1. 認証エラー
```
Error: Invalid username or password
```
**解決方法**: ユーザーが存在するか、パスワードが正しいか確認

#### 2. トークン期限切れ
```
Error: Token has expired
```
**解決方法**: リフレッシュトークンを使用して新しいトークンを取得

#### 3. コールバックURL不一致
```
Error: redirect_uri_mismatch
```
**解決方法**: User Pool ClientのコールバックURL設定を確認

### ログ確認方法
```bash
# CloudWatch Logsでログ確認
aws logs describe-log-groups --log-group-name-prefix /aws/cognito
```

## コスト最適化

### 開発環境
- 最小限のリソース使用
- 不要な機能は無効化
- 定期的なリソースクリーンアップ

### 本番環境
- 適切なスケーリング設定
- 不要なログの削除
- コスト監視の設定

## 次のステップ

1. **カスタムドメイン設定**: Route 53でのドメイン設定
2. **SES統合**: カスタムメール送信設定
3. **MFA設定**: 多要素認証の有効化
4. **SAML連携**: エンタープライズ認証の統合
5. **監視強化**: 詳細なメトリクスとアラート設定

## 参考資料

- [AWS Cognito User Pools Documentation](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools.html)
- [Terraform AWS Provider - Cognito](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cognito_user_pool)
- [Cognito Best Practices](https://docs.aws.amazon.com/cognito/latest/developerguide/best-practices.html)
