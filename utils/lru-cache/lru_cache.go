package lrucache

import (
	"sync"

	dll "github.com/guilherme13c/go-search/utils/doubly-linked-list"
)

type LRUCache[K comparable, V any] struct {
	capacity uint
	size     uint

	m map[K]*dll.DoublyLinkedListNode[K, V]
	l dll.DoublyLinkedList[K, V]

	mu sync.Mutex
}

func NewLRUCache[K comparable, V any](capacity uint) *LRUCache[K, V] {
	ll := dll.NewDoublyLinkedList[K, V]()

	return &LRUCache[K, V]{
		capacity: capacity,
		size:     0,
		m:        make(map[K]*dll.DoublyLinkedListNode[K, V], capacity),
		l:        *ll,
		mu:       sync.Mutex{},
	}
}

func (self *LRUCache[K, V]) Get(key K) (*V, bool) {
	self.mu.Lock()
	defer self.mu.Unlock()

	node, ok := self.m[key]
	if !ok {
		return nil, false
	}

	self.l.Remove(node)
	self.l.Insert(node)

	return &node.Value, true
}

func (self *LRUCache[K, V]) Put(key K, value V) {
	self.mu.Lock()
	defer self.mu.Unlock()

	node, ok := self.m[key]
	if ok {
		self.l.Remove(node)
		self.l.Insert(node)
		node.Value = value
	}
	if self.size >= self.capacity {
		lru := self.l.Tail.Prev
		self.l.Remove(lru)
		delete(self.m, lru.Key)
		self.size--
	}
	newNode := dll.NewDoublyLinkedListNode(key, value)
	self.m[key] = newNode
	self.l.Insert(newNode)
	self.size++
}
