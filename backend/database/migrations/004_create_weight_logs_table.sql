-- 体重記録テーブル
CREATE TABLE weight_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    weight DOUBLE PRECISION NOT NULL,
    logged_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- インデックス
CREATE INDEX idx_weight_logs_user_id ON weight_logs(user_id);
CREATE INDEX idx_weight_logs_logged_at ON weight_logs(logged_at);
CREATE INDEX idx_weight_logs_user_logged_at ON weight_logs(user_id, logged_at);

-- コメント
COMMENT ON TABLE weight_logs IS 'ユーザーの体重記録';
COMMENT ON COLUMN weight_logs.weight IS '体重 (kg)';
COMMENT ON COLUMN weight_logs.logged_at IS '記録日時';
