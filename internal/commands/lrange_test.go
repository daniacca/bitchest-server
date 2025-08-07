package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestLRangeCommand(t *testing.T) {
	t.Run("LRANGE should return an empty array if the key does not exist", func(t *testing.T) {
		store := db.NewDB()
		command := LRangeCommand{}
		
		response, err := command.Execute([]string{"key", "0", "1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{}), response)
		}
	})

	t.Run("LRANGE should return an error if the key exists but is not a list", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.StringValue{ Val: "value" })
		command := LRangeCommand{}
		
		_, err := command.Execute([]string{"key", "0", "1"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		
		if err.Error() != "wrong type for 'LRANGE'" {
			t.Errorf("Expected 'wrong type for 'LRANGE', got %s", err.Error())
		}
	})

	t.Run("LRANGE should return an empty array if the start is greater than the length of the list", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LRangeCommand{}
		
		response, err := command.Execute([]string{"key", "3", "4"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{}), response)
		}
	})

	t.Run("LRANGE should return array with the last element if the stop is greater than the length of the list", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LRangeCommand{}

		response, err := command.Execute([]string{"key", "0", "4"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{"value1", "value2", "value3"}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{"value1", "value2", "value3"}), response)
		}
	})

	t.Run("LRANGE should return an empty array if the start is greater than the stop", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LRangeCommand{}

		response, err := command.Execute([]string{"key", "2", "1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{}), response)
		}
	})

	t.Run("LRANGE should return the correct elements if the start and stop are positive", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LRangeCommand{}

		response, err := command.Execute([]string{"key", "0", "2"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{"value1", "value2", "value3"}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{"value1", "value2", "value3"}), response)
		}
	})

	t.Run("LRANGE should return the correct elements if the start and stop are negative", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LRangeCommand{}

		response, err := command.Execute([]string{"key", "-2", "-1"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Array([]string{"value2", "value3"}) {
			t.Errorf("Expected %s, got %s", protocol.Array([]string{"value2", "value3"}), response)
		}
	})
}