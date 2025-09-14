package server

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foodspb "github.com/HirotakaUchishiba/calomeal_mvp/proto/foods/v1"
)

// FoodService implements the foods.v1.FoodService gRPC service
type FoodService struct {
	foodspb.UnimplementedFoodServiceServer
	db *sql.DB
}

// NewFoodService creates a new FoodService instance
func NewFoodService(db *sql.DB) *FoodService {
	return &FoodService{
		db: db,
	}
}

// SearchFoods searches for foods by keyword query
func (s *FoodService) SearchFoods(ctx context.Context, req *foodspb.SearchFoodsRequest) (*foodspb.SearchFoodsResponse, error) {
	// デフォルトの制限値を設定
	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// SQLクエリの構築（全文検索）
	query := `
		SELECT id, name, brand, calories, protein, carbohydrate, fat
		FROM foods 
		WHERE name ILIKE $1 OR brand ILIKE $1
		ORDER BY 
			CASE 
				WHEN name ILIKE $2 THEN 1
				WHEN brand ILIKE $2 THEN 2
				ELSE 3
			END,
			name
		LIMIT $3
	`

	// 検索パターンの準備
	searchPattern := "%" + req.Query + "%"
	exactPattern := req.Query

	rows, err := s.db.QueryContext(ctx, query, searchPattern, exactPattern, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to search foods: %v", err)
	}
	defer rows.Close()

	var foods []*foodspb.Food
	for rows.Next() {
		var food foodspb.Food
		var brand sql.NullString

		err := rows.Scan(
			&food.Id,
			&food.Name,
			&brand,
			&food.Calories,
			&food.Protein,
			&food.Carbohydrate,
			&food.Fat,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan food: %v", err)
		}

		// ブランドがnullの場合は空文字列を設定
		if brand.Valid {
			food.Brand = brand.String
		}

		foods = append(foods, &food)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating foods: %v", err)
	}

	return &foodspb.SearchFoodsResponse{
		Foods: foods,
	}, nil
}

// GetFoodById retrieves a specific food by ID
func (s *FoodService) GetFoodById(ctx context.Context, req *foodspb.GetFoodByIdRequest) (*foodspb.GetFoodByIdResponse, error) {
	if req.FoodId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Food ID is required")
	}

	query := `
		SELECT id, name, brand, calories, protein, carbohydrate, fat
		FROM foods 
		WHERE id = $1
	`

	var food foodspb.Food
	var brand sql.NullString

	err := s.db.QueryRowContext(ctx, query, req.FoodId).Scan(
		&food.Id,
		&food.Name,
		&brand,
		&food.Calories,
		&food.Protein,
		&food.Carbohydrate,
		&food.Fat,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Food with ID %s not found", req.FoodId)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get food: %v", err)
	}

	// ブランドがnullの場合は空文字列を設定
	if brand.Valid {
		food.Brand = brand.String
	}

	return &foodspb.GetFoodByIdResponse{
		Food: &food,
	}, nil
}
