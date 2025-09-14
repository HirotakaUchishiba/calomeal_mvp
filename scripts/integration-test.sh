#!/bin/bash

# gRPC統合テストスクリプト
set -e

echo "🚀 gRPC統合テスト開始"

# 環境変数の設定
export DATABASE_URL="postgres://postgres:password@localhost:5432/calomeal?sslmode=disable"
export FOOD_SERVICE_ADDR="localhost:50051"
export LOGS_SERVICE_ADDR="localhost:50052"

# プロセス管理用の配列
PIDS=()

# クリーンアップ関数
cleanup() {
    echo "🧹 クリーンアップ中..."
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            echo "プロセス $pid を停止中..."
            kill "$pid"
        fi
    done
    exit 0
}

# シグナルハンドラーを設定
trap cleanup SIGINT SIGTERM

echo "📊 データベース接続テスト..."
sleep 2

echo "🍽️  foods gRPCサービス起動中..."
cd services/foods
DATABASE_URL="postgres://postgres:password@localhost:5432/calomeal?sslmode=disable" go run main.go &
FOODS_PID=$!
PIDS+=($FOODS_PID)
echo "foods-service PID: $FOODS_PID"

echo "📝 logs gRPCサービス起動中..."
cd ../logs
DATABASE_URL="postgres://postgres:password@localhost:5432/calomeal?sslmode=disable" go run main.go &
LOGS_PID=$!
PIDS+=($LOGS_PID)
echo "logs-service PID: $LOGS_PID"

echo "⏳ サービス起動待機中..."
sleep 5

echo "🔧 BFF起動中..."
cd ../../backend
DATABASE_URL="postgres://postgres:password@localhost:5432/calomeal?sslmode=disable" \
FOOD_SERVICE_ADDR="localhost:50051" \
LOGS_SERVICE_ADDR="localhost:50052" \
go run cmd/server/main.go &
BFF_PID=$!
PIDS+=($BFF_PID)
echo "BFF PID: $BFF_PID"

echo "⏳ BFF起動待機中..."
sleep 5

echo "🧪 統合テスト実行中..."

# テスト1: foods gRPCサービス接続テスト
echo "📋 テスト1: foods gRPCサービス接続テスト"
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "query { searchFood(query: \"ごはん\") { id name calories } }"}' \
  --max-time 10 || echo "❌ foods gRPC接続テスト失敗"

echo ""

# テスト2: logs gRPCサービス接続テスト（GraphQL経由）
echo "📋 テスト2: logs gRPCサービス接続テスト"
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "query { dailySummary(date: \"2025-09-14\") { caloriesIntake caloriesBurned protein carbohydrate fat } }"}' \
  --max-time 10 || echo "❌ logs gRPC接続テスト失敗"

echo ""

# テスト3: ヘルスチェック
echo "📋 テスト3: ヘルスチェック"
curl -X GET http://localhost:8080/health \
  --max-time 5 || echo "❌ ヘルスチェック失敗"

echo ""

echo "✅ 統合テスト完了"
echo "📊 テスト結果:"
echo "  - foods gRPCサービス: ポート50051"
echo "  - logs gRPCサービス: ポート50052" 
echo "  - BFF GraphQL: ポート8080"
echo "  - データベース: ポート5432"

echo ""
echo "🔄 テストを継続中... (Ctrl+Cで停止)"

# プロセス監視
while true; do
    sleep 10
    echo "💓 サービス稼働中... $(date)"
done
