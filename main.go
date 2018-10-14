package main

import (
	"github.com/skiy/comicFetch/controller"
	"github.com/skiy/comicFetch/library"
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
		db.Init(Conf.Mysql.Host, Conf.Mysql.User, Conf.Mysql.Password, Conf.Mysql.Name, Conf.Mysql.Char)
	} else if s.Datatype == "sqlite" {
		db.Init("", "", "", Conf.Sqlite.Name, "")
	}

	dbh, err := db.Connect()
	defer dbh.Close()

	if err != nil {
		log.Fatalln("Db connect fail", err)
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
	Comic.Conf = Conf
	Comic.Construct()
}
