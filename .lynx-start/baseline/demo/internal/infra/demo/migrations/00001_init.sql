-- +goose Up
CREATE TABLE IF NOT EXISTS demo_items (
    id         text PRIMARY KEY,
    title      text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_demo_items_created_at ON demo_items (created_at DESC, id DESC);

-- +goose Down
DROP TABLE IF EXISTS demo_items;
