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

	t.Run("Index should return an error if the index is out of range", func(t *testing.T) {
		queue := NewQueue()
		_, err := queue.Index(0)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Index should return the correct item if the index is in range", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("item1")
		queue.Push("item2")
		queue.Push("item3")

		item, err := queue.Index(1)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if item != "item2" {
			t.Errorf("Expected item2, got %s", item)
		}
	})

	t.Run("Set should return an error if the index is out of range", func(t *testing.T) {
		queue := NewQueue()
		err := queue.Set(0, "item1")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		queue.Push("item1")
		err = queue.Set(1, "item2")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		err = queue.Set(-1, "item2")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("Set should set the correct item if the index is in range", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("item1")
		queue.Push("item2")
		queue.Push("item3")

		err := queue.Set(1, "item4")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if queue.GetItems()[1] != "item4" {
			t.Errorf("Expected item4, got %s", queue.GetItems()[1])
		}
	})

	t.Run("Remove should remove all occurrences when count is 0", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("a")
		queue.Push("b")
		queue.Push("a")
		queue.Push("c")
		queue.Push("a")
		removed := queue.Remove("a", 0)
		if removed != 3 {
			t.Errorf("Expected 3 items removed, got %d", removed)
		}
		expected := []string{"b", "c"}
		got := queue.GetItems()
		if len(got) != len(expected) {
			t.Errorf("Expected queue length %d, got %d", len(expected), len(got))
		}
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("Remove should remove up to count occurrences from head to tail when count > 0", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("a")
		queue.Push("b")
		queue.Push("a")
		queue.Push("c")
		queue.Push("a")
		
		removed := queue.Remove("a", 2)
		if removed != 2 {
			t.Errorf("Expected 2 items removed, got %d", removed)
		}
		
		expected := []string{"b", "c", "a"}
		got := queue.GetItems()
		if len(got) != len(expected) {
			t.Errorf("Expected queue length %d, got %d", len(expected), len(got))
		}

		for i, v := range expected {
			if got[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("Remove should remove up to count occurrences from tail to head when count < 0", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("a")
		queue.Push("b")
		queue.Push("a")
		queue.Push("c")
		queue.Push("a")
		
		removed := queue.Remove("a", -2)
		if removed != 2 {
			t.Errorf("Expected 2 items removed, got %d", removed)
		}
		
		expected := []string{"a", "b", "c"}
		got := queue.GetItems()
		if len(got) != len(expected) {
			t.Errorf("Expected queue length %d, got %d", len(expected), len(got))
		}
		
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("Remove should do nothing if value is not present", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("a")
		queue.Push("b")
		queue.Push("c")
		
		removed := queue.Remove("x", 0)
		if removed != 0 {
			t.Errorf("Expected 0 items removed, got %d", removed)
		}
		
		expected := []string{"a", "b", "c"}
		got := queue.GetItems()
		
		if len(got) != len(expected) {
			t.Errorf("Expected queue length %d, got %d", len(expected), len(got))
		}
		
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("Remove should handle empty queue", func(t *testing.T) {
		queue := NewQueue()
		removed := queue.Remove("a", 0)
		if removed != 0 {
			t.Errorf("Expected 0 items removed, got %d", removed)
		}
		if queue.GetLength() != 0 {
			t.Errorf("Expected queue length 0, got %d", queue.GetLength())
		}
	})

	t.Run("Remove should remove all occurrence when count is 0", func(t *testing.T) {
		queue := NewQueue()
		queue.Push("a")
		queue.Push("b")
		queue.Push("a")
		queue.Push("c")
		queue.Push("a")
		queue.Push("d")
		queue.Push("a")

		removed := queue.Remove("a", 0)
		if removed != 4 {
			t.Errorf("Expected 4 items removed, got %d", removed)
		}
		
		expected := []string{"b", "c", "d"}
		got := queue.GetItems()
		if len(got) != len(expected) {
			t.Errorf("Expected queue length %d, got %d", len(expected), len(got))
		}

		for i, v := range expected {
			if got[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, got[i])
			}
		}
	})
}
