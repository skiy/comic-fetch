package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skiy/comicFetch/library"
	"github.com/skiy/comicFetch/model"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	Conf       library.Config
	comicModel model.Comic
)

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

	comicModel.Db = dbh

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", indexFunc)
	router.GET("/chapter/:id/:order", chapterFunc)
	router.GET("/picture/:id/:cid", pictureFunc)

	router.Static("/images", Conf.Image.Path)

	router.GET("/origin", originFunc)

	port := s.WebPort
	if port <= 0 {
		port = 35001
	}

	portStr := fmt.Sprintf(":%d", port)
	//router.Run()
	err = router.Run(portStr)
	if err != nil {
		router.Run(":65431")
	}
}

/**
获取漫画列表
*/
func indexFunc(c *gin.Context) {
	comicList := comicModel.GetBookList(0)

	if len(comicList) > 0 {
		var comicArr []model.TbBooks
		for _, book := range comicList {
			if len(book.ImageUrl) > 0 {
				book.ImageUrl = fmt.Sprintf("/images/%s", book.ImageUrl)
			} else {
				book.ImageUrl = fmt.Sprintf("/origin?url=%s&referer=%s", url.QueryEscape(book.OriginImageUrl), url.QueryEscape(book.OriginUrl))
			}

			comicArr = append(comicArr, book)
		}
		comicList = comicArr
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "第五城市漫画网",
		"list":  comicList,
	})
}

/**
获取漫画章节
*/
func chapterFunc(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		c.String(404, "漫画编号参数错误")
		if err != nil {
			log.Println(err)
		}
		return
	}

	order := c.Param("order")
	orderBy := "0"

	if order == "1" {
		orderBy = "1"
	}

	comicInfo := comicModel.GetBookList(id)
	if comicInfo == nil {
		c.String(404, "漫画不存在")
		return
	}

	title := fmt.Sprintf("%s - 第五城市漫画网", comicInfo[0].Name)

	chapterList := comicModel.GetChapterList(id)
	if len(chapterList) > 0 && orderBy == "0" {
		var chapterArr []model.TbChapter
		for i := len(chapterList) - 1; i > 0; i-- {
			chapterArr = append(chapterArr, chapterList[i])
		}
		chapterList = chapterArr
	}

	c.HTML(http.StatusOK, "chapter.html", gin.H{
		"title":       title,
		"info":        comicInfo[0],
		"chapterList": chapterList,
	})
}

/**
获取漫画图片
*/
func pictureFunc(c *gin.Context) {
	bidStr := c.Param("id")
	bid, err := strconv.Atoi(bidStr)

	if err != nil || bid <= 0 {
		c.String(404, "漫画编号参数错误")
		if err != nil {
			log.Println(err)
		}
		return
	}

	cidStr := c.Param("cid")
	cid, err := strconv.Atoi(cidStr)

	if err != nil || cid <= 0 {
		c.String(404, "章节编号参数错误")
		if err != nil {
			log.Println(err)
		}
		return
	}

	comicInfo := comicModel.GetBookList(bid)
	if comicInfo == nil {
		c.String(404, "漫画不存在")
		return
	}

	title := comicInfo[0].Name + " - 第五城市漫画网"

	imageList := comicModel.GetImages(bid, cid)

	if len(imageList) > 0 {
		var imageArr []model.TbImages
		for _, image := range imageList {
			if len(image.ImageUrl) > 0 {
				image.ImageUrl = fmt.Sprintf("/images/%s", image.ImageUrl)
			} else {
				image.ImageUrl = fmt.Sprintf("/origin?url=%s&referer=%s", url.QueryEscape(image.OriginUrl), url.QueryEscape(comicInfo[0].OriginUrl))
			}
			imageArr = append(imageArr, image)
		}
		imageList = imageArr
	}

	c.HTML(http.StatusOK, "picture.html", gin.H{
		"title": title,
		"list":  imageList,
	})
}

/**
远程图片
*/
func originFunc(c *gin.Context) {
	type Url struct {
		Url     string `form:"url"`
		Referer string `form:"referer"`
	}

	var url2 Url

	err := c.BindQuery(&url2)
	if err != nil {
		fmt.Println(err)
		return
	}

	url3, err := url.QueryUnescape(url2.Url)
	if err != nil {
		fmt.Println(err)
		return
	}

	referer, err := url.QueryUnescape(url2.Referer)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := library.OriginFile(url3, referer)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)
	newStr := buf.String()

	c.String(200, newStr)
}
