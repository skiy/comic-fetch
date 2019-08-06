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
	"regexp"
	"strings"
)

// Manhuaniu 漫画牛
type Manhuaniu struct {
	Books  *model.TbBooks
	WebURL string
	ResURL string
}

// NewManhuaniu Manhuaniu init
func NewManhuaniu(books *model.TbBooks) *Manhuaniu {
	t := &Manhuaniu{}
	t.Books = books
	t.ResURL = "https://res.nbhbzl.com/"
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
	chapterURLList, err := t.ToFetchChapterList()
	if err != nil {
		return err
	}

	if len(chapterURLList) == 0 {
		return errors.New("获取不到章节数据")
	}

	log := ulog.ReadLog()
	log.Println(chapterURLList)

	db := udb.GetDatabase()

	// 从数据库中获取已采集的章节列表
	chapters := ([]model.TbChapters)(nil)
	if err = db.Table(config.TbNameChapters).Structs(&chapters); err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = nil
	}
	log.Println(chapters)

	// 这里应该用 channel 并发获取章节数据
	for _, chapterURL := range chapterURLList {
		fullChapterURL := t.WebURL + chapterURL
		log.Println(fullChapterURL)

		chapterName, imageURLList, err := t.ToFetchChapter(fullChapterURL)
		if err != nil {
			log.Warningf("章节: %s, 图片抓取失败: %v", fullChapterURL, err)
			continue
		}

		if len(imageURLList) == 0 {
			log.Warningf("章节: %s, URL: %s, 无图片资源", chapterName, fullChapterURL)
			continue
		}

		log.Println(chapterName, len(imageURLList))

		for _, imageURL := range imageURLList {
			log.Println(t.ResURL + imageURL)
		}
		//break
	}

	return
}

// ToFetchChapterList 采集章节 URL 列表
func (t *Manhuaniu) ToFetchChapterList() (chapterURLList g.SliceStr, err error) {
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

// ToFetchChapter 获取章节内容
func (t *Manhuaniu) ToFetchChapter(chapterURL string) (chapterName string, imageURLList g.SliceStr, err error) {
	doc, err := fetch.PageSource(chapterURL, "utf-8")
	if err != nil {
		return
	}

	script2Text := doc.Find("script").Eq(2).Text()

	pregImages := `images\\/[^"]*`
	re, _ := regexp.Compile(pregImages)
	images := re.FindAllString(script2Text, -1)

	if images == nil {
		return
	}

	for _, image := range images {
		imageURLList = append(imageURLList, strings.ReplaceAll(image, "\\", ""))
	}

	script22Text := doc.Find("script").Eq(22).Text()

	pregInfo := `SinMH\.initChapter\(([^;]*)\)`
	re2, _ := regexp.Compile(pregInfo)
	infos := re2.FindStringSubmatch(script22Text)

	if len(infos) == 2 {
		infoStr := strings.ReplaceAll(infos[1], `"`, "")
		infoArr := strings.Split(infoStr, ",")

		if len(infoArr) == 4 {
			chapterName = infoArr[1]
		}
	}

	return
}
