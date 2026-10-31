package set

import "maps"

type Set[T comparable] map[T]struct{}

func New[T comparable](elems ...T) Set[T] {
	s := make(Set[T], len(elems))
	s.Add(elems...)
	return s
}

func (s Set[T]) Add(elems ...T) {
	for _, e := range elems {
		s[e] = struct{}{}
	}
}

func (s Set[T]) Remove(elems ...T) {
	for _, e := range elems {
		delete(s, e)
	}
}

func (s Set[T]) Contains(e T) bool {
	_, ok := s[e]
	return ok
}

func (s Set[T]) Equals(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}
	for x := range s {
		if !other.Contains(x) {
			return false
		}
	}
	return true
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}
