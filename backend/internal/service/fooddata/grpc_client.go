package fooddata

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	foodspb "github.com/HirotakaUchishiba/calomeal_mvp/proto/foods/v1"
)

// GRPCClient implements the fooddata Service interface using gRPC
type GRPCClient struct {
	client foodspb.FoodServiceClient
	conn   *grpc.ClientConn
}

// NewGRPCClient creates a new gRPC client for the food service
func NewGRPCClient(foodServiceAddr string) (*GRPCClient, error) {
	// gRPC接続の確立
	conn, err := grpc.Dial(foodServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to food service: %w", err)
	}

	client := foodspb.NewFoodServiceClient(conn)

	return &GRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// SearchFood searches for foods using the gRPC service
func (c *GRPCClient) SearchFood(ctx context.Context, query string) ([]Food, error) {
	req := &foodspb.SearchFoodsRequest{
		Query: query,
		Limit: 20,
	}

	resp, err := c.client.SearchFoods(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC search foods failed: %w", err)
	}

	// gRPCレスポンスを内部のFood型に変換
	var foods []Food
	for _, pbFood := range resp.Foods {
		food := Food{
			ID:           int64(0), // gRPCでは文字列IDなので、必要に応じて変換
			Name:         pbFood.Name,
			Brand:        &pbFood.Brand,
			Calories:     pbFood.Calories,
			Protein:      pbFood.Protein,
			Carbohydrate: pbFood.Carbohydrate,
			Fat:          pbFood.Fat,
		}
		foods = append(foods, food)
	}

	return foods, nil
}

// GetFoodByID retrieves a food by ID using the gRPC service
func (c *GRPCClient) GetFoodByID(ctx context.Context, foodID string) (*Food, error) {
	req := &foodspb.GetFoodByIdRequest{
		FoodId: foodID,
	}

	resp, err := c.client.GetFoodById(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC get food by ID failed: %w", err)
	}

	// gRPCレスポンスを内部のFood型に変換
	food := &Food{
		ID:           0, // gRPCでは文字列IDなので、必要に応じて変換
		Name:         resp.Food.Name,
		Brand:        &resp.Food.Brand,
		Calories:     resp.Food.Calories,
		Protein:      resp.Food.Protein,
		Carbohydrate: resp.Food.Carbohydrate,
		Fat:          resp.Food.Fat,
	}

	return food, nil
}
