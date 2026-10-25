-- +goose Up
CREATE TABLE IF NOT EXISTS task_tasks (
    id              text PRIMARY KEY,
    organization_id text NOT NULL,
    creator_id      text NOT NULL,
    title           text NOT NULL,
    description     text NOT NULL DEFAULT '',
    completed_at    timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_task_tasks_organization
    ON task_tasks (organization_id, created_at DESC);
-- ADR-0008：organization_id 是租户键（只存 ID 不设外键，契约④）；
-- 本域无业务级唯一约束（同名任务合法）——若有，必含 organization_id。
-- creator_id 只存 ID：完成与删除的裁决依据（v1 资源所有权语义）。

-- +goose Down
DROP TABLE IF EXISTS task_tasks;
