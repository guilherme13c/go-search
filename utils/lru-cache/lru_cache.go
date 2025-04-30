package lrucache

import (
	dll "github.com/guilherme13c/go-search/utils/doubly-linked-list"
)

type lrucache[K comparable, V any] struct {
	capacity uint
	size     uint

	m map[K]*dll.DoublyLinkedListNode[K, V]
	l dll.DoublyLinkedList[K, V]
}

func (self *lrucache[K, V]) Get(key K) (*V, bool) {
	node, ok := self.m[key]
	if !ok {
		return nil, false
	}

	self.l.Remove(node)
	self.l.Insert(node)

	return &node.Value, true
}

func (self *lrucache[K, V]) Put(key K, value V) {
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
	newNode := &dll.DoublyLinkedListNode[K, V]{Key: key, Value: value}
	self.m[key] = newNode
	self.l.Insert(newNode)
	self.size++
}
