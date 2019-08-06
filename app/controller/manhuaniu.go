package controller

import (
	"database/sql"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/library/fetch"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gf-utils/udb"
	"github.com/skiy/gf-utils/ulog"
)

// Manhuaniu 漫画牛
type Manhuaniu struct {
	Books  *model.TbBooks
	WebURL string
}

// NewManhuaniu Manhuaniu init
func NewManhuaniu(books *model.TbBooks) *Manhuaniu {
	t := &Manhuaniu{}
	t.Books = books
	return t
}

// ToFetch 采集
func (t *Manhuaniu) ToFetch() (err error) {
	web := webURL[t.Books.OriginFlag]
	if len(web) < t.Books.OriginWebType {
		return errors.New("runtime error: index out of range for origin_web_type")
	}

	t.WebURL = web[t.Books.OriginWebType]

	// 采集章节列表
	chapterURLList, err := t.ToFetchChapter()
	if err != nil {
		return err
	}

	if len(chapterURLList) == 0 {
		return errors.New("获取不到章节数据")
	}

	ulog.ReadLog().Println(chapterURLList)

	db := udb.GetDatabase()

	// 从数据库中获取已采集的章节列表
	chapters := ([]model.TbChapters)(nil)
	if err = db.Table(config.TbNameChapters).Structs(&chapters); err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = nil
	}
	ulog.ReadLog().Println(chapters)

	// 这里应该用 channel 并发获取章节数据
	for _, chapterURL := range chapterURLList {
		ulog.ReadLog().Println(t.WebURL + chapterURL)
	}

	return
}

// ToFetchChapter 采集章节 URL 列表
func (t *Manhuaniu) ToFetchChapter() (chapterURLList g.SliceStr, err error) {
	doc, err := fetch.PageSource(t.Books.OriginURL, "utf-8")
	if err != nil {
		return nil, err
	}

	doc.Find("#chapter-list-1 li a").Each(func(i int, aa *goquery.Selection) {
		chapterURL, exist := aa.Attr("href")
		if !exist {
			return
		}

		//chapterURLList = append(chapterURLList, t.WebURL + chapterURL)
		chapterURLList = append(chapterURLList, chapterURL)
	})

	return
}
