package main

import (
	"code.aliyun.com/skiystudy/comicFetch/controller"
	"code.aliyun.com/skiystudy/comicFetch/library"
	"log"
)

func main() {
	db := new(library.Database)
	db.Init("localhost", "root", "123456", "comic", "utf8")
	dbh, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	new(controller.Init).Construct(dbh)
}
