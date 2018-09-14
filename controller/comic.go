package controller

import (
	"code.aliyun.com/skiystudy/comicFetch/model"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

type Init struct {
	Model model.Comic
	Redis redis.Client
}

func (t *Init) Construct() {
	t.getComicList()
	//t.addBook(11106, "")
}

func (t *Init) getComicList() {
	comicList := t.Model.GetBookList(0)
	fmt.Println(comicList)

	comicTip := "漫画：%s (%d), %s\n"
	for _, value := range comicList {
		//fmt.Println(index, value)

		if value.Status != 0 {
			s := "暂停更新"
			if value.Status == 2 {
				s = "完结"
			}

			fmt.Printf(comicTip, value.Name, value.Id, s)
			continue
		}

		fmt.Printf(comicTip, value.Name, value.Id, "正在更新……")

		if strings.EqualFold(value.OriginWeb, "漫画160") {
			var mh mh160
			mh.db = t.Model.Db
			mh.id = value.OriginBookId
			mh.originImageUrl = value.OriginImageUrl
			mh.Init()
		}
	}
}

func (t *Init) addBook(id int, source string) {
	var mh mh160
	mh.db = t.Model.Db
	mh.id = id
	mh.Init()
}
