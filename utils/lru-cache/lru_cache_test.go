package lrucache

import (
	"testing"
)

func TestGetMiss(t *testing.T) {
	cache := NewLruCache[string, bool](32)
	if v, ok := cache.Get("foo"); ok || v != nil {
		t.Errorf("Get on empty cache = (%v, %v); want (nil, false)", v, ok)
	}
}

func TestPutAndGetHit(t *testing.T) {
	cache := NewLruCache[string, string](32)
	cache.Put("a", "alpha")
	cache.Put("b", "beta")

	if v, ok := cache.Get("a"); !ok || *v != "alpha" {
		t.Errorf("Get(\"a\") = (%v, %v); want (\"alpha\", true)", v, ok)
	}
	if v, ok := cache.Get("b"); !ok || *v != "beta" {
		t.Errorf("Get(\"b\") = (%v, %v); want (\"beta\", true)", v, ok)
	}
}

func TestUpdateExisting(t *testing.T) {
	cache := NewLruCache[int, int](2)
	cache.Put(1, 100)
	cache.Put(2, 200)
	
	if v, ok := cache.Get(1); !ok || *v != 100 {
		t.Errorf("After update, Get(1) = (%v, %v); want (100, true)", v, ok)
	}

	cache.Put(1, 101)
	cache.Put(3, 300)

	if v, ok := cache.Get(1); !ok || *v != 101 {
		t.Errorf("After update, Get(1) = (%v, %v); want (101, true)", v, ok)
	}
	if _, ok := cache.Get(2); ok {
		t.Error("Expected key 2 to be evicted, but still present")
	}
	if v, ok := cache.Get(3); !ok || *v != 300 {
		t.Errorf("Get(3) = (%v, %v); want (300, true)", v, ok)
	}
}

func TestEvictionOrder(t *testing.T) {
	cache := NewLruCache[string, string](3)
	cache.Put("a", "A")
	cache.Put("b", "B")
	cache.Put("c", "C")
	if _, _ = cache.Get("b"); true {
	}
	if _, _ = cache.Get("a"); true {
	}
	cache.Put("d", "D")
	if _, ok := cache.Get("c"); ok {
		t.Error("Expected \"c\" to be evicted")
	}
	cache.Put("e", "E")
	if _, ok := cache.Get("b"); ok {
		t.Error("Expected \"b\" to be evicted")
	}
	for _, key := range []string{"a", "d", "e"} {
		if v, ok := cache.Get(key); !ok {
			t.Errorf("Expected key %q to be present", key)
		} else {
			want := map[string]string{"a": "A", "d": "D", "e": "E"}[key]
			if *v != want {
				t.Errorf("Get(%q) = %q; want %q", key, *v, want)
			}
		}
	}
}
