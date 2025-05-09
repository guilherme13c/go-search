package doublylinkedlist

import (
	"testing"
)

func collectForward[K comparable, V any](dl *DoublyLinkedList[K, V]) []struct {
	Key   K
	Value V
} {
	var res []struct {
		Key   K
		Value V
	}
	for cur := dl.Head.Next; cur != dl.Tail; cur = cur.Next {
		res = append(res, struct {
			Key   K
			Value V
		}{cur.Key, cur.Value})
	}
	return res
}

func collectBackward[K comparable, V any](dl *DoublyLinkedList[K, V]) []struct {
	Key   K
	Value V
} {
	var res []struct {
		Key   K
		Value V
	}
	for cur := dl.Tail.Prev; cur != dl.Head; cur = cur.Prev {
		res = append(res, struct {
			Key   K
			Value V
		}{cur.Key, cur.Value})
	}
	return res
}

func TestNewNode(t *testing.T) {
	n := NewDoublyLinkedListNode("a", 42)
	if n.Key != "a" || n.Value != 42 {
		t.Errorf("NewDoublyLinkedListNode: got (%v,%v); want (\"a\",42)", n.Key, n.Value)
	}
	if n.Next != nil || n.Prev != nil {
		t.Error("New node should have Next and Prev == nil")
	}
}

func TestNewListSentinels(t *testing.T) {
	dl := NewDoublyLinkedList[string, int]()
	if dl.Head.Next != dl.Tail {
		t.Error("Head.Next should point to Tail")
	}
	if dl.Tail.Prev != dl.Head {
		t.Error("Tail.Prev should point to Head")
	}
}

func TestInsertSingle(t *testing.T) {
	dl := NewDoublyLinkedList[string, int]()
	n := NewDoublyLinkedListNode("x", 100)
	dl.Insert(n)

	if dl.Head.Next != n {
		t.Error("After insert, head.Next should be the new node")
	}
	if n.Prev != dl.Head || n.Next != dl.Tail {
		t.Error("Node pointers incorrect after insert")
	}
	if dl.Tail.Prev != n {
		t.Error("Tail.Prev should be the new node")
	}

	got := collectForward(dl)
	if len(got) != 1 || got[0].Key != "x" || got[0].Value != 100 {
		t.Errorf("collectForward = %v; want [(x,100)]", got)
	}
}

func TestInsertMultiple(t *testing.T) {
	dl := NewDoublyLinkedList[int, string]()
	n1 := NewDoublyLinkedListNode(1, "a")
	n2 := NewDoublyLinkedListNode(2, "b")
	n3 := NewDoublyLinkedListNode(3, "c")
	dl.Insert(n1)
	dl.Insert(n2)
	dl.Insert(n3)

	fwd := collectForward(dl)
	expectedFwd := []struct {
		Key   int
		Value string
	}{{3, "c"}, {2, "b"}, {1, "a"}}
	if len(fwd) != 3 {
		t.Fatalf("len(fwd) = %d; want 3", len(fwd))
	}
	for i, e := range expectedFwd {
		if fwd[i] != e {
			t.Errorf("fwd[%d] = %v; want %v", i, fwd[i], e)
		}
	}

	bwd := collectBackward(dl)
	expectedBwd := []struct {
		Key   int
		Value string
	}{{1, "a"}, {2, "b"}, {3, "c"}}
	if len(bwd) != 3 {
		t.Fatalf("len(bwd) = %d; want 3", len(bwd))
	}
	for i, e := range expectedBwd {
		if bwd[i] != e {
			t.Errorf("bwd[%d] = %v; want %v", i, bwd[i], e)
		}
	}
}

func TestRemoveMiddle(t *testing.T) {
	dl := NewDoublyLinkedList[int, int]()
	n1 := NewDoublyLinkedListNode(1, 10)
	n2 := NewDoublyLinkedListNode(2, 20)
	n3 := NewDoublyLinkedListNode(3, 30)

	dl.Insert(n1)
	dl.Insert(n2)
	dl.Insert(n3)

	dl.Remove(n2)

	fwd := collectForward(dl)
	if len(fwd) != 2 {
		t.Fatalf("len after remove = %d; want 2", len(fwd))
	}
	if fwd[0].Key != 3 || fwd[1].Key != 1 {
		t.Errorf("forward after remove = %v; want [3,1]", fwd)
	}

	if n3.Next != n1 {
		t.Errorf("n3.Next = %v; want %v", n3.Next, n1)
	}
	if n1.Prev != n3 {
		t.Errorf("n1.Prev = %v; want %v", n1.Prev, n3)
	}
}

func TestRemoveAll(t *testing.T) {
	dl := NewDoublyLinkedList[string, string]()
	a := NewDoublyLinkedListNode("a", "A")
	b := NewDoublyLinkedListNode("b", "B")

	dl.Insert(a)
	dl.Insert(b)

	dl.Remove(a)
	dl.Remove(b)

	if dl.Head.Next != dl.Tail || dl.Tail.Prev != dl.Head {
		t.Error("list not empty after removing all nodes")
	}

	if got := collectForward(dl); len(got) != 0 {
		t.Errorf("collectForward = %v; want empty", got)
	}
}
