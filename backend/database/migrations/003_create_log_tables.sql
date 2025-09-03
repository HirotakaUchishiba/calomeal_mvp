-- backend/database/migrations/003_create_log_tables.sql

-- 食事記録を格納するテーブル
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

CREATE INDEX idx_food_logs_user_id_logged_at ON food_logs (user_id, logged_at DESC);

-- 運動記録を格納するテーブル
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

COMMENT ON COLUMN exercise_logs.calories_burned IS 'ユーザーが入力した消費カロリー(キロカロリー)';

CREATE INDEX idx_exercise_logs_user_id_logged_at ON exercise_logs (user_id, logged_at DESC);