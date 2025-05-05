package doublylinkedlist

type DoublyLinkedListNode[K any, V any] struct {
	Key   K
	Value V

	Next *DoublyLinkedListNode[K, V]
	Prev *DoublyLinkedListNode[K, V]
}

func NewDoublyLinkedListNode[K any, V any](key K, value V) *DoublyLinkedListNode[K, V] {
	return &DoublyLinkedListNode[K, V]{
		Key:   key,
		Value: value,
		Next:  nil,
		Prev:  nil,
	}
}

type DoublyLinkedList[K any, V any] struct {
	Head *DoublyLinkedListNode[K, V]
	Tail *DoublyLinkedListNode[K, V]
}

func NewDoublyLinkedList[K any, V any]() *DoublyLinkedList[K, V] {
	res := &DoublyLinkedList[K, V]{
		Head: &DoublyLinkedListNode[K, V]{},
		Tail: &DoublyLinkedListNode[K, V]{},
	}

	res.Head.Next = res.Tail
	res.Tail.Prev = res.Head

	return res
}

func (self *DoublyLinkedList[K, V]) Insert(node *DoublyLinkedListNode[K, V]) {
	tmp := self.Head.Next
	self.Head.Next = node
	node.Prev = self.Head
	node.Next = tmp
	tmp.Prev = node
}

func (self *DoublyLinkedList[K, V]) Remove(node *DoublyLinkedListNode[K, V]) {
	next := node.Next
	prev := node.Prev

	node.Prev = nil
	node.Next = nil

	prev.Next = next
	next.Prev = prev
}
