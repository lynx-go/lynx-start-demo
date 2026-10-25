-- +goose Up
CREATE TABLE IF NOT EXISTS organization_organizations (
    id         text PRIMARY KEY,
    name       text NOT NULL,
    slug       text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS organization_members (
    organization_id text NOT NULL,
    user_id         text NOT NULL,
    is_owner        boolean NOT NULL DEFAULT false,
    joined_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_organization_members_user ON organization_members (user_id);
-- 契约④：user_id 只存 ID，不设指向 identity_users 的外键（迁移链隔离）。
-- is_owner 是组织结构标记（owner-only 管理授权依据）；角色/权限语义
-- 归 access 模块，不在此落库。

-- +goose Down
DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS organization_organizations;
