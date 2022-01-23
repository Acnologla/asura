package main

import (
	"asura/src/adapter"
	"asura/src/database"
	"asura/src/entities"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

func Process(id disgord.Snowflake) {
	x := time.Now()
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 500; j++ {
				err := adapter.UserAdapterImpl.UpdateUser(id, func(u entities.User) entities.User {
					u.Pity++
					if u.Pity >= 5000 {
						fmt.Println(time.Since(x).Seconds())
					}
					return u
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
}

func main() {

	godotenv.Load()
	database.Connect(database.GetEnvConfig())
	adapter.Init()
	go Process(12)
	go Process(10)
	fmt.Println("executing")
	for {
	}
}
