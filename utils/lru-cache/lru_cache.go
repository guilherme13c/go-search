package lrucache

import (
	"sync"

	dll "github.com/guilherme13c/go-search/utils/doubly-linked-list"
)

type LruCache[K comparable, V any] interface {
	Get(K) (*V, bool)
	Put(K, V)
}

type lruCache[K comparable, V any] struct {
	capacity int

	m map[K]*dll.DoublyLinkedListNode[K, V]
	l dll.DoublyLinkedList[K, V]

	mu sync.Mutex
}

func NewLruCache[K comparable, V any](capacity uint) LruCache[K, V] {
	ll := dll.NewDoublyLinkedList[K, V]()

	return &lruCache[K, V]{
		capacity: int(capacity),
		m:        make(map[K]*dll.DoublyLinkedListNode[K, V], capacity),
		l:        *ll,
		mu:       sync.Mutex{},
	}
}

func (self *lruCache[K, V]) Get(key K) (*V, bool) {
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

func (self *lruCache[K, V]) Put(key K, value V) {
	self.mu.Lock()
	defer self.mu.Unlock()

	node, ok := self.m[key]
	if ok {
		self.l.Remove(node)
	}
	self.m[key] = dll.NewDoublyLinkedListNode(key, value)
	self.l.Insert(self.m[key])
	if len(self.m) > self.capacity {
		lru := self.l.Tail.Prev
		self.l.Remove(lru)
		delete(self.m, lru.Key)
	}
}
