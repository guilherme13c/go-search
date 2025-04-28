package queue

import (
	"sync"
)

type Queue[T any] interface {
	Put(*T)
	Get() *T
	Len() int
}

type queue[T any] struct {
	lock sync.RWMutex
	data []*T
	size int
}

func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		lock: sync.RWMutex{},
		data: make([]*T, 0),
		size: 0,
	}
}

func (self *queue[T]) Put(elem *T) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.data = append(self.data, elem)
	self.size += 1
}

func (self *queue[T]) Get() *T {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if self.size == 0 {
		return nil
	}

	res := self.data[0]
	self.data = self.data[1:]
	self.size -= 1
	return res
}

func (self *queue[T]) Len() int {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.size
}
