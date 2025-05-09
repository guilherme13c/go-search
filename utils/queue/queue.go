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

type randomizedQueue[T any] struct {
	mu     sync.RWMutex
	data   []T
	size   int
	random *rand.Rand
}

func NewQueue[T any](seed int64) Queue[T] {
	return &randomizedQueue[T]{
		mu:     sync.RWMutex{},
		data:   make([]T, 0),
		size:   0,
		random: rand.New(rand.NewSource(seed)),
	}
}

func (self *randomizedQueue[T]) Put(elem T) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.data = append(self.data, elem)
	self.size += 1
}

func (self *randomizedQueue[T]) Get() (T, bool) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.size == 0 {
		return *new(T), false
	}

	idx := self.random.Intn(self.size)
	elem := self.data[idx]

	self.data[idx] = self.data[self.size-1]
	self.data = self.data[:self.size-1]
	self.size--

	return elem, true
}

func (self *randomizedQueue[T]) Len() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.size
}
