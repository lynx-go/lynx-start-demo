// Package demo 是样例模块的用例层：编排领域与端口（事务/事件/仓储）。
package demo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/lynx-go/lynx-start/kernel/crud"
	"github.com/lynx-go/lynx-start/kernel/events"
	"github.com/lynx-go/lynx-start/kernel/uow"
	"github.com/lynx-go/lynx-start-demo/internal/domain/demo"
)

// UseCase 是样例模块的用例集。
type UseCase struct {
	repo   demo.Repository
	tx     uow.Runner
	events events.Publisher
	now    func() time.Time
}

// New 构造用例（Assemble 调用）。
func New(repo demo.Repository, tx uow.Runner, publisher events.Publisher) *UseCase {
	return &UseCase{repo: repo, tx: tx, events: publisher, now: time.Now}
}

// Create 建项并在同一工作单元内发布事件。
func (u *UseCase) Create(ctx context.Context, title string) (*demo.Item, error) {
	var out *demo.Item
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		item, err := demo.NewItem(newID(), title, u.now().UTC())
		if err != nil {
			return err
		}
		if err := u.repo.Create(ctx, item); err != nil {
			return err
		}
		out = item
		return u.events.Publish(ctx, envelope("demo.item.created", item))
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get 取单项。
func (u *UseCase) Get(ctx context.Context, id string) (*demo.Item, error) {
	return u.repo.Get(ctx, id)
}

// Rename 改名。
func (u *UseCase) Rename(ctx context.Context, id, title string) (*demo.Item, error) {
	var out *demo.Item
	err := u.tx.Run(ctx, func(ctx context.Context) error {
		item, err := u.repo.Get(ctx, id)
		if err != nil {
			return err
		}
		if err := item.Rename(title, u.now().UTC()); err != nil {
			return err
		}
		if err := u.repo.Update(ctx, item); err != nil {
			return err
		}
		out = item
		return u.events.Publish(ctx, envelope("demo.item.renamed", item))
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Delete 删除（幂等；不存在不报错）。
func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.tx.Run(ctx, func(ctx context.Context) error {
		if err := u.repo.Delete(ctx, id); err != nil {
			return err
		}
		return u.events.Publish(ctx, events.Envelope{
			Name:    "demo.item.deleted",
			Module:  "demo",
			Payload: []byte(`{"id":"` + id + `"}`),
		})
	})
}

// ListResult 是列表用例的产物。
type ListResult struct {
	Items     []*demo.Item
	NextToken string
	TotalSize int32
}

// List 分页列表（AIP-158 offset 方言）。
func (u *UseCase) List(ctx context.Context, pageSize int32, pageToken string) (*ListResult, error) {
	params, err := crud.ParseListParams(pageSize, pageToken, "", "")
	if err != nil {
		return nil, err
	}
	items, total, err := u.repo.List(ctx, demo.ListQuery{PageSize: params.PageSize, Offset: params.Offset})
	if err != nil {
		return nil, err
	}
	page, err := crud.FinalizeOffsetPage(params, total, len(items))
	if err != nil {
		return nil, err
	}
	return &ListResult{Items: items, NextToken: page.NextToken, TotalSize: int32(total)}, nil
}

// envelope 构造事件信封（载荷最小事实集，消费方按目录文档解码）。
func envelope(name string, item *demo.Item) events.Envelope {
	return events.Envelope{
		Name:    name,
		Module:  "demo",
		Payload: []byte(`{"id":"` + item.ID + `","title":` + quote(item.Title) + `}`),
	}
}

func quote(s string) string { return `"` + s + `"` }

// newID 生成 128bit 随机 hex ID（样例模块自足实现；真实模块可用 idgen 端口）。
func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic("demo: rand: " + err.Error())
	}
	return hex.EncodeToString(b[:])
}
