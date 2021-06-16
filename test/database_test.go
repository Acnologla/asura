package test

import (
	"asura/src/database"
	"context"
	"testing"
)

func TestDatabaseConnect(t *testing.T) {
	err := database.Init()
	if err != nil {
		t.Error(err)
	}
	t.Run("TestDatabaseQuery", func(t *testing.T) {
		var acc database.User
		if err := database.Database.NewRef("users/395542409959964672").Get(context.Background(), &acc); err != nil {
			t.Error(err)
		}
	})
	t.Run("TestDatabaseIsBanned", func(t *testing.T) {
		if database.IsBanned(0) {
			t.Errorf("This must return false")
		}
	})
}
