-- +goose Up
CREATE TABLE IF NOT EXISTS audit_records (
    id            text PRIMARY KEY,
    actor_id      text NOT NULL DEFAULT '',
    actor_email   text NOT NULL DEFAULT '',
    action        text NOT NULL,
    resource_type text NOT NULL DEFAULT '',
    resource_id   text NOT NULL DEFAULT '',
    payload       jsonb,
    source        text NOT NULL DEFAULT 'api',
    occurred_at   timestamptz NOT NULL
);
-- 审计查询三个主入口：时间轴、按主体、按动作（均带时间降序）。
CREATE INDEX IF NOT EXISTS idx_audit_records_occurred ON audit_records (occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_records_actor ON audit_records (actor_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_records_action ON audit_records (action, occurred_at DESC);

-- +goose Down
DROP TABLE IF EXISTS audit_records;
