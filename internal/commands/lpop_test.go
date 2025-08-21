package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestLPopCommand(t *testing.T) {
	t.Run("LPOP should return a NIL reply if the key does not exist", func(t *testing.T) {
		store := db.NewDB()
		command := LPopCommand{}		
		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("LPOP should return a NIL reply if the list is empty", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		store.Set("key", exisistingList)
		command := LPopCommand{}

		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("LPOP should return a BULK reply if the list is not empty", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		store.Set("key", exisistingList)
		command := LPopCommand{}

		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Bulk("value1") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("value1"), response)
		}
	})

	t.Run("LPOP should return an array of elements if count is provided", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		store.Set("key", exisistingList)
		command := LPopCommand{}

		response, err := command.Execute([]string{"key", "2"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{"value1", "value2"}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{"value1", "value2"}), response)
		}
	})

	t.Run("LPOP should return an error if the key exists but is not a list", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.StringValue{ Val: "value" })
		command := LPopCommand{}		
		
		_, err := command.Execute([]string{"key"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		
		if err.Error() != "wrong type for 'LPOP'" {
			t.Errorf("Expected 'wrong type for 'LPOP', got %s", err.Error())
		}
	})

	t.Run("LPOP should return an error if no key is provided", func(t *testing.T) {
		store := db.NewDB()
		command := LPopCommand{}

		_, err := command.Execute([]string{}, store)
		if err == nil {
			t.Errorf("Expected error for no arguments, got nil")
		}

		if err.Error() != "wrong number of arguments for 'LPOP'" {
			t.Errorf("Expected 'wrong number of arguments for 'LPOP', got %s", err.Error())
		}
	})

	t.Run("LPOP should return an error if too many arguments are provided", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.ListValue{ Items: db.Queue{} })
		command := LPopCommand{}

		_, err := command.Execute([]string{"key", "10", "extra"}, store)
		if err == nil {
			t.Errorf("Expected error for too many arguments, got nil")
		}

		if err.Error() != "wrong number of arguments for 'LPOP'" {
			t.Errorf("Expected 'wrong number of arguments for 'LPOP', got %s", err.Error())
		}
	})

	t.Run("LPOP is a write command", func(t *testing.T) {
		cmd := &LPopCommand{}
		if !cmd.IsWrite() {
			t.Errorf("Expected IsWrite to return true, got false")
		}
	})
}