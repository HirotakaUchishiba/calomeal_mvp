このドキュメントは、マイクロサービス分割と gRPC 契約の設計・実装・運用方針を定義します。

## 1. サービス境界

- **usersvc**: ユーザー作成/取得、オンボーディング管理（Cognito と連携）
- **foodsvc**: 内部食品DB検索（全文検索）、食品メタデータ管理
- **logsvc**: 食事/運動/体重の記録 CRUD
- **analytics**: 集計/レポート（日次サマリー）

## 2. Proto 設計原則

- 後方互換性: フィールド番号は不変、削除は非推奨化→段階的削除
- エラー: gRPC ステータスを準拠して返す
- セキュリティ: JWT 由来の `x-user-id` をメタデータ必須、必要に応じ `x-email`
- 観測性: トレースIDをメタデータで伝播（`x-trace-id`）

## 3. サンプル .proto（抜粋）

```proto
syntax = "proto3";
package foods.v1;
option go_package = "github.com/yourorg/yourrepo/proto/foods/v1;foodspb";

service FoodService {
  rpc SearchFoods(SearchFoodsRequest) returns (SearchFoodsResponse) {}
}

message SearchFoodsRequest {
  string query = 1;
  int32 limit = 2;
}

message Food {
  string id = 1;
  string name = 2;
  string brand = 3;
  double calories = 4;
  double protein = 5;
  double carbohydrate = 6;
  double fat = 7;
}

message SearchFoodsResponse {
  repeated Food foods = 1;
}
```

## 4. メタデータ規約

- `authorization`: Bearer トークン（BFF→サービスで必要に応じ検証/信頼境界に応じて省略可）
- `x-user-id`: 必須、Cognito の `sub`
- `x-email`: 任意
- `x-trace-id`: 任意（なければ生成）

## 5. タイムアウト/リトライ/回路遮断

- クライアントは 1-2 秒のデッドライン設定を必須
- `Unavailable` のみ指数バックオフで 2-3 回リトライ
- サーキットブレーカー/バルクヘッドは将来導入検討

## 6. デプロイ/サービスディスカバリ

- 各サービスは ECS Fargate 上で稼働、Cloud Map に登録
- ネットワークはプライベートサブネット、BFF のみ ALB 経由で公開

## 7. 監査/ログ

- すべてのミューテーションは `x-user-id` と操作内容を構造化ログで記録

## 8. バージョニング

- 破壊的変更は `v2` パッケージとして追加。並行稼働期間を設ける


