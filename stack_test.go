package main

import (
	"testing"
)

func TestPeekPop(t *testing.T) {

	var s Stack

	if !s.IsEmpty() {
		t.Errorf("Stack should be empty")
	}

	s = s.Push(1)

	peeked := s.Peek()
	if peeked != 1 {
		t.Errorf("Peek() = %d; want 1", peeked)
	}

	if s.IsEmpty() {
		t.Errorf("Stack shouldn't be empty")
	}

	s, peeked = s.Pop()
	if peeked != 1 {
		t.Errorf("Pop() = %d; want 1", peeked)
	}

	if !s.IsEmpty() {
		t.Errorf("Stack should be empty")
	}
}
