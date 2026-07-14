package store

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		ttl      time.Duration
		expected string
		wantOk   bool
	}{
		{
			name:     "simple set and get",
			key:      "key1",
			value:    "hello",
			ttl:      0,
			expected: "hello",
			wantOk:   true,
		},
		{
			name:     "expired key",
			key:      "key2",
			value:    "world",
			ttl:      -1 * time.Second,
			expected: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := New()
			st.Set(tt.key, tt.value, tt.ttl)
			got, ok := st.Get(tt.key)

			if ok != tt.wantOk {
				t.Errorf("Get() ok = %v, want %v", ok, tt.wantOk)
			}

			if got != tt.expected {
				t.Errorf("Get() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDel(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(st *Store)
		key       string
		wantDel   bool
		wantGetOk bool
	}{
		{
			name: "delete existing key",
			setup: func(st *Store) {
				st.Set("a", "1", 0)
			},
			key:       "a",
			wantDel:   true,
			wantGetOk: false,
		},
		{
			name: "delete non-existing key",
			setup: func(st *Store) {
				st.Set("a", "1", 0)
			},
			key:       "b",
			wantDel:   false,
			wantGetOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := New()
			tt.setup(st)
			got := st.Del(tt.key)
			if got != tt.wantDel {
				t.Errorf("Del() = %v, want %v", got, tt.wantDel)
			}
			_, ok := st.Get(tt.key)
			if ok != tt.wantGetOk {
				t.Errorf("Get() after Del ok = %v, want %v", ok, tt.wantGetOk)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(st *Store)
		expected []string
	}{
		{
			name: "only alive keys",
			setup: func(st *Store) {
				st.Set("a", "1", 0)              // бессрочно
				st.Set("b", "2", 10*time.Second) // ещё не истекло
				st.Set("c", "3", -1*time.Second) // уже истекло
			},
			expected: []string{"a", "b"},
		},
		{
			name: "all keys expired",
			setup: func(st *Store) {
				st.Set("a", "1", -1*time.Hour)
				st.Set("b", "2", -1*time.Minute)
			},
			expected: []string{},
		},
		{
			name:     "empty store",
			setup:    func(st *Store) {},
			expected: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := New()
			tt.setup(st)
			keys := st.Keys()
			if len(keys) != len(tt.expected) {
				t.Errorf("Keys() len = %d, want %d, got %v", len(keys), len(tt.expected), keys)
				return
			}

			keySet := make(map[string]bool)
			for _, k := range keys {
				keySet[k] = true
			}
			for _, k := range tt.expected {
				if !keySet[k] {
					t.Errorf("Keys() missing expected key: %s, got %v", k, keys)
				}
			}
		})
	}
}
