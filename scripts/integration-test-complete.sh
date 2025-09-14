#!/bin/bash

# 統合テストスクリプト - 全サービスの動作確認
# 使用方法: ./scripts/integration-test-complete.sh

set -e

echo "🚀 統合テスト開始: 全サービスの動作確認"
echo "================================================"

# 環境変数設定
export DATABASE_URL="postgres://postgres:password@localhost:5432/calomeal?sslmode=disable"
export FOOD_SERVICE_ADDR="localhost:50051"
export LOGS_SERVICE_ADDR="localhost:50052"
export ANALYTICS_SERVICE_ADDR="localhost:50053"

# プロセスIDを保存する配列
PIDS=()

# クリーンアップ関数
cleanup() {
    echo ""
    echo "🧹 クリーンアップ中..."
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            echo "プロセス $pid を終了中..."
            kill "$pid"
        fi
    done
    echo "✅ クリーンアップ完了"
}

# シグナルハンドラーを設定
trap cleanup EXIT INT TERM

# データベース接続確認
echo "📊 データベース接続確認中..."
if ! docker exec calomeal_mvp-db-1 psql -U postgres -d calomeal -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ データベース接続に失敗しました"
    exit 1
fi
echo "✅ データベース接続成功"

# foods gRPCサービス起動
echo "🍽️  foods gRPCサービス起動中..."
cd services/foods
go run main.go &
FOODS_PID=$!
PIDS+=($FOODS_PID)
echo "foods-service PID: $FOODS_PID"
cd ../..

# 少し待機
sleep 2

# logs gRPCサービス起動
echo "📝 logs gRPCサービス起動中..."
cd services/logs
go run main.go &
LOGS_PID=$!
PIDS+=($LOGS_PID)
echo "logs-service PID: $LOGS_PID"
cd ../..

# 少し待機
sleep 2

# analytics gRPCサービス起動
echo "📊 analytics gRPCサービス起動中..."
cd services/analytics
go run main.go &
ANALYTICS_PID=$!
PIDS+=($ANALYTICS_PID)
echo "analytics-service PID: $ANALYTICS_PID"
cd ../..

# 少し待機
sleep 3

# BFF GraphQLサービス起動
echo "🔧 BFF GraphQLサービス起動中..."
cd backend
go run cmd/server/main.go &
BFF_PID=$!
PIDS+=($BFF_PID)
echo "BFF PID: $BFF_PID"
cd ..

# サービス起動待機
echo "⏳ サービス起動待機中..."
sleep 5

# サービス接続テスト
echo ""
echo "🔍 サービス接続テスト開始"
echo "================================================"

# foodsサービス接続テスト
echo "🍽️  foodsサービス接続テスト..."
if curl -s http://localhost:50051/health > /dev/null 2>&1; then
    echo "✅ foodsサービス接続成功"
else
    echo "❌ foodsサービス接続失敗"
fi

# logsサービス接続テスト
echo "📝 logsサービス接続テスト..."
if curl -s http://localhost:50052/health > /dev/null 2>&1; then
    echo "✅ logsサービス接続成功"
else
    echo "❌ logsサービス接続失敗"
fi

# analyticsサービス接続テスト
echo "📊 analyticsサービス接続テスト..."
if curl -s http://localhost:50053/health > /dev/null 2>&1; then
    echo "✅ analyticsサービス接続成功"
else
    echo "❌ analyticsサービス接続失敗"
fi

# BFF GraphQLサービス接続テスト
echo "🔧 BFF GraphQLサービス接続テスト..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ BFF GraphQLサービス接続成功"
else
    echo "❌ BFF GraphQLサービス接続失敗"
fi

# GraphQLクエリテスト
echo ""
echo "🔍 GraphQLクエリテスト開始"
echo "================================================"

# 基本的なGraphQLクエリテスト
echo "📊 基本的なGraphQLクエリテスト..."
GRAPHQL_QUERY='{"query": "query { __schema { types { name } } }"}'
if curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$GRAPHQL_QUERY" \
    http://localhost:8080/query > /dev/null 2>&1; then
    echo "✅ GraphQLスキーマクエリ成功"
else
    echo "❌ GraphQLスキーマクエリ失敗"
fi

# フロントエンド起動
echo ""
echo "🌐 フロントエンド起動中..."
cd frontend
npm run dev &
FRONTEND_PID=$!
PIDS+=($FRONTEND_PID)
echo "フロントエンド PID: $FRONTEND_PID"
cd ..

# フロントエンド起動待機
echo "⏳ フロントエンド起動待機中..."
sleep 10

# フロントエンド接続テスト
echo "🌐 フロントエンド接続テスト..."
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo "✅ フロントエンド接続成功"
else
    echo "❌ フロントエンド接続失敗"
fi

echo ""
echo "🎉 統合テスト完了！"
echo "================================================"
echo "📊 サービス状況:"
echo "  - PostgreSQL: http://localhost:5432"
echo "  - foods gRPC: http://localhost:50051"
echo "  - logs gRPC: http://localhost:50052"
echo "  - analytics gRPC: http://localhost:50053"
echo "  - BFF GraphQL: http://localhost:8080"
echo "  - フロントエンド: http://localhost:5173"
echo ""
echo "🔍 テスト方法:"
echo "  1. ブラウザで http://localhost:5173 にアクセス"
echo "  2. ログインしてダッシュボードを確認"
echo "  3. アナリティクスページでanalytics機能をテスト"
echo ""
echo "⏹️  テスト終了時は Ctrl+C を押してください"

# ユーザーがCtrl+Cを押すまで待機
wait
