package controller

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/library/fetch"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gf-utils/ulog"
)

// Manhuaniu 漫画牛
type Manhuaniu struct {
	Books *model.TbBooks
}

// NewManhuaniu Manhuaniu init
func NewManhuaniu(books *model.TbBooks) *Manhuaniu {
	t := &Manhuaniu{}
	t.Books = books
	return t
}

// ToFetchChapter 采集章节 URL 列表
func (t *Manhuaniu) ToFetchChapter() (err error) {
	doc, err := fetch.PageSource(t.Books.OriginURL, "utf-8")
	if err != nil {
		return err
	}

	var chapterURLList g.SliceStr
	doc.Find("#chapter-list-1 li a").Each(func(i int, aa *goquery.Selection) {
		chapterURL, exist := aa.Attr("href")
		if !exist {
			return
		}

		chapterURLList = append(chapterURLList, chapterURL)
	})
	ulog.ReadLog().Println(chapterURLList)

	return
}
