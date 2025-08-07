package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestLRemCommand_Execute(t *testing.T) {
	store := db.NewDB()
	key := "mylist"
	list := &db.ListValue{}
	list.Items.Push("a")
	list.Items.Push("b")
	list.Items.Push("a")
	list.Items.Push("c")
	list.Items.Push("a")
	store.Set(key, list)

	cmd := &LRemCommand{}

	t.Run("removes all occurrences when count is 0", func(t *testing.T) {
		// Reset list
		list := &db.ListValue{}
		list.Items.Push("a")
		list.Items.Push("b")
		list.Items.Push("a")
		list.Items.Push("c")
		list.Items.Push("a")
		store.Set(key, list)

		resp, err := cmd.Execute([]string{key, "0", "a"}, store)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != protocol.Integer(3) {
			t.Errorf("expected 3 removed, got %s", resp)
		}
		got := list.Items.GetItems()
		expected := []string{"b", "c"}
		if len(got) != len(expected) {
			t.Errorf("expected %d items, got %d", len(expected), len(got))
		}
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("expected %s at %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("removes up to count from head to tail when count > 0", func(t *testing.T) {
		// Reset list
		list := &db.ListValue{}
		list.Items.Push("a")
		list.Items.Push("b")
		list.Items.Push("a")
		list.Items.Push("c")
		list.Items.Push("a")
		store.Set(key, list)

		resp, err := cmd.Execute([]string{key, "2", "a"}, store)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != protocol.Integer(2) {
			t.Errorf("expected 2 removed, got %s", resp)
		}
		got := list.Items.GetItems()
		expected := []string{"b", "c", "a"}
		if len(got) != len(expected) {
			t.Errorf("expected %d items, got %d", len(expected), len(got))
		}
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("expected %s at %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("removes up to count from tail to head when count < 0", func(t *testing.T) {
		// Reset list
		list := &db.ListValue{}
		list.Items.Push("a")
		list.Items.Push("b")
		list.Items.Push("a")
		list.Items.Push("c")
		list.Items.Push("a")
		store.Set(key, list)

		resp, err := cmd.Execute([]string{key, "-2", "a"}, store)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != protocol.Integer(2) {
			t.Errorf("expected 2 removed, got %s", resp)
		}
		got := list.Items.GetItems()
		expected := []string{"a", "b", "c"}
		if len(got) != len(expected) {
			t.Errorf("expected %d items, got %d", len(expected), len(got))
		}
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("expected %s at %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("does nothing if value is not present", func(t *testing.T) {
		// Reset list
		list := &db.ListValue{}
		list.Items.Push("a")
		list.Items.Push("b")
		list.Items.Push("c")
		store.Set(key, list)

		resp, err := cmd.Execute([]string{key, "0", "x"}, store)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != protocol.Integer(0) {
			t.Errorf("expected 0 removed, got %s", resp)
		}
		got := list.Items.GetItems()
		expected := []string{"a", "b", "c"}
		if len(got) != len(expected) {
			t.Errorf("expected %d items, got %d", len(expected), len(got))
		}
		for i, v := range expected {
			if got[i] != v {
				t.Errorf("expected %s at %d, got %s", v, i, got[i])
			}
		}
	})

	t.Run("returns 0 if key does not exist", func(t *testing.T) {
		resp, err := cmd.Execute([]string{"nokey", "0", "a"}, store)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != protocol.Integer(0) {
			t.Errorf("expected 0 removed, got %s", resp)
		}
	})

	t.Run("returns error for wrong type", func(t *testing.T) {
		store.Set("notalist", &db.StringValue{ Val: "value" })
		_, err := cmd.Execute([]string{"notalist", "0", "a"}, store)
		if err == nil {
			t.Errorf("expected error for wrong type")
		}
	})

	t.Run("returns error for invalid count", func(t *testing.T) {
		_, err := cmd.Execute([]string{key, "notanint", "a"}, store)
		if err == nil {
			t.Errorf("expected error for invalid count")
		}
	})

	t.Run("returns error for wrong number of arguments", func(t *testing.T) {
		_, err := cmd.Execute([]string{key, "1"}, store)
		if err == nil {
			t.Errorf("expected error for wrong number of arguments")
		}
	})
}
