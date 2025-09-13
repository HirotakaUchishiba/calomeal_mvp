# 本番環境デプロイ情報

## 🌐 本番環境アクセス情報

### フロントエンド
- **S3静的ウェブサイト**: http://calomeal-frontend-prd.s3-website-ap-northeast-1.amazonaws.com
- **CloudFront**: https://d3tk27aswcyo81.cloudfront.net
- **ステータス**: ✅ 動作中

### バックエンド
- **GraphQL API**: http://calomeal-prd-alb-461151354.ap-northeast-1.elb.amazonaws.com/query
- **ヘルスチェック**: http://calomeal-prd-alb-461151354.ap-northeast-1.elb.amazonaws.com/health
- **GraphQL Playground**: http://calomeal-prd-alb-461151354.ap-northeast-1.elb.amazonaws.com/health
- **ステータス**: ✅ 動作中

## 🔧 インフラ構成

### AWS リソース
- **ECS Fargate**: calomeal-prd-cluster
- **ALB**: calomeal-prd-alb
- **RDS**: PostgreSQL (calomeal-subnet-group)
- **S3**: calomeal-frontend-prd
- **CloudFront**: E3CPN565VVM3BG
- **Cognito**: ap-northeast-1_V8T3ojIla

### 認証設定
- **User Pool ID**: ap-northeast-1_V8T3ojIla
- **Client ID**: 1llp147jgjsk9g8lgaiqtfir5h
- **Region**: ap-northeast-1

## 🚀 デプロイ手順

### フロントエンドデプロイ
```bash
# デプロイスクリプトを実行
./scripts/deploy-frontend.sh
```

### 手動デプロイ
```bash
# 1. フロントエンドをビルド
cd frontend
VITE_GRAPHQL_ENDPOINT=http://calomeal-prd-alb-461151354.ap-northeast-1.elb.amazonaws.com/query \
VITE_COGNITO_USER_POOL_ID=ap-northeast-1_V8T3ojIla \
VITE_COGNITO_CLIENT_ID=1llp147jgjsk9g8lgaiqtfir5h \
VITE_AWS_REGION=ap-northeast-1 \
npm run build

# 2. S3にデプロイ
aws s3 sync dist/ s3://calomeal-frontend-prd --delete

# 3. CloudFrontキャッシュを無効化
aws cloudfront create-invalidation --distribution-id E3CPN565VVM3BG --paths "/*"
```

## 📊 本番環境デプロイ状況

| 項目 | 状況 | 詳細 |
|------|------|------|
| バックエンド | ✅ 100% | GraphQL API 正常動作 |
| フロントエンド | ✅ 100% | S3 + CloudFront で配信 |
| 認証システム | ✅ 100% | Cognito + JWT 統合 |
| データベース | ✅ 100% | RDS PostgreSQL 稼働中 |
| インフラ | ✅ 100% | ECS + ALB + VPC 構成 |
| SSL設定 | ✅ 95% | 設定完了、ドメイン検証待ち |

## 🎯 次のステップ

1. **ユーザーテスト開始**: 実際のユーザーによる価値検証
2. **gRPCマイクロサービス化**: フェーズ2の実装
3. **React Native アプリ**: フェーズ7の実装
4. **観測性・信頼性向上**: フェーズ5の実装

## 🔍 トラブルシューティング

### フロントエンドが表示されない場合
1. S3バケットのパブリックアクセス設定を確認
2. CloudFrontの配信状況を確認
3. ブラウザのキャッシュをクリア

### API接続エラーの場合
1. ALBのヘルスチェックを確認
2. ECSサービスの状態を確認
3. セキュリティグループの設定を確認

## 📝 更新履歴

- **2025-09-13**: フロントエンド本番デプロイ完了
- **2025-09-13**: バックエンド本番デプロイ完了
- **2025-09-13**: SSL/HTTPS設定完了
