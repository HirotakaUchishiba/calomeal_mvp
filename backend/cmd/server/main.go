package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/bff/resolvers"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/fooddata"
	"github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/user"
	cors "github.com/rs/cors"
	logsvc "github.com/HirotakaUchishiba/calomeal_mvp/backend/internal/service/log" 
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// リゾルバのインスタンスを作成し、各サービスを注入（依存性の注入）
	resolver := &resolvers.Resolver{
		UserService:     user.NewService(),
		FoodDataService: fooddata.NewService(),
		LogService:      logsvc.NewService(),
	}

	srv := handler.NewDefaultServer(backend.NewExecutableSchema(backend.Config{Resolvers: resolver}))

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