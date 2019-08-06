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
	"strconv"
	"strings"
	"time"
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
	log := ulog.ReadLog()
	log.Printf("\n正在采集漫画: %s\n源站: %s\n源站漫画URL: %s\n", t.Books.Name, t.Books.OriginWeb, t.Books.OriginURL)

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

	// 章节转Map
	chapterStatusMap := map[int]model.TbChapters{}
	for _, chapter := range chapters {
		//log.Println(chapter)
		chapterStatusMap[chapter.OriginID] = chapter
	}

	orderID := len(chapters)
	//log.Println(chapters, "orderID: ", orderID)
	log.Println("orderID: ", orderID)

	// 这里应该用 channel 并发获取章节数据
	for _, chapterURL := range chapterURLList {
		preg := `\/([0-9]*).html`
		re, err := regexp.Compile(preg)
		if err != nil {
			log.Warningf("章节ID正则执行失败: %v, URL: %s", err, chapterURL)
			continue
		}
		chapterIDs := re.FindStringSubmatch(chapterURL)
		if len(chapterIDs) != 2 {
			log.Warningf("章节ID提取失败: %v, URL: %s", err, chapterURL)
			continue
		}
		chapterIDStr := chapterIDs[1]
		chapterID, err := strconv.Atoi(chapterIDStr)
		if err != nil {
			log.Fatalf("章节ID(%s)转Int型失败: %v", chapterIDStr, err)
			continue
		}

		// 章节是否存在
		chapterInfo, ok := chapterStatusMap[chapterID]
		// 章节已存在
		if ok {
			// 章节采集已成功 或 章节停止采集
			if chapterInfo.Status == 0 || chapterInfo.Status == 2 {
				continue
			}
		}

		fullChapterURL := t.WebURL + chapterURL
		log.Println(fullChapterURL, chapterID)

		chapterName, imageURLList, err := t.ToFetchChapter(fullChapterURL)
		if err != nil {
			log.Warningf("章节: %s, 图片抓取失败: %v", fullChapterURL, err)
			continue
		}

		if len(imageURLList) == 0 {
			log.Warningf("章节: %s, URL: %s, 无图片资源", chapterName, fullChapterURL)
			continue
		}

		var episodeID int
		preg2 := `第([0-9]*)[话章]`
		re2 := regexp.MustCompile(preg2)
		episodeIDs := re2.FindStringSubmatch(chapterName)

		if len(episodeIDs) > 1 {
			episodeID, _ = strconv.Atoi(strings.Trim(episodeIDs[1], ""))
		}

		log.Println(chapterName, len(imageURLList), episodeID)
		// 保存章节

		for _, imageURL := range imageURLList {
			log.Println(t.ResURL + imageURL)
		}

		status := 0
		timestamp := time.Now().Unix()

		// 存在章节(原来采集失败的章节)则变更采集状态
		if ok {
			chapterInfo.Status = status
			chapterInfo.UpdatedAt = timestamp

			if _, err := db.Table(config.TbNameChapters).Data(chapterInfo).Where(g.Map{"id": chapterInfo.ID}).Update(); err != nil {
				log.Warningf("章节: %s, URL: %s, 更新失败", chapterName, fullChapterURL)
			}
		} else { // 未存在章节, 则新增章节
			chapter := model.TbChapters{}
			chapter.BookID = t.Books.ID
			chapter.EpisodeID = episodeID
			chapter.Title = chapterName
			chapter.OrderID = orderID
			chapter.OriginID = chapterID
			chapter.Status = status
			chapter.OriginURL = fullChapterURL
			chapter.CreatedAt = timestamp
			chapter.UpdatedAt = timestamp

			if _, err := db.Table(config.TbNameChapters).Data(chapter).Insert(); err != nil {
				log.Warningf("章节: %s, URL: %s, 保存失败", chapterName, fullChapterURL)
			} else {
				orderID++
			}
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
	re, err := regexp.Compile(pregImages)
	if err != nil {
		return "", nil, err
	}
	images := re.FindAllString(script2Text, -1)

	if images == nil {
		return
	}

	for _, image := range images {
		imageURLList = append(imageURLList, strings.ReplaceAll(image, "\\", ""))
	}

	script22Text := doc.Find("script").Eq(22).Text()

	pregInfo := `SinMH\.initChapter\(([^;]*)\)`
	re2, err := regexp.Compile(pregInfo)
	if err != nil {
		return "", nil, err
	}
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
