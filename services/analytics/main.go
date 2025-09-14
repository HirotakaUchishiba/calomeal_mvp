package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HirotakaUchishiba/calomeal_mvp/proto/analytics/v1"
	"github.com/HirotakaUchishiba/calomeal_mvp/services/analytics/internal/server"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// データベース接続
	dbURL := getenv("DATABASE_URL", "postgres://postgres:password@localhost:5432/calomeal?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// データベース接続テスト
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to database")

	// gRPCサーバーの設定
	port := getenv("PORT", "50053")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	// gRPCサーバーの作成
	grpcServer := grpc.NewServer()

	// AnalyticsServiceの実装を登録
	analyticsService := server.NewAnalyticsService(db)
	analyticspb.RegisterAnalyticsServiceServer(grpcServer, analyticsService)

	log.Printf("Analytics gRPC server starting on port %s", port)

	// グレースフルシャットダウンの設定
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve:", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down analytics gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Analytics gRPC server stopped")
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
