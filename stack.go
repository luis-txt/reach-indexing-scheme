package main

import (
	"fmt"
	"os"
)

type Stack[T any] struct {
	data []T
}

func CreateStack[T any](n int) *Stack[T] {
	return &Stack[T]{
		make([]T, 0, max(n, 1)),
	}
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) ClearStack() {
	s.data = s.data[:0]
}

func (s *Stack[T]) Pop() T {
	l := len(s.data)
	if l == 0 {
		fmt.Fprintln(os.Stderr, "Cannot Pop. Stack is already empty.")
		os.Exit(1)
	}
	top := s.data[l-1]
	s.data = s.data[:l-1]
	return top
}

func (s *Stack[T]) Peek() T {
	l := len(s.data)
	if l == 0 {
		fmt.Fprintln(os.Stderr, "Cannot Peek. Stack is already empty.")
		os.Exit(1)
	}
	return (s.data)[l-1]
}
