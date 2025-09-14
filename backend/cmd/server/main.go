package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/bff/middleware"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/bff/resolvers"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/analytics"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/auth"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/fooddata"
	logsvc "github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/log"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/user"
	_ "github.com/lib/pq"
	cors "github.com/rs/cors"
)

const defaultPort = "8080"

// getDSN returns database connection string from environment variables
func getDSN() string {
	// DATABASE_URLが設定されている場合はそれを使用
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	// 個別の環境変数から構築
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

	// 認証サービスを初期化
	authService := auth.NewService()

	// 認証ミドルウェアを初期化
	middleware.InitAuthMiddleware(authService)

	// gRPCクライアントの初期化
	foodServiceAddr := getenv("FOOD_SERVICE_ADDR", "localhost:50051")
	foodDataClient, err := fooddata.NewGRPCClient(foodServiceAddr)
	if err != nil {
		log.Fatal("Failed to create food gRPC client:", err)
	}
	defer foodDataClient.Close()

	// logs gRPCクライアントの初期化
	logsServiceAddr := getenv("LOGS_SERVICE_ADDR", "localhost:50052")
	logsClient, err := logsvc.NewGRPCClient(logsServiceAddr)
	if err != nil {
		log.Fatal("Failed to create logs gRPC client:", err)
	}
	defer logsClient.Close()

	// analytics gRPCクライアントの初期化
	analyticsServiceAddr := getenv("ANALYTICS_SERVICE_ADDR", "localhost:50053")
	analyticsClient, err := analytics.NewGRPCClient(analyticsServiceAddr)
	if err != nil {
		log.Fatal("Failed to create analytics gRPC client:", err)
	}
	defer analyticsClient.Close()

	// リゾルバのインスタンスを作成し、各サービスを注入（依存性の注入）
	resolver := &resolvers.Resolver{
		UserService:      user.NewService(db),
		FoodDataService:  foodDataClient,  // gRPCクライアントを使用
		LogService:       logsClient,      // gRPCクライアントを使用
		AnalyticsService: analyticsClient, // gRPCクライアントを使用
	}

	cfg := backend.Config{
		Resolvers: resolver,
		Directives: backend.DirectiveRoot{
			Auth: middleware.Auth, // ← ディレクティブを紐付け
		},
	}

	srv := handler.NewDefaultServer(backend.NewExecutableSchema(cfg))

	// --- ここからCORSの設定 ---
	// フロントエンドの開発サーバーからのアクセスを許可する設定
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:5176"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	// --- ここまでCORSの設定 ---

	// GraphQLのプレイグラウンド（開発用のテスト画面）をルートURLに設定
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// 実際のAPIエンドポイントである /query を、CORSミドルウェアを通して設定
	http.Handle("/query", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// サーバーを起動
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
