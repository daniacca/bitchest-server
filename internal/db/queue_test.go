package db

import (
	"testing"
)

func TestQueue(t *testing.T) {
	t.Run("Push should add items to the end of the queue", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetLength() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetLength())
		}

		if queue.GetItems()[0] != "item1" {
			t.Errorf("Expected item1, got %s", queue.GetItems()[0])
		}
	})

	t.Run("Pop should remove items from the end of the queue", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetLength() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetLength())
		}

		item, err := queue.Pop()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if item != "item2" {
			t.Errorf("Expected item2, got %s", item)
		}

		if queue.GetLength() != 1 {
			t.Errorf("Expected size to be 1, got %d", queue.GetLength())
		}
	})

	t.Run("Pop should return an error if the queue is empty", func(t *testing.T) {
		queue := NewQueue()
		_, err := queue.Pop()
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "queue is empty" {
			t.Errorf("Expected error message 'queue is empty', got %v", err)
		}
	})

	t.Run("Shift should return an error if the queue is empty", func(t *testing.T) {
		queue := NewQueue()
		_, err := queue.Shift()
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Shift should remove items from the beginning of the queue", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetLength() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetLength())
		}

		item, err := queue.Shift()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if item != "item1" {
			t.Errorf("Expected item1, got %s", item)
		}

		if queue.GetLength() != 1 {
			t.Errorf("Expected size to be 1, got %d", queue.GetLength())
		}
	})

	t.Run("Unshift should add items to the beginning of the queue", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetLength() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetLength())
		}	

		queue.Unshift("item0")

		if queue.GetLength() != 3 {
			t.Errorf("Expected size to be 3, got %d", queue.GetLength())
		}

		if queue.GetItems()[0] != "item0" {
			t.Errorf("Expected item0, got %s", queue.GetItems()[0])
		}

		if queue.GetItems()[1] != "item1" {
			t.Errorf("Expected item1, got %s", queue.GetItems()[1])
		}
	})
}
