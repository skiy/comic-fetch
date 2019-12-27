package command

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/library/lfetch"
	"github.com/skiy/comic-fetch/app/library/lnotify"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/lcommon"
	"github.com/skiy/gfutils/lfilepath"
	"github.com/skiy/gfutils/llog"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type base struct {
	Books    *model.TbBooks
	WebURL   string
	ResURL   string
	Notified bool

	Prep prep
}

type prep struct {
	Book,
	SiteURL,
	ChapterList,
	Chapter,
	ChapterURL,
	ChapterPath,
	ImageStr,
	ImagesURL string
}

// AddBook Add new comic
func (t *base) AddBook(siteURL string) (err error) {
	t.WebURL = siteURL
	t.Books.OriginURL = fmt.Sprintf(t.Prep.SiteURL, siteURL, t.Books.OriginBookID)

	if err = t.ToFetchBook(); err != nil {
		return err
	}

	timeNow := time.Now().Unix()
	t.Books.UpdatedAt = timeNow
	t.Books.CreatedAt = timeNow

	bookModel := model.NewBooks()
	bookRes, err := bookModel.AddData(t.Books)
	if err != nil {
		return err
	}

	cfg := lcfg.GetCfg()
	notifyType := cfg.GetInt("notify.type")
	notifyNewBook := cfg.GetBool("notify.book")

	if notifyNewBook {
		notify := lnotify.NewNotifyMessage(1)

		// 钉钉通知
		if notifyType == 1 {
			if err := notify.Dingtalk(t.Books.Name, t.Books.OriginWeb, t.Books.OriginImageURL, t.Books.OriginURL); err != nil {
				llog.Log.Warningf("新增漫画通知失败: %v", err)
			}

			t.Notified = true
		}
	}

	t.Books.ID, _ = bookRes.LastInsertId()

	return t.ToFetch()
}

// ToFetch 采集
func (t *base) ToFetch() (err error) {
	log := llog.Log

	web, ok := config.WebURL[t.Books.OriginFlag]
	if ok {
		t.WebURL = web[t.Books.OriginWebType]
	} else {
		return errors.New("index out of range for origin_web_type: " + t.Books.OriginFlag)
	}

	if t.WebURL == "" {
		return errors.New("WebURL is nil: " + t.Books.OriginFlag)
	}

	// 采集章节列表
	chapterURLList, err := t.ToFetchChapterList()
	if err != nil {
		return err
	}

	if len(chapterURLList) == 0 {
		return errors.New("获取不到章节数据")
	}

	// 从数据库中获取已采集的章节列表
	chapterModel := model.NewChapters()
	chapterRes, err := chapterModel.GetData(g.Map{"book_id": t.Books.ID}, "")
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		err = nil
	}

	var chapters []model.TbChapters
	chapterStatusMap := map[int]model.TbChapters{}

	if chapterRes != nil {
		if err := chapterRes.ToStructs(&chapters); err != nil {
			return err
		}

		// 章节转 Map
		for _, chapter := range chapters {
			//log.Println(chapter)
			chapterStatusMap[chapter.OriginID] = chapter
		}
	}

	orderID := len(chapters) + 1
	cfg := lcfg.Get()

	imageLocal := cfg.GetBool("image.local")
	filePath := cfg.GetString("image.path")
	nametype := cfg.GetString("image.nametype")

	// 这里应该用 channel 并发获取章节数据
	for _, chapterURL := range chapterURLList {
		preg := t.Prep.ChapterURL
		re, err := regexp.Compile(preg)
		if err != nil {
			log.Warningf("章节 ID 正则执行失败: %v, URL: %s", err, chapterURL)
			continue
		}
		chapterIDs := re.FindStringSubmatch(chapterURL)
		if len(chapterIDs) != 2 {
			log.Warningf("章节 ID 提取失败: %v, URL: %s", err, chapterURL)
			continue
		}
		chapterOriginIDStr := chapterIDs[1]
		chapterOriginID, err := strconv.Atoi(chapterOriginIDStr)
		if err != nil {
			log.Warningf("章节 ID (%s) 转 Int 型失败: %v", chapterOriginIDStr, err)
			continue
		}

		// 章节是否存在
		chapterInfo, ok := chapterStatusMap[chapterOriginID]
		// 章节已存在
		if ok {
			// 章节采集已成功 或 章节停止采集
			if chapterInfo.Status == 0 || chapterInfo.Status == 2 {
				continue
			}
		}

		fullChapterURL := t.WebURL + chapterURL
		log.Debugf("\n\n███████████████████████████████████████████████████████████████████████████\n[URL] %s", fullChapterURL)

		chapterName, imageURLList, err := t.ToFetchChapter(fullChapterURL)
		if err != nil {
			log.Warningf("章节图片抓取失败: %v", err)
			continue
		}

		if len(imageURLList) == 0 {
			log.Warningf("章节名称: %s, 无图片资源", chapterName)
			continue
		}

		var episodeID int
		preg2 := `^([0-9]*)`
		re2 := regexp.MustCompile(preg2)
		episodeIDs := re2.FindStringSubmatch(chapterName)
		if len(episodeIDs) < 2 || episodeIDs[1] == "" {
			preg2 := `第([0-9]*)[话章]`
			re2 := regexp.MustCompile(preg2)
			episodeIDs = re2.FindStringSubmatch(chapterName)
		}

		if len(episodeIDs) > 1 {
			episodeID, _ = strconv.Atoi(strings.Trim(episodeIDs[1], ""))
		}

		log.Debugf("[Title] %s, [Image Count] %d", chapterName, len(imageURLList))

		status := 1 // 默认失败状态
		timestamp := time.Now().Unix()

		if !ok { // 新章节
			chapter := model.TbChapters{}
			chapter.BookID = t.Books.ID
			chapter.EpisodeID = episodeID
			chapter.Title = chapterName
			chapter.OrderID = orderID
			chapter.OriginID = chapterOriginID
			chapter.Status = status
			chapter.OriginURL = fullChapterURL
			chapter.CreatedAt = timestamp
			chapter.UpdatedAt = timestamp

			if res, err := chapterModel.AddData(chapter); err != nil {
				log.Warningf("新章节: %s, 保存失败", chapterName)
				continue
			} else {
				orderID++
				chapterInfo.ID, _ = res.LastInsertId()

				// 未通知过
				if !t.Notified {
					cfg := lcfg.Get()
					notifyType := cfg.GetInt("notify.type")
					notifyNewBook := cfg.GetBool("notify.book")

					if notifyNewBook {
						notify := lnotify.NewNotifyMessage(2)

						// 钉钉通知
						if notifyType == 1 {
							if err := notify.Dingtalk(t.Books.Name, chapter.Title, t.Books.OriginImageURL, chapter.OriginURL); err != nil {
								log.Warningf("更新漫画通知失败: %v", err)
							}

							t.Notified = true
						}
					}
				}
			}
		}

		imageModel := model.NewImages()

		var imageDataArr []model.TbImages
		for index, imageOriginURL := range imageURLList {
			imageOrderID := index + 1
			fullImageOriginURL := t.ResURL + imageOriginURL

			log.Debugf("[IMAGE URL] %s", fullImageOriginURL)

			var imageSize int64
			var imageURL string
			isRemote := 1

			// 图片本地化
			if imageLocal {
				if res, err := lfetch.GetResponse(fullImageOriginURL, fullChapterURL); err != nil {
					log.Warningf("远程获取图片本地化失败: %v", err)
				} else {
					fileName := fmt.Sprintf("%d-%d-%d", t.Books.ID, chapterInfo.ID, imageOrderID)
					if strings.EqualFold(nametype, "md5") {
						if name, err := gmd5.Encrypt(fileName); err != nil {
							log.Warningf("图片本地化文件名 MD5 加密失败: %v", err)
						} else {
							fileName = name
						}
					}

					fileExt := lfilepath.Ext(fullImageOriginURL)
					if fileExt == "" {
						fileExt = ".jpg"
					}

					// 真实保存的文件名
					fullFileName := fmt.Sprintf("%s/%s%s", filePath, fileName, fileExt)

					if imageFile, err := gfile.Create(fullFileName); err != nil {
						log.Warningf("本地化图片创建失败: %v", err)
					} else { // 创建成功
						imageSize, err = io.Copy(imageFile, res.Response.Body)
						if err != nil {
							log.Warningf("本地化图片保存失败: %v", err)
							if err := os.Remove(fullFileName); err != nil {
								log.Warningf("本地文件(%s)删除失败: %v", fullFileName, err)
							}
						} else { // 图片保存到本地成功
							isRemote = 0
							imageURL = fullFileName
						}
					}
				}
			}

			imageDataArr = append(imageDataArr, model.TbImages{
				ID:        0,
				BookID:    t.Books.ID,
				ChapterID: chapterInfo.ID,
				EpisodeID: episodeID,
				ImageURL:  imageURL,
				OriginURL: fullImageOriginURL,
				Size:      imageSize,
				OrderID:   imageOrderID,
				IsRemote:  isRemote,
				CreatedAt: timestamp,
			})
		}

		// break
		// 保存图片数据
		if _, err := imageModel.AddDataBatch(imageDataArr, 0); err != nil {
			log.Warningf("图片批量保存到数据库失败: %v", err)
		} else { // 图片保存成功
			chapterData := g.Map{
				"status":     0,
				"updated_at": time.Now().Unix(),
			}

			if _, err := chapterModel.UpdateData(chapterData, g.Map{"id": chapterInfo.ID}); err != nil {
				log.Warningf("章节(%d): %s, URL: %s, 状态(0)更新失败", chapterInfo.ID, chapterName, fullChapterURL)
			}
		}

		//break
	}

	return
}

// ToFetchBook 获取漫画信息
func (t *base) ToFetchBook() (err error) {
	doc, err := lfetch.PageSource(t.Books.OriginURL, "utf-8")
	if err != nil {
		return err
	}

	bookInfo := doc.Find(t.Prep.Book)

	if src, ok := bookInfo.Attr("src"); ok {
		if src != "" {
			t.Books.OriginImageURL = src
		}
	}

	if name, ok := bookInfo.Attr("alt"); ok {
		if name != "" {
			t.Books.Name = name
			return nil
		}
	}

	return errors.New("漫画标题获取失败")
}

// ToFetchChapterList 采集章节 URL 列表
func (t *base) ToFetchChapterList() (chapterURLList g.SliceStr, err error) {
	doc, err := lfetch.PageSource(t.Books.OriginURL, "utf-8")
	if err != nil {
		return nil, err
	}

	doc.Find(t.Prep.ChapterList).Each(func(i int, aa *goquery.Selection) {
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
func (t *base) ToFetchChapter(chapterURL string) (chapterName string, imageURLList g.SliceStr, err error) {
	doc, err := lfetch.PageSource(chapterURL, "utf-8")
	if err != nil {
		return
	}

	var reg10, reg20, reg30, reg31 *regexp.Regexp

	scriptNodes := doc.Find("script")

	scriptCount := scriptNodes.Length()
	scriptNameText := scriptNodes.Eq(scriptCount - 1).Text()

	reg10, err = regexp.Compile(t.Prep.Chapter)
	if err != nil {
		return
	}
	infos := reg10.FindStringSubmatch(scriptNameText)

	if len(infos) == 2 {
		infoStr := strings.ReplaceAll(infos[1], `"`, "")
		infoArr := strings.Split(infoStr, ",")

		if len(infoArr) == 4 {
			chapterName = infoArr[1]
		}
	}

	llog.Log.Debugf("正在抓取章节 (%s) 图片", chapterName)

	// 采集图片
	chapterPath := ""
	imageStr := ""
	for i := 0; i <= scriptCount; i++ {
		scriptImageText := scriptNodes.Eq(i).Text()

		reg20, err = regexp.Compile(t.Prep.ChapterPath)
		if err != nil {
			continue
		}

		path := reg20.FindStringSubmatch(scriptImageText)
		if len(path) == 2 {
			chapterPath = path[1]
		}

		reg30, err = regexp.Compile(t.Prep.ImageStr)
		if err != nil {
			continue
		}

		imageStrs := reg30.FindStringSubmatch(scriptImageText)
		if len(imageStrs) == 2 {
			imageStr = imageStrs[1]
		}

		if imageStr == "" {
			err = errors.New("图片列表字符串获取失败")
			continue
		}

		reg31, err = regexp.Compile(t.Prep.ImagesURL)
		if err != nil {
			llog.Log.Debugf("script-%d 未采到图片: %s", i, err.Error())
			continue
		}

		imagesList := reg31.FindAllStringSubmatch(imageStr, -1)
		if imagesList == nil {
			continue
		}

		for _, images := range imagesList {
			if len(images) == 2 {
				imageURLList = append(imageURLList, "/"+strings.TrimLeft(lcommon.Stripslashes(chapterPath+images[1]), "/"))
			}
		}

		if len(imageURLList) > 0 {
			break
		}
	}

	return
}
