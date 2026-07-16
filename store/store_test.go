package store

import (
	"context"
	"testing"
	"time"
)

func TestSet_Get(t *testing.T) {
	s := New()
	s.Set("name", "Alice", 0)

	val, ok := s.Get("name")
	if !ok {
		t.Fatal("ожидали ok=true, получили false")
	}
	if val != "Alice" {
		t.Errorf("ожидали Alice, получили %s", val)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := New()

	_, ok := s.Get("несуществующий")
	if ok {
		t.Fatal("ожидали ok=false, получили true")
	}
}

func TestSet_TTL_Expired(t *testing.T) {
	s := New()
	s.Set("token", "abc123", 50*time.Millisecond)

	_, ok := s.Get("token")
	if !ok {
		t.Fatal("ключ должен быть доступен сразу после Set")
	}

	time.Sleep(100 * time.Millisecond)

	_, ok = s.Get("token")
	if ok {
		t.Fatal("ключ должен быть недоступен после истечения TTL")
	}
}

func TestDel(t *testing.T) {
	s := New()
	s.Set("key", "value", 0)

	existed := s.Del("key")
	if !existed {
		t.Error("Del должен вернуть true для существующего ключа")
	}

	_, ok := s.Get("key")
	if ok {
		t.Error("ключ должен быть удалён")
	}

	existed = s.Del("key")
	if existed {
		t.Error("Del должен вернуть false для несуществующего ключа")
	}
}

func TestKeys(t *testing.T) {
	s := New()
	s.Set("a", "1", 0)
	s.Set("b", "2", 0)
	s.Set("c", "3", 50*time.Millisecond) // протухнет

	time.Sleep(100 * time.Millisecond)

	keys := s.Keys()

	if len(keys) != 2 {
		t.Errorf("ожидали 2 ключа, получили %d: %v", len(keys), keys)
	}
}

func TestCleaner(t *testing.T) {
	s := New()
	s.Set("temp", "value", 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.StartCleaner(ctx, 30*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	_, ok := s.Get("temp")
	if ok {
		t.Error("ключ должен быть удалён очистителем")
	}
}
