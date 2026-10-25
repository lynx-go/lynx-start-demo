-- +goose Up
CREATE TABLE IF NOT EXISTS notification_notifications (
    id         text PRIMARY KEY,
    user_id    text NOT NULL,
    type       text NOT NULL,
    title      text NOT NULL,
    body       text NOT NULL DEFAULT '',
    payload    jsonb NOT NULL DEFAULT '{}'::jsonb,
    read_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_notification_notifications_user
    ON notification_notifications (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notification_notifications_unread
    ON notification_notifications (user_id) WHERE read_at IS NULL;
-- 契约④：user_id 只存 ID，不设跨模块外键（迁移链隔离）。
-- payload 是事件原始载荷快照（跳转等用途）；title/body 是渲染快照。

-- +goose Down
DROP TABLE IF EXISTS notification_notifications;
