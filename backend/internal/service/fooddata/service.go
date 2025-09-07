// backend/internal/service/fooddata/service.go
package fooddata

import (
	"context"
	"database/sql"
)

// Foodはデータベースのfoodsテーブルのレコードを表します
type Food struct {
	ID           int64
	Name         string
	Brand        *string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
}

// Serviceは食品データ関連のビジネスロジックのインターフェースです
type Service interface {
	SearchFood(ctx context.Context, query string) ([]Food, error)
	// TODO: GetFoodByIDのようなメソッドも後で必要になります
}

type service struct {
	db *sql.DB
}

// NewServiceは新しいfooddataサービスインスタンスを作成します
func NewService(db *sql.DB) Service {
	return &service{db: db}
}

// SearchFoodはキーワードで食品を検索します
func (s *service) SearchFood(ctx context.Context, query string) ([]Food, error) {
	// PostgreSQLの全文検索とLIKE検索を組み合わせて食品を検索
	// 日本語検索に対応するため、LIKE検索も併用
	const searchQuery = `
		SELECT id, name, brand, calories, protein, carbohydrate, fat
		FROM foods 
		WHERE to_tsvector('simple', name) @@ to_tsquery('simple', $1)
		   OR name ILIKE $2
		ORDER BY name
		LIMIT 20
	`

	likePattern := "%" + query + "%"
	rows, err := s.db.QueryContext(ctx, searchQuery, query, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []Food
	for rows.Next() {
		var food Food
		err := rows.Scan(
			&food.ID,
			&food.Name,
			&food.Brand,
			&food.Calories,
			&food.Protein,
			&food.Carbohydrate,
			&food.Fat,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return foods, nil
}
