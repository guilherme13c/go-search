package set

import (
	"crypto/md5"
	"encoding/json"
)

type Set[T any] interface {
	Add(T)
	Remove(T)
	Contains(T) bool
}

type set[T any] struct {
	data map[string]struct{}
}

func NewSet[T any]() Set[T] {
	return &set[T]{
		data: map[string]struct{}{},
	}
}

func (self *set[T]) Add(element T) {
	k := self.getKey(element)

	self.data[k] = struct{}{}
}

func (self *set[T]) Remove(element T) {
	k := self.getKey(element)

	delete(self.data, k)
}

func (self *set[T]) Contains(element T) bool {
	k := self.getKey(element)

	_, ok := self.data[k]

	return ok
}

func (self *set[T]) getKey(element T) string {
	b, err := json.Marshal(element)
	if err != nil {
		panic(err)
	}

	h := md5.New()

	h.Write(b)

	return string(h.Sum(nil))
}
