package set

import (
	"testing"
)

type Person struct {
	Name string
	Age  int
}

type Bad struct {
	ch chan int
}

func TestIntSet_AddContainsRemove(t *testing.T) {
	s := NewSet[int]()

	if s.Contains(1) {
		t.Error("empty set should not contain 1")
	}

	s.Add(1)
	if !s.Contains(1) {
		t.Error("set should contain 1 after Add(1)")
	}

	s.Add(1)
	if !s.Contains(1) {
		t.Error("set should still contain 1 after Add(1) again")
	}

	s.Remove(1)
	if s.Contains(1) {
		t.Error("set should not contain 1 after Remove(1)")
	}

	s.Remove(2)
	if s.Contains(2) {
		t.Error("set should not contain 2")
	}
}

func TestStringSet(t *testing.T) {
	s := NewSet[string]()
	s.Add("foo")
	s.Add("bar")
	if !s.Contains("foo") || !s.Contains("bar") {
		t.Error("set should contain both \"foo\" and \"bar\"")
	}
	s.Remove("foo")
	if s.Contains("foo") {
		t.Error("set should not contain \"foo\" after removal")
	}
}

func TestStructSet(t *testing.T) {
	s := NewSet[Person]()
	alice := Person{Name: "Alice", Age: 30}
	bob := Person{Name: "Bob", Age: 25}

	s.Add(alice)
	s.Add(bob)

	if !s.Contains(alice) {
		t.Error("set should contain Alice")
	}
	if !s.Contains(bob) {
		t.Error("set should contain Bob")
	}

	s.Remove(alice)
	if s.Contains(alice) {
		t.Error("set should not contain Alice after removal")
	}
}
