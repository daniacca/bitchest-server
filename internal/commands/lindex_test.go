package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestLIndexCommand(t *testing.T) {
	t.Run("LINDEX should return a null bulk if the key does not exist", func(t *testing.T) {
		store := db.NewDB()
		command := LIndexCommand{}
		response, err := command.Execute([]string{"key", "0"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("LINDEX should return an error if the key exists but is not a list", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.StringValue{ Val: "value" })
		command := LIndexCommand{}
		_, err := command.Execute([]string{"key", "0"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("LINDEX should return an error if the index is not an integer", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.ListValue{ Items: db.Queue{} })
		command := LIndexCommand{}
		_, err := command.Execute([]string{"key", "not-an-integer"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("LINDEX should return a null bulk if the index is out of range", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.ListValue{ Items: db.Queue{} })
		command := LIndexCommand{}
		response, err := command.Execute([]string{"key", "1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("LINDEX should return the correct item if the index is in range", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LIndexCommand{}
		
		response, err := command.Execute([]string{"key", "0"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Bulk("value1") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("value1"), response)
		}

		response, err = command.Execute([]string{"key", "1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Bulk("value2") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("value2"), response)
		}

		response, err = command.Execute([]string{"key", "2"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Bulk("value3") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("value3"), response)
		}
	})

	t.Run("LINDEX should return the correct item if the index is negative", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LIndexCommand{}
		
		response, err := command.Execute([]string{"key", "-1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Bulk("value3") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("value3"), response)
		}
	})
}