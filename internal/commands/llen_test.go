package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestLLenCommand(t *testing.T) {
	t.Run("LLEN should return 0 if the key does not exist", func(t *testing.T) {
		store := db.NewDB()
		command := LLenCommand{}
		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Integer(0) {
			t.Errorf("Expected 0, got %s", response)
		}
	})

	t.Run("LLEN should return the length of the list if the key exists", func(t *testing.T) {
		store := db.NewDB()
		exisistingList := &db.ListValue{ Items: db.Queue{} }
		exisistingList.Items.Push("value1")
		exisistingList.Items.Push("value2")
		exisistingList.Items.Push("value3")
		store.Set("key", exisistingList)
		command := LLenCommand{}

		response, err := command.Execute([]string{"key"}, store)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response != protocol.Integer(3) {
			t.Errorf("Expected 3, got %s", response)
		}
	})

	t.Run("LLEN should return an error if the key exists but is not a list", func(t *testing.T) {
		store := db.NewDB()
		store.Set("key", &db.StringValue{ Val: "value" })
		command := LLenCommand{}		
		
		_, err := command.Execute([]string{"key"}, store)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		
		if err.Error() != "wrong type for 'LLEN'" {
			t.Errorf("Expected 'wrong type for 'LLEN', got %s", err.Error())
		}
	})
}