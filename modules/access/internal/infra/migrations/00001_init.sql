-- +goose Up
CREATE TABLE IF NOT EXISTS access_role_assignments (
    user_id    text NOT NULL,
    role       text NOT NULL,
    granted_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role)
);
CREATE INDEX IF NOT EXISTS idx_access_role_assignments_user ON access_role_assignments (user_id);
-- 契约④：user_id 只存 ID，不设指向 identity_users 的外键（迁移链隔离）。
-- 角色目录（名字→权限集）在配置段声明，不落库。

-- +goose Down
DROP TABLE IF EXISTS access_role_assignments;
