package queue

import (
	"reflect"
	"testing"
)

func TestPutAndLen(t *testing.T) {
	q := NewQueue[int](0)

	if got := q.Len(); got != 0 {
		t.Errorf("Len initially = %d; want 0", got)
	}

	q.Put(10)
	q.Put(20)
	if got := q.Len(); got != 2 {
		t.Errorf("Len after 2 puts = %d; want 2", got)
	}
}

func TestGetEmpty(t *testing.T) {
	q := NewQueue[int](0)
	v, ok := q.Get()
	if ok {
		t.Errorf("Get on empty queue returned ok=true, v=%q; want ok=false", v)
	}
}

func TestGetReducesLenAndReturnsItem(t *testing.T) {
	q := NewQueue[int](0)
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		q.Put(v)
	}

	seen := make(map[int]bool)
	originalLen := q.Len()

	for i := range originalLen {
		v, ok := q.Get()
		if !ok {
			t.Fatalf("expected ok=true on Get(); got false at iteration %d", i)
		}
		seen[v] = true
		if got := q.Len(); got != originalLen-(i+1) {
			t.Errorf("Len after %d Gets = %d; want %d", i+1, got, originalLen-(i+1))
		}
	}

	for _, v := range values {
		if !seen[v] {
			t.Errorf("value %d was never returned by Get()", v)
		}
	}

	if _, ok := q.Get(); ok {
		t.Error("expected ok=false on Get() after draining queue")
	}
}

func TestDeterministicRemovalOrder(t *testing.T) {
	q := NewQueue[int](0)
	for i := range 5 {
		q.Put(i)
	}

	var order []int
	for q.Len() > 0 {
		v, _ := q.Get()
		order = append(order, v)
	}

	expected := []int{4, 2, 1, 0, 3}
	if !reflect.DeepEqual(order, expected) {
		t.Errorf("removal order = %v; want %v", order, expected)
	}
}

