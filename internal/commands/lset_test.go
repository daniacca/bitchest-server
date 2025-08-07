package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestLSetCommand(t *testing.T) {
	t.Run("LSET should return a null bulk if the key does not exist", func(t *testing.T) {
		store := db.NewDB()
		command := LSetCommand{}
		response, err := command.Execute([]string{"key", "0", "value"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("LSET should return an error if the key exists but is not a list", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.StringValue{Val: "value"})
		command := LSetCommand{}
		_, err := command.Execute([]string{"key", "0", "value"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("LSET should return an error if the index is not an integer", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.ListValue{Items: db.Queue{}})
		command := LSetCommand{}
		_, err := command.Execute([]string{"key", "not-an-integer", "value"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("LSET should return an error if the wrong number of arguments is given", func(t *testing.T) {
		store := db.NewDB()
		command := LSetCommand{}
		_, err := command.Execute([]string{"key", "0"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		_, err = command.Execute([]string{"key", "0", "value", "extra"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("LSET should return an error if the index is out of range", func(t *testing.T) {
		store := db.NewDB()
		list := &db.ListValue{Items: db.Queue{}}
		list.Items.Push("a")
		store.Set("key", list)
		command := LSetCommand{}
		_, err := command.Execute([]string{"key", "1", "b"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		_, err = command.Execute([]string{"key", "-2", "b"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("LSET should set the correct item if the index is in range", func(t *testing.T) {
		store := db.NewDB()
		list := &db.ListValue{Items: db.Queue{}}
		list.Items.Push("a")
		list.Items.Push("b")
		list.Items.Push("c")
		store.Set("key", list)
		command := LSetCommand{}
		response, err := command.Execute([]string{"key", "1", "x"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response != protocol.Bulk("OK") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("OK"), response)
		}
		// Check that the value was set
		val, _ := store.Get("key")
		lv := val.(*db.ListValue)
		if lv.Items.GetItems()[1] != "x" {
			t.Errorf("Expected item at index 1 to be 'x', got %s", lv.Items.GetItems()[1])
		}
	})

	t.Run("LSET should set the correct item if the index is negative", func(t *testing.T) {
		store := db.NewDB()
		list := &db.ListValue{Items: db.Queue{}}
		list.Items.Push("a")
		list.Items.Push("b")
		list.Items.Push("c")
		store.Set("key", list)
		command := LSetCommand{}
		response, err := command.Execute([]string{"key", "-1", "z"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if response != protocol.Bulk("OK") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("OK"), response)
		}
		val, _ := store.Get("key")
		lv := val.(*db.ListValue)
		if lv.Items.GetItems()[2] != "z" {
			t.Errorf("Expected item at index 2 to be 'z', got %s", lv.Items.GetItems()[2])
		}
	})
}
