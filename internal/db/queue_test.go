package db

import (
	"testing"
)

func TestQueue(t *testing.T) {
	t.Run("Push should add items to the end of the queue", func(t *testing.T) {
		queue := Queue{}
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetSize() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetSize())
		}

		if queue.GetItems()[0] != "item1" {
			t.Errorf("Expected item1, got %s", queue.GetItems()[0])
		}
	})

	t.Run("Pop should remove items from the end of the queue", func(t *testing.T) {
		queue := Queue{}
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetSize() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetSize())
		}

		item := queue.Pop()
		if item != "item2" {
			t.Errorf("Expected item2, got %s", item)
		}

		if queue.GetSize() != 1 {
			t.Errorf("Expected size to be 1, got %d", queue.GetSize())
		}
	})

	t.Run("Shift should remove items from the beginning of the queue", func(t *testing.T) {
		queue := Queue{}
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetSize() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetSize())
		}

		item := queue.Shift()
		if item != "item1" {
			t.Errorf("Expected item1, got %s", item)
		}

		if queue.GetSize() != 1 {
			t.Errorf("Expected size to be 1, got %d", queue.GetSize())
		}
	})

	t.Run("Unshift should add items to the beginning of the queue", func(t *testing.T) {
		queue := Queue{}
		queue.Push("item1")
		queue.Push("item2")

		if queue.GetSize() != 2 {
			t.Errorf("Expected size to be 2, got %d", queue.GetSize())
		}

		queue.Unshift("item0")

		if queue.GetSize() != 3 {
			t.Errorf("Expected size to be 3, got %d", queue.GetSize())
		}

		if queue.GetItems()[0] != "item0" {
			t.Errorf("Expected item0, got %s", queue.GetItems()[0])
		}

		if queue.GetItems()[1] != "item1" {
			t.Errorf("Expected item1, got %s", queue.GetItems()[1])
		}
	})
}
