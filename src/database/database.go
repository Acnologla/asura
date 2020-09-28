package database

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
	"github.com/andersfylling/disgord"
	"log"
	"errors"
	"os"
)

var Database *db.Client

type User struct {
	Avatars   []string `json:"avatar"`
	Usernames []string `json:"username"`
}

func Init() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config := &firebase.Config{
		DatabaseURL: fmt.Sprintf("https://%s.firebaseio.com/", os.Getenv("FIREBASE_PROJECT_ID")),
	}
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CONFIG")))
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
		return err
	}
	Database, err = app.Database(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}


func GetUserDB(id disgord.Snowflake) (User, error) {
	var acc User
	err := Database.NewRef(fmt.Sprintf("users/%d",id)).Get(context.Background(), &acc);
	if err != nil {	
		return acc, errors.New("Not bro")
	}
	return acc, nil
}