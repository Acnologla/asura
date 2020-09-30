package cache
import (
	"log"
	"os"
	"time"
	"fmt"
	"context"
	"encoding/json"
    "github.com/go-redis/redis/v8"

)

var client *redis.Client
func Init(){
	client = redis.NewClient(&redis.Options{
		Addr: "redis-14039.c91.us-east-1-3.ec2.cloud.redislabs.com:14039",
		Password: os.Getenv("REDIS"),
		DB: 0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err !=nil{
		log.Fatal(err)
	}
}

func Exists(val string) bool{
	value := client.Exists(context.Background(),val)
	return value.Val() ==  1
}

func Delete(val string){
	client.Del(context.Background(),val)
}

func Set(key string,value interface{},duration time.Duration){
	result, err := json.Marshal(value)
    if err != nil {
        fmt.Println(err)
	}
	client.Set(context.Background(),key,result,duration)
}

func Get(val string) []byte{
	value, _ := client.Get(context.Background(),val).Result()
	return []byte(value)
}