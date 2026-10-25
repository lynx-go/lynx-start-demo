package notification

import (
	"testing"
	"time"
)

var now = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNewNotification(t *testing.T) {
	n, err := NewNotification("n1", "u1", IdentityUserPasswordChanged, "密码已更改", "正文", map[string]string{"id": "u1"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if n.Unread() != true {
		t.Errorf("new notification should be unread")
	}
	if n.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", n.CreatedAt, now)
	}
}

func TestNewNotificationValidation(t *testing.T) {
	cases := []struct {
		name       string
		id, userID string
		typ, title string
		wantErr    bool
	}{
		{"empty id", "", "u1", "t", "title", true},
		{"empty userID", "n1", "", "t", "title", true},
		{"empty type", "n1", "u1", "", "title", true},
		{"empty title", "n1", "u1", "t", "", true},
		{"ok minimal", "n1", "u1", "t", "title", false},
	}
	for _, c := range cases {
		_, err := NewNotification(c.id, c.userID, c.typ, c.title, "", nil, now)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err = %v, wantErr %v", c.name, err, c.wantErr)
		}
	}
}

func TestMarkReadIdempotent(t *testing.T) {
	n, _ := NewNotification("n1", "u1", "t", "title", "", nil, now)
	n.MarkRead(now.Add(time.Minute))
	if n.Unread() {
		t.Fatal("should be read after MarkRead")
	}
	first := *n.ReadAt
	n.MarkRead(now.Add(2 * time.Minute)) // 幂等：不覆盖首次已读时间
	if !n.ReadAt.Equal(first) {
		t.Errorf("ReadAt = %v, want first %v", n.ReadAt, first)
	}
}
