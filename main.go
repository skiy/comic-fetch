package main

import (
	"code.aliyun.com/skiystudy/comicFetch/controller"
	"code.aliyun.com/skiystudy/comicFetch/library"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var Conf library.Config

func main() {

	Conf.ReadConfig()

	s := Conf.Setting

	db := new(library.Database)
	db.Datatype = s.Datatype

	if s.Datatype == "mysql" {
		fmt.Println(Conf.Mysql)
		db.Init(Conf.Mysql.Host, Conf.Mysql.User, Conf.Mysql.Password, Conf.Mysql.Name, Conf.Mysql.Char)
	} else if s.Datatype == "sqlite" {
		db.Init("", "", "", Conf.Sqlite.Name, "")
	}

	dbh, err := db.Connect()
	defer dbh.Close()

	if err != nil {
		log.Fatalln(err)
	}

	Comic := new(controller.Init)

	cache := redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Host + ":" + Conf.Redis.Port,
		Password: Conf.Redis.Password, // no password set
		DB:       Conf.Redis.Db,       // use default DB
	})

	pong, err := cache.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
	} else {
		Comic.Cache = cache
	}

	Comic.Model.Db = dbh
	Comic.Construct()
}
