// Package migrations 内嵌 audit 模块迁移（独立版本表 goose_audit_db，
// 契约④：迁移链隔离）。
package migrations

import "embed"

// FS 是迁移文件集（goosekit 以本目录为 root）。
//
//go:embed *.sql
var FS embed.FS
