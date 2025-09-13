#!/bin/bash

# フロントエンド本番デプロイスクリプト
# CaloMeal MVP Frontend Production Deployment Script

set -e

echo "🚀 Starting CaloMeal Frontend Production Deployment..."

# 環境変数の設定
export VITE_GRAPHQL_ENDPOINT="http://calomeal-prd-alb-461151354.ap-northeast-1.elb.amazonaws.com/query"
export VITE_COGNITO_USER_POOL_ID="ap-northeast-1_V8T3ojIla"
export VITE_COGNITO_CLIENT_ID="1llp147jgjsk9g8lgaiqtfir5h"
export VITE_AWS_REGION="ap-northeast-1"

# S3バケット名
S3_BUCKET="calomeal-frontend-prd"
CLOUDFRONT_DISTRIBUTION_ID="E3CPN565VVM3BG"

echo "📦 Building frontend for production..."
cd frontend
npm run build

echo "📤 Deploying to S3..."
aws s3 sync dist/ s3://$S3_BUCKET --delete

echo "🔄 Invalidating CloudFront cache..."
aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DISTRIBUTION_ID --paths "/*"

echo "✅ Frontend deployment completed!"
echo ""
echo "🌐 Access URLs:"
echo "  S3 Website: http://calomeal-frontend-prd.s3-website-ap-northeast-1.amazonaws.com"
echo "  CloudFront: https://d3tk27aswcyo81.cloudfront.net"
echo "  Backend API: http://calomeal-prd-alb-461151354.ap-northeast-1.elb.amazonaws.com/query"
