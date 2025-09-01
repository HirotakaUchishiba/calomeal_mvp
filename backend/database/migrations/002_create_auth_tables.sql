-- backend/database/migrations/002_create_auth_tables.sql

-- ユーザーの基本情報を格納するテーブル
CREATE TABLE users (
    id UUID PRIMARY KEY, -- Cognitoのsubと一致させるためUUIDを維持
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ユーザーのプロフィール情報（身長、体重など）を格納するテーブル
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    height DOUBLE PRECISION,
    weight DOUBLE PRECISION,
    activity_level VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ユーザーの目標（目標体重など）を格納するテーブル
CREATE TABLE user_goals (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    target_weight DOUBLE PRECISION,
    target_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);