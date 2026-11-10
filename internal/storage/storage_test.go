package storage

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)
	v, ok := s.Get("foo")
	if !ok || v != "bar" {
		t.Fatalf("Get = %q, %v; want bar, true", v, ok)
	}
}

func TestGetMissing(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)
	s.Delete("foo")
	_, ok := s.Get("foo")
	if ok {
		t.Fatal("key should be deleted")
	}
}

func TestTTLExpires(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 100*time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	_, ok := s.Get("foo")
	if ok {
		t.Fatal("key should be expired")
	}
}

func TestTTLValue(t *testing.T) {
	s := New()
	s.Set("foo", "bar", time.Hour)
	if ttl := s.TTL("foo"); ttl <= 0 {
		t.Fatalf("TTL = %d; want > 0", ttl)
	}
}

func TestTTLNoKey(t *testing.T) {
	s := New()
	if ttl := s.TTL("missing"); ttl != -2 {
		t.Fatalf("TTL = %d; want -2", ttl)
	}
}

func TestTTLNoExpiry(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)
	if ttl := s.TTL("foo"); ttl != -1 {
		t.Fatalf("TTL = %d; want -1", ttl)
	}
}
