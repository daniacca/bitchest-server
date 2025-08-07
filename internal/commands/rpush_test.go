package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestRPushCommand(t *testing.T) {
	t.Run("RPush should return error if no arguments are provided", func(t *testing.T) {
		store := db.NewDB()
		cmd := &RPushCommand{}

		_, err := cmd.Execute([]string{}, store)
		if err == nil {
			t.Error("Expected error for wrong number of arguments, got nil")
		}
		
		if err.Error() != "wrong number of arguments for 'RPUSH'" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("RPush should return error if the key is not a list", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.StringValue{ Val: "value" })
		cmd := &RPushCommand{}
		
		_, err := cmd.Execute([]string{"key", "value1"}, store)
		if err == nil {
			t.Error("Expected error for wrong type, got nil")
		}
		
		if err.Error() != "wrong type for 'RPUSH'" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("RPush should create a new list if the key doesn't exist", func(t *testing.T) {
		store := db.NewDB()
		cmd := &RPushCommand{}
		
		result, err := cmd.Execute([]string{"key", "value1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != ":1\r\n" {
			t.Errorf("Expected 1, got %s", result)
		}
		
		val, ok := store.Get("key")
		if !ok {
			t.Error("Expected key to exist, got nil")
		}
		
		list, ok := val.(*db.ListValue)
		if !ok {
			t.Error("Expected list value, got nil")
		}

		if list.Items.GetLength() != 1 {
			t.Errorf("Expected 1 items, got %d", list.Items.GetLength())
		}
		
		if list.Items.GetItems()[0] != "value1" {
			t.Errorf("Expected value1, got %s", list.Items.GetItems()[0])
		}
	})

	t.Run("RPush should add values to an existing list", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		store.Set("key", exisistingList)

		cmd := &RPushCommand{}
		result, err := cmd.Execute([]string{"key", "value2"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != ":2\r\n" {
			t.Errorf("Expected 2, got %s", result)
		}

		val, ok := store.Get("key")
		if !ok {
			t.Error("Expected key to exist, got nil")
		}

		list, ok := val.(*db.ListValue)
		if !ok {
			t.Error("Expected list value, got nil")
		}

		if list.Items.GetLength() != 2 {
			t.Errorf("Expected 2 items, got %d", list.Items.GetLength())
		}
		
		if list.Items.GetItems()[0] != "value1" {
			t.Errorf("Expected value1, got %s", list.Items.GetItems()[0])
		}
		
		if list.Items.GetItems()[1] != "value2" {
			t.Errorf("Expected value2, got %s", list.Items.GetItems()[1])
		}
	})

	t.Run("RPush should add all values to the list", func(t *testing.T) {
		store := db.NewDB()
		cmd := &RPushCommand{}
		result, err := cmd.Execute([]string{"key", "value1", "value2", "value3"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != ":3\r\n" {
			t.Errorf("Expected 3, got %s", result)
		}

		val, ok := store.Get("key")
		if !ok {
			t.Error("Expected key to exist, got nil")
		}	

		list, ok := val.(*db.ListValue)
		if !ok {
			t.Error("Expected list value, got nil")
		}

		if list.Items.GetLength() != 3 {
			t.Errorf("Expected 3 items, got %d", list.Items.GetLength())
		}

		if list.Items.GetItems()[0] != "value1" {
			t.Errorf("Expected value1, got %s", list.Items.GetItems()[0])
		}

		if list.Items.GetItems()[1] != "value2" {
			t.Errorf("Expected value2, got %s", list.Items.GetItems()[1])
		}

		if list.Items.GetItems()[2] != "value3" {
			t.Errorf("Expected value3, got %s", list.Items.GetItems()[2])
		}
	})
}