package organization

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNewOrganization(t *testing.T) {
	if _, err := NewOrganization("o1", "", "acme", t0); !errors.Is(err, ErrEmptyName) {
		t.Errorf("empty name: err = %v, want ErrEmptyName", err)
	}
	for _, bad := range []string{"", "ab", "Acme", "acme_corp", "-acme", "acme-", strings.Repeat("a", 65)} {
		if _, err := NewOrganization("o1", "Acme Inc", bad, t0); !errors.Is(err, ErrInvalidSlug) {
			t.Errorf("slug %q: err = %v, want ErrInvalidSlug", bad, err)
		}
	}
	for _, good := range []string{"abc", "acme", "acme-corp", "acme-2", "a" + strings.Repeat("b", 62) + "c"} {
		if _, err := NewOrganization("o1", "Acme Inc", good, t0); err != nil {
			t.Errorf("slug %q: err = %v, want ok", good, err)
		}
	}
	o, err := NewOrganization("o1", "Acme Inc", "acme", t0)
	if err != nil {
		t.Fatal(err)
	}
	if o.CreatedAt != t0 || o.UpdatedAt != t0 {
		t.Errorf("timestamps = %v/%v, want t0", o.CreatedAt, o.UpdatedAt)
	}
}

func TestRename(t *testing.T) {
	o, err := NewOrganization("o1", "Acme Inc", "acme", t0)
	if err != nil {
		t.Fatal(err)
	}
	t1 := t0.Add(time.Hour)
	if err := o.Rename("Acme Corporation", t1); err != nil {
		t.Fatal(err)
	}
	if o.Name != "Acme Corporation" || o.UpdatedAt != t1 {
		t.Errorf("rename result = %+v", o)
	}
	if err := o.Rename("", t1); !errors.Is(err, ErrEmptyName) {
		t.Errorf("empty name: err = %v, want ErrEmptyName", err)
	}
}

func TestNewMember(t *testing.T) {
	if _, err := NewMember("", "u1", true, t0); err == nil {
		t.Error("empty orgID: want error")
	}
	if _, err := NewMember("o1", "", false, t0); err == nil {
		t.Error("empty userID: want error")
	}
	m, err := NewMember("o1", "u1", true, t0)
	if err != nil {
		t.Fatal(err)
	}
	if !m.IsOwner || m.JoinedAt != t0 {
		t.Errorf("member = %+v", m)
	}
}
