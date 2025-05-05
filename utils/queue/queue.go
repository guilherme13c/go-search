package queue

import (
	"math/rand"
	"sync"
)

type Queue[T any] interface {
	Put(T)
	Get() (T, bool)
	Len() int
}

type queue[T any] struct {
	mu   sync.RWMutex
	data []T
	size int
}

func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		mu:   sync.RWMutex{},
		data: make([]T, 0),
		size: 0,
	}
}

func (self *queue[T]) Put(elem T) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.data = append(self.data, elem)
	self.size += 1
}

func (q *queue[T]) Get() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return *new(T), false
	}

	idx := rand.Intn(q.size)
	elem := q.data[idx]

	q.data[idx] = q.data[q.size-1]
	q.data = q.data[:q.size-1]
	q.size--

	return elem, true
}

func (self *queue[T]) Len() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.size
}
