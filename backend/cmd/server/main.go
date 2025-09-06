package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/bff/resolvers"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/fooddata"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/user"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/bff/middleware"
	cors "github.com/rs/cors"
	logsvc "github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/log" 
)

const defaultPort = "8080"

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

	// リゾルバのインスタンスを作成し、各サービスを注入（依存性の注入）
	resolver := &resolvers.Resolver{
		UserService:     user.NewService(),
		FoodDataService: fooddata.NewService(),
		LogService:      logsvc.NewService(db),
	}

	cfg := backend.Config{
		Resolvers: resolver,
		Directives: backend.DirectiveRoot{
		  Auth: middleware.Auth, // ← ディレクティブを紐付け
		},
	  }

	srv := handler.NewDefaultServer(backend.NewExecutableSchema(cfg))

	// --- ここからCORSの設定 ---
	// フロントエンドの開発サーバーである http://localhost:5173 からのアクセスを許可する設定
	c := cors.New(cors.Options{
		AllowedOrigins:[]string{"http://localhost:5173"},
		AllowedMethods:[]string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:[]string{"Authorization", "Content-Type"},
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