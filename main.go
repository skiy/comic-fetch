package main

import (
	"code.aliyun.com/skiystudy/comicFetch/controller"
	"code.aliyun.com/skiystudy/comicFetch/library"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

func main() {
	db := new(library.Database)
	db.Init("localhost", "root", "123456", "comic", "utf8")
	dbh, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	Comic := new(controller.Init)
	Comic.Model.Db = dbh
	Comic.Construct()

}
