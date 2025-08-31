-- backend/database/migrations/001_create_core_tables.sql

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

-- キーワード検索を高速化するためのインデックス
CREATE INDEX idx_foods_name ON foods USING GIN (to_tsvector('simple', name));