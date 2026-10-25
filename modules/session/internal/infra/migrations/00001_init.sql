-- +goose Up
CREATE TABLE IF NOT EXISTS session_sessions (
    id                text PRIMARY KEY,
    user_id           text NOT NULL,
    refresh_jti       text NOT NULL,
    status            text NOT NULL DEFAULT 'active',
    created_at        timestamptz NOT NULL DEFAULT now(),
    last_refreshed_at timestamptz NOT NULL DEFAULT now(),
    expires_at        timestamptz NOT NULL,
    revoked_at        timestamptz,
    revoke_reason     text
);
CREATE INDEX IF NOT EXISTS idx_session_sessions_user ON session_sessions (user_id);
-- 契约④：不设指向 identity_users 的外键（迁移链隔离）。

-- +goose Down
DROP TABLE IF EXISTS session_sessions;
