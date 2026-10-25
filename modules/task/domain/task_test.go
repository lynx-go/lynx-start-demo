package task

import (
	"errors"
	"testing"
	"time"
)

var now = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNewTask(t *testing.T) {
	task, err := NewTask("t1", "o1", "u1", "写周报", "周五前", now)
	if err != nil {
		t.Fatal(err)
	}
	if task.Completed() {
		t.Error("new task should be incomplete")
	}
	if task.CreatedAt != now || task.UpdatedAt != now {
		t.Errorf("timestamps = %v/%v", task.CreatedAt, task.UpdatedAt)
	}
}

func TestNewTaskValidation(t *testing.T) {
	for _, c := range []struct {
		name                    string
		id, org, creator, title string
	}{
		{"empty id", "", "o1", "u1", "title"},
		{"empty organization", "t1", "", "u1", "title"},
		{"empty creator", "t1", "o1", "", "title"},
	} {
		if _, err := NewTask(c.id, c.org, c.creator, c.title, "", now); err == nil {
			t.Errorf("%s: expected error", c.name)
		}
	}
	if _, err := NewTask("t1", "o1", "u1", "", "", now); !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("empty title: err = %v, want ErrEmptyTitle", err)
	}
}

func TestCompleteIdempotent(t *testing.T) {
	task, _ := NewTask("t1", "o1", "u1", "title", "", now)
	task.Complete(now.Add(time.Hour))
	first := *task.CompletedAt
	task.Complete(now.Add(2 * time.Hour)) // 幂等：不覆盖首次完成时间
	if !task.CompletedAt.Equal(first) {
		t.Errorf("CompletedAt = %v, want first %v", task.CompletedAt, first)
	}
	if !task.UpdatedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("UpdatedAt = %v, want first complete time", task.UpdatedAt)
	}
}
