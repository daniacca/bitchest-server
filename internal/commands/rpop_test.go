package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestRPopCommand(t *testing.T) {
	t.Run("RPOP should return a NIL reply if the key does not exist", func(t *testing.T) {
		store := db.NewDB()
		command := RPopCommand{}		
		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("RPOP should return a NIL reply if the list is empty", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		store.Set("key", exisistingList)
		command := RPopCommand{}

		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.NullBulk() {
			t.Errorf("Expected %s, got %s", protocol.NullBulk(), response)
		}
	})

	t.Run("RPOP should return a BULK reply if the list is not empty", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		store.Set("key", exisistingList)
		command := RPopCommand{}

		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Bulk("value1") {
			t.Errorf("Expected %s, got %s", protocol.Bulk("value1"), response)
		}
	})

	t.Run("RPOP should return an array of elements if count is provided", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		store.Set("key", exisistingList)
		command := RPopCommand{}

		response, err := command.Execute([]string{"key", "2"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{"value2", "value1"}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{"value2", "value1"}), response)
		}
	})
}