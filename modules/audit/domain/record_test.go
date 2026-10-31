package audit

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	rec, err := NewRecord("id-1", "u1", "u1@example.com", "order.created", "order", "o-9",
		json.RawMessage(`{"amount":100}`), SourceAPI, time.Unix(1000, 0))
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if rec.Source != SourceAPI || rec.Action != "order.created" {
		t.Fatalf("record = %+v", rec)
	}
	// source 缺省 SourceAPI。
	rec2, err := NewRecord("id-2", "", "", "x.y.z", "", "", nil, "", time.Unix(1, 0))
	if err != nil {
		t.Fatalf("new2: %v", err)
	}
	if rec2.Source != SourceAPI {
		t.Fatalf("default source = %q", rec2.Source)
	}
}

func TestNewRecordValidation(t *testing.T) {
	if _, err := NewRecord("id", "", "", "", "", "", nil, "", time.Unix(1, 0)); !errors.Is(err, ErrEmptyAction) {
		t.Fatalf("want ErrEmptyAction, got %v", err)
	}
	if _, err := NewRecord("id", "", "", "a.b", "", "", json.RawMessage(`{bad`), "", time.Unix(1, 0)); err == nil {
		t.Fatal("want invalid-JSON error")
	}
	if _, err := NewRecord("id", "", "", "a.b", "", "", json.RawMessage(`"`+strings.Repeat("x", MaxPayloadBytes)+`"`), "", time.Unix(1, 0)); !errors.Is(err, ErrPayloadTooLarge) {
		t.Fatalf("want ErrPayloadTooLarge, got %v", err)
	}
}
