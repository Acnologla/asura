package test

import (
	"testing"
	"asura/src/database"
	"context"
) 


func TestDatabaseConnect(t *testing.T){
	err := database.Init()
	if err !=nil{
		t.Error(err)
	}
	t.Run("TestDatabaseQuery",func(t *testing.T){
		var acc database.User
		if err := database.Database.NewRef("users/395542409959964672").Get(context.Background(), &acc); err != nil {
		  t.Error(err)
		}
	})
}
