// backend/internal/service/fooddata/service.go
package fooddata

import "context"

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
	SearchFood(ctx context.Context, query string) (Food, error)
	// TODO: GetFoodByIDのようなメソッドも後で必要になります
}

type service struct {
	// TODO: ここにデータベース接続(リポジトリ)を保持します
}

// NewServiceは新しいfooddataサービスインスタンスを作成します
func NewService() Service {
	return &service{}
}

// SearchFoodはキーワードで食品を検索します
func (s *service) SearchFood(ctx context.Context, query string) (Food, error) {
	// 【メンターズノート】
	// ここで、設計資料で指定されたPostgreSQLの全文検索を実装します。
	// `to_tsvector`でドキュメントをベクトル化し、`to_tsquery`でクエリをベクトル化し、
	// `@@`演算子でマッチングを行います。GINインデックスにより、この検索は非常に高速です。
	// 実際のDBライブラリ(GORMなど)では以下のようなクエリを生成します:
	// SELECT * FROM foods WHERE to_tsvector('simple', name) @@ to_tsquery('simple', 'your_query') LIMIT 20;

	// 現時点ではダミーデータを返します
	// TODO: 実際のデータベース検索ロジックを実装
	if query == "ごはん" {
		return Food{
			ID: 2, Name: "ごはん", Brand: nil, Calories: 168, Protein: 2.5, Carbohydrate: 37, Fat: 0.3,
		}, nil
	}
	return Food{}, nil
}
