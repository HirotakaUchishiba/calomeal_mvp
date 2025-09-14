package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	foodspb "github.com/HirotakaUchishiba/calomeal_mvp/proto/foods/v1"
	"github.com/HirotakaUchishiba/calomeal_mvp/services/foods/internal/server"
	_ "github.com/lib/pq"
)

const defaultPort = "50051"

// getDSN returns database connection string from environment variables
func getDSN() string {
	host := getenv("DB_HOST", "db")
	port := getenv("DB_PORT", "5432")
	user := getenv("POSTGRES_USER", "postgres")
	pass := getenv("POSTGRES_PASSWORD", "postgres")
	name := getenv("POSTGRES_DB", "calomeal")
	ssl := getenv("DB_SSLMODE", "disable")
	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=" + ssl
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// DB接続
	db, err := sql.Open("postgres", getDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// gRPCサーバーの作成
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// FoodServiceの実装を登録
	foodService := server.NewFoodService(db)
	foodspb.RegisterFoodServiceServer(s, foodService)

	// リフレクションを有効化（開発用）
	reflection.Register(s)

	log.Printf("🚀 Foods service starting on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
