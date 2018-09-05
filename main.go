package main

import (
	"code.aliyun.com/skiystudy/comicFetch/library"
	"code.aliyun.com/skiystudy/comicFetch/source"
	"log"
)

func main() {
	db := new(library.Database)
	db.Init("localhost", "root", "123456", "comic", "utf8")
	dbh, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	comic := new(source.Mh160)
	comic.Init(dbh)
}
