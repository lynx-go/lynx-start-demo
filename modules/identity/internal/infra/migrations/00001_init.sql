-- +goose Up
CREATE TABLE IF NOT EXISTS identity_users (
    id            text PRIMARY KEY,
    email         text NOT NULL,
    password_hash text NOT NULL,
    display_name  text NOT NULL,
    status        text NOT NULL DEFAULT 'active',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);
-- 归一化（小写）由领域层强制；唯一索引保证兜底。
CREATE UNIQUE INDEX IF NOT EXISTS idx_identity_users_email ON identity_users (email);

-- +goose Down
DROP TABLE IF EXISTS identity_users;
