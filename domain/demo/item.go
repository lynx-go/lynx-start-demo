// Package demo 是样例模块的领域层：聚合、不变量与仓储端口。
// 契约：不 import infra/app/api；跨模块只允许 import 他模块 domain 包。
package demo

import (
	"time"

	"github.com/lynx-go/lynx-start/kernel/errcode"
)

// space 是本模块的错误码区段（契约⑤）。
var space = errcode.MustRegister("demo")

// 哨兵错误：app 层转 grpc/status 时携带 ErrorInfo（Domain=demo）。
var (
	ErrItemNotFound = space.New("ITEM_NOT_FOUND", "demo item not found")
	ErrEmptyTitle   = space.New("EMPTY_TITLE", "title must not be empty")
)

// Item 是样例聚合（单实体；不变量：title 非空）。
type Item struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewItem 构造新聚合。
func NewItem(id, title string, now time.Time) (*Item, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	return &Item{ID: id, Title: title, CreatedAt: now, UpdatedAt: now}, nil
}

// Rename 改名（更新时间戳由领域维护）。
func (i *Item) Rename(title string, now time.Time) error {
	if title == "" {
		return ErrEmptyTitle
	}
	i.Title = title
	i.UpdatedAt = now
	return nil
}

// ListQuery 是列表查询参数（分页方言见 kernel/crud）。
type ListQuery struct {
	PageSize int32
	Offset   int
}
