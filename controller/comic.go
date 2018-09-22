package controller

import (
	"code.aliyun.com/skiystudy/comicFetch/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

type Init struct {
	Model model.Comic
	Cache *redis.Client
}

/**
抓取漫画初始化
*/
func (t *Init) Construct() {
	t.newBooks()
	t.getComicList()
}

/**
检测新书
*/
func (t *Init) newBooks() {
	comicList := t.Model.GetBookList(0)

	cacheKey := "newbooks"

	var v string
	var err error

	if t.Cache != nil {
		v, err = t.Cache.Do("get", cacheKey).String()
	} else {
		//file 方式
	}

	//空数据
	if v == "" {
		return
	}

	if err == nil {
		//fmt.Println(v)
		type m struct {
			Id   int    `json:"id"`
			Flag string `json:"flag"`
		}
		var m1 []m
		str := []byte(v)
		err = json.Unmarshal(str, &m1)

		if err == nil && m1 != nil {
			for _, v2 := range m1 {
				//fmt.Println(v2)
				for _, value := range comicList {
					//fmt.Println(value.OriginBookId, v2.Id, value.OriginBookId == v2.Id)
					if value.OriginBookId == v2.Id && value.OriginFlag == v2.Flag {
						goto Next
					}
				}

				t.addMh160Book(v2.Id)
				fmt.Printf("新增漫画ID(%d), 来源：<<%s>>\n", v2.Id, v2.Flag)
			Next:
			}

			t.Cache.Do("del", cacheKey)
		}
	}
}

/**
获取漫画列表
*/
func (t *Init) getComicList() {
	comicList := t.Model.GetBookList(0)
	//fmt.Println(comicList)

	comicTip := "\n漫画：%s (%d), %s\n"
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

		if strings.EqualFold(value.OriginFlag, "mh160") {
			var mh mh160
			mh.db = t.Model.Db
			mh.id = value.OriginBookId
			mh.originImageUrl = value.OriginImageUrl
			mh.Init()
		}
	}
}

/**
添加漫画160的漫画
*/
func (t *Init) addMh160Book(id int) {
	var mh mh160
	mh.db = t.Model.Db
	mh.id = id
	mh.Init()
}
