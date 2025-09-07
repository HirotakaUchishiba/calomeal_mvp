このドキュメントは、「Lean MVPブループリント」で定義されたアーキテクチャに基づき、サーバーサイドの実装仕様を詳細に定義します。

## 4.1. APIコントラクト定義 (GraphQLスキーマ)

Lean MVPのスコープに合わせて、GraphQLスキーマを大幅に簡素化します。バーコード検索、カスタムメニュー関連のクエリとミューテーションは完全に削除します。

```graphql
# backend/schema.graphqls

directive @auth on FIELD_DEFINITION

type Query {
  "指定された日付のユーザーの日次サマリーを取得します"
  dailySummary(date: String!): DailySummary! @auth

  "キーワードで食品を検索します (内部データベース)"
  searchFood(query: String!): [FoodSearchResult!]! @auth
}

type Mutation {
  "ユーザーの初回プロフィールと目標を設定します"
  completeOnboarding(profile: UserProfileInput!, goal: UserGoalInput!): User! @auth

  "食事を記録します"
  logFood(input: LogFoodInput!): FoodLog! @auth

  "運動を記録します (消費カロリーは手動入力)"
  logExercise(input: LogExerciseInput!): ExerciseLog! @auth

  "体重を記録します"
  logWeight(weight: Float!, date: String!): WeightLog! @auth
}

# (User, DailySummary, FoodLog, FoodSearchResult, 各種Input型などの型定義は省略)
# ...
```

## 4.2. サービスロジックの簡素化

ビジネス仮説の検証速度を優先するため、各サービスの責務を簡素化します。

* **FoodDataServiceの責務変更**:
    * 外部API (Open Food Facts)との通信ロジックは**完全に削除**します。
    * このサービスは、内部の`foods`テーブルを検索するシンプルなリポジトリとして機能します。

* **LogServiceの責務変更**:
    * **運動記録**: METS法に基づく消費カロリー自動計算ロジックを**完全に削除**します。`logExercise`ハンドラは、ユーザーが入力した値をそのままデータベースに保存するだけの単純な処理となります。
    * **食事記録**: ユーザーの摂取量に応じた栄養価の計算ロジックは維持します。これはアプリケーションのコア機能です。

## 4.3. データベーススキーマ (物理定義)

開発速度とシンプルさを優先し、スキーマを簡素化します。カスタムメニュー関連のテーブルは削除し、内部食品データベース用の`foods`テーブルを新たに追加します。

#### `foods` テーブル (新規追加)

```sql
-- MVP用の内部食品データベース
CREATE TABLE foods (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    brand TEXT,
    -- 100gあたりの栄養価
    calories DOUBLE PRECISION NOT NULL,
    protein DOUBLE PRECISION NOT NULL,
    carbohydrate DOUBLE PRECISION NOT NULL,
    fat DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- キーワード検索を高速化するためのGINインデックス
CREATE INDEX idx_foods_name ON foods USING GIN (to_tsvector('simple', name));
```

#### `users` テーブル

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY, -- Cognitoのsubと一致させるためUUIDを維持
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `food_logs` テーブル (簡素化)

```sql
CREATE TABLE food_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_name TEXT NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    unit VARCHAR(50) NOT NULL,
    calories DOUBLE PRECISION NOT NULL,
    protein DOUBLE PRECISION NOT NULL,
    carbohydrate DOUBLE PRECISION NOT NULL,
    fat DOUBLE PRECISION NOT NULL,
    logged_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `exercise_logs` テーブル (簡素化)

```sql
CREATE TABLE exercise_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_name TEXT NOT NULL,
    duration_minutes INT NOT NULL,
    calories_burned DOUBLE PRECISION NOT NULL, -- ユーザー入力値を格納
    logged_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
*(user_profiles, user_goals, weight_logs テーブルも同様に主キーと数値型を簡素化)*

## 4.4. 詳細ビジネスロジック仕様

### 4.4.1. 食品検索 (`searchFood`クエリ)

* **検索アルゴリズム**: PostgreSQLの全文検索機能を利用します。具体的には `to_tsvector('simple', name) @@ to_tsquery('simple', :query)` の形式でクエリを実行します。
* **結果の制限**: パフォーマンスとUXを考慮し、検索結果は最大20件に制限します。

### 4.4.2. 食事記録時の栄養計算 (`logFood`ミューテーション)

* **前提**: `foods`テーブルに格納されている栄養価は、すべて100gあたりの値です。
* **計算式**:
    ```
    recorded_calories = (foods.calories / 100) * input.quantity
    ```
    (タンパク質、脂質、炭水化物も同様)
* この計算結果を`food_logs`テーブルの各栄養素カラムに保存します。

### 4.4.3. 日次サマリー (`dailySummary`クエリ)

* **データ取得範囲**: 指定された日付の`00:00:00`から`23:59:59`までのに`logged_at`を持つ`food_logs`と`exercise_logs`をユーザーIDで絞り込み、取得します。
* **集計ロジック**:
    * **総摂取カロリー**: `food_logs`の`calories`カラムの合計値。
    * **総消費カロリー**: `exercise_logs`の`calories_burned`カラムの合計値。
    * **PFC摂取量**: `food_logs`の`protein`, `fat`, `carbohydrate`カラムのそれぞれの合計値。
* これらの集計結果を`DailySummary`型としてフロントエンドに返却します。