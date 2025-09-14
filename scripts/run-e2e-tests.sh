#!/bin/bash

# E2Eテスト実行スクリプト
# 使用方法: ./scripts/run-e2e-tests.sh

set -e

echo "🧪 E2Eテスト実行開始"
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

# データベース起動確認
echo "📊 データベース起動確認中..."
if ! docker exec calomeal_mvp-db-1 psql -U postgres -d calomeal -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ データベースが起動していません。Docker Composeで起動してください。"
    echo "   docker-compose up -d db"
    exit 1
fi
echo "✅ データベース起動確認完了"

# バックエンドサービス起動
echo "🔧 バックエンドサービス起動中..."

# foods gRPCサービス起動
echo "🍽️  foods gRPCサービス起動中..."
cd services/foods
go run main.go &
FOODS_PID=$!
PIDS+=($FOODS_PID)
echo "foods-service PID: $FOODS_PID"
cd ../..

sleep 2

# logs gRPCサービス起動
echo "📝 logs gRPCサービス起動中..."
cd services/logs
go run main.go &
LOGS_PID=$!
PIDS+=($LOGS_PID)
echo "logs-service PID: $LOGS_PID"
cd ../..

sleep 2

# analytics gRPCサービス起動
echo "📊 analytics gRPCサービス起動中..."
cd services/analytics
go run main.go &
ANALYTICS_PID=$!
PIDS+=($ANALYTICS_PID)
echo "analytics-service PID: $ANALYTICS_PID"
cd ../..

sleep 2

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
sleep 10

# サービス接続確認
echo "🔍 サービス接続確認中..."

# BFF GraphQLサービス接続確認
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ BFF GraphQLサービス接続成功"
else
    echo "❌ BFF GraphQLサービス接続失敗"
    exit 1
fi

# フロントエンド起動
echo "🌐 フロントエンド起動中..."
cd frontend
npm run dev &
FRONTEND_PID=$!
PIDS+=($FRONTEND_PID)
echo "フロントエンド PID: $FRONTEND_PID"
cd ..

# フロントエンド起動待機
echo "⏳ フロントエンド起動待機中..."
sleep 15

# フロントエンド接続確認
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo "✅ フロントエンド接続成功"
else
    echo "❌ フロントエンド接続失敗"
    exit 1
fi

echo ""
echo "🎭 Playwright E2Eテスト実行開始"
echo "================================================"

# Playwrightテスト実行
cd /Users/hirotaka/Desktop/calomeal_mvp

# テスト実行
echo "📋 認証機能テスト実行中..."
npx playwright test tests/e2e/auth.spec.ts --reporter=html

echo "📋 ハッピーパステスト実行中..."
npx playwright test tests/e2e/happy-path.spec.ts --reporter=html

echo "📋 ログ機能テスト実行中..."
npx playwright test tests/e2e/logging.spec.ts --reporter=html

echo "📋 Analytics機能テスト実行中..."
npx playwright test tests/e2e/analytics.spec.ts --reporter=html

echo ""
echo "🎉 E2Eテスト完了！"
echo "================================================"
echo "📊 テスト結果レポート:"
echo "  - HTMLレポート: playwright-report/index.html"
echo "  - テスト結果: test-results/"
echo ""
echo "🔍 テスト対象機能:"
echo "  - 認証機能（ログイン・サインアップ・パスワードリセット）"
echo "  - ハッピーパス（新規ユーザー登録から食事記録まで）"
echo "  - ログ機能（食事・運動・体重記録）"
echo "  - Analytics機能（栄養サマリー・トレンド・体重進捗・カロリーバランス）"
echo ""
echo "⏹️  テスト終了時は Ctrl+C を押してください"

# ユーザーがCtrl+Cを押すまで待機
wait
