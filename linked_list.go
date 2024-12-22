package main

import (
	"fmt"
	"os"
)

type ListNode[T any] struct {
	val  T
	prev *ListNode[T]
	next *ListNode[T]
}

type LinkedList[T any] struct {
	first *ListNode[T]
	last  *ListNode[T]
	n     int
}

func createLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{nil, nil, 0}
}

func createListNode[T any](v T) *ListNode[T] {
	n := ListNode[T]{v, nil, nil}
	return &n
}

func (l *LinkedList[T]) Add(node *ListNode[T]) {
	if node == nil {
		fmt.Fprintln(os.Stderr, "Cannot add a nil node to ", *l, ".")
		os.Exit(1)
	}
	if l.n == 0 {
		// List is empty
		l.first = node
		l.last = node
		l.n = 1
		return
	}
	l.last.next = node
	node.prev = l.last
	l.last = node
	l.n++
}

func (l *LinkedList[T]) Unlink(node *ListNode[T]) {
	if node == nil {
		fmt.Fprintln(os.Stderr, "Cannot unlink a nil node.")
		os.Exit(1)
	}
	if l.n == 0 {
		// List is empty
		fmt.Fprintln(os.Stderr, "Cannot unlink", *node, ". LinkedList ", *l, "is already empty.")
		os.Exit(1)
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if l.last == node {
		l.last = node.prev
	}
	if l.first == node {
		l.first = node.next
	}
	node.prev = nil
	node.next = nil
	l.n--
}

func (l *LinkedList[T]) Print() {
	fmt.Print("(n = ", l.n, ", f: ", l.first.val, ", l: ", l.last.val, "): ")
	for n := l.first; n != nil; n = n.next {
		fmt.Print("[", n.val, "]")
		if n.next != nil {
			fmt.Print("-")
		}
	}
	fmt.Println()
}
