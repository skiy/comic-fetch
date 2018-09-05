package source

import (
	"code.aliyun.com/skiystudy/comicFetch/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Mh160 struct {
	id    int
	url   string
	model model.Mh160
}

func (t *Mh160) Init(db *gorm.DB) {
	t.model.Db = db

	t.id = 31512
	t.url = "https://m.mh160.com" //手机版
	//t.url = "https://www.mh160.com" //PC版
	t.mobileChapter()
}

/**
获取源码
*/
func (t *Mh160) fetchSource(url string) (doc *goquery.Document) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	dec := mahonia.NewDecoder("gb18030")
	rd := dec.NewReader(res.Body)

	doc, err = goquery.NewDocumentFromReader(rd)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

/**
PC端获取章节
*/
/*
func (t *Mh160) pcChapter() {
	bookUrl := fmt.Sprintf("/kanmanhua/%d/", t.id)

	doc := t.fetchSource(t.url + bookUrl)

	var bookName string
	doc.Find(".intro_l").Each(func(i int, s *goquery.Selection) {
		bookName = s.Find(".title h1").Text()
		//fmt.Println(bookName)
	})

	doc.Find("#pic-list").Each(func(i int, s *goquery.Selection) {
		bookName = s.Find(".title h1").Text()
		fmt.Println(bookName)
	})

	//doc.Find(".chapter-list ul li").Each(func(i int, s *goquery.Selection) { //手机版
	doc.Find(".plist ul li").Each(func(i int, s *goquery.Selection) {
		chapterName := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")

		//fmt.Printf("Review %d: - %s - %s \n", i, chapterName, url)

		counts := t.countImage(url)
		t.detail(url, book.Id, chapter.Id, bookName, chapterName, counts)
	})
}
*/

/**
移动端获取章节
*/
func (t *Mh160) mobileChapter() {
	bookUrl := t.url + fmt.Sprintf("/kanmanhua/%d/", t.id)

	doc := t.fetchSource(bookUrl)

	nowTime := time.Now().Unix()

	var bookName string
	doc.Find(".main-bar").Each(func(i int, s *goquery.Selection) {
		bookName = s.Find("h1").Text()
		//fmt.Println(bookName)
	})

	book := t.model.Table.Books
	t.model.Db.Where("name = ?", bookName).First(&book)

	if book.Id == 0 {
		books := t.model.Table.Books
		//books.Id = t.id
		books.Name = bookName
		books.Status = 0
		books.OriginUrl = bookUrl
		books.OriginWeb = "漫画160"
		books.OriginBookId = t.id
		books.UpdatedAt = nowTime
		books.CreatedAt = nowTime

		book = t.model.CreateBook(books)
	}

	chapter_1 := t.model.Table.Chapter
	t.model.Db.Limit(1).Order("chapter_id desc").Find(&chapter_1)
	//fmt.Println(chapter_1)

	doc.Find(".chapter-list ul li").Each(func(i int, s *goquery.Selection) { //手机版
		chapterName := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")

		//fmt.Printf("Review %d: - %s - %s \n", i, chapterName, url)

		preg := `第([0-9]*)话`
		re := regexp.MustCompile(preg)
		test := re.FindStringSubmatch(chapterName)

		if len(test) < 2 {
			log.Fatalf("获取章节ID失败: %s %s", url, chapterName)
		}

		chapterNum, err := strconv.Atoi(test[1])
		if err != nil {
			log.Fatalf("章节转Int型失败: %s %s", test[1], chapterName)
		}

		//fmt.Println(chapterNum, chapter_1.ChapterId)
		if chapterNum > chapter_1.ChapterId {
			chapter := t.model.Table.Chapter
			chapter.Bid = book.Id
			chapter.ChapterId = chapterNum
			chapter.Title = chapterName
			chapter.OriginUrl = t.url + url
			chapter.CreatedAt = nowTime

			chapter = t.model.CreateChapter(chapter)

			counts := t.countImage(url)
			t.detail(url, book.Id, chapter.Id, chapterNum, bookName, chapterName, counts)
		}
	})
}

/**
获取共几话
*/
func (t *Mh160) countImage(url string) (counts int) {
	fetchUrl := t.url + url
	//fmt.Println(fetchUrl)

	var err error
	doc := t.fetchSource(fetchUrl)

	doc.Find(".main-bar").Each(func(i int, s *goquery.Selection) {
		imagePage := s.Find(".manga-page").Text()

		preg := `([0-9\/]*)/([0-9\/]*)P`
		re := regexp.MustCompile(preg)
		test := re.FindStringSubmatch(imagePage)

		if len(test) < 3 {
			log.Fatalf("获取章节话数失败: %s %s", url, imagePage)
		}

		counts, err = strconv.Atoi(test[2])
		if err != nil {
			log.Fatal(err)
		}
	})

	return counts
}

/**
获取漫画图片
*/
func (t *Mh160) detail(url string, bookId, chapterId, chapterNum int, bookName, chapterName string, counts int) {
	imgUrl := "https://mhpic5.lineinfo.cn/mh160tuku/w/%s_%d/%s_%s/"

	//preg := `2[0-9-\s:]*`
	preg := `/([0-9\/]*)/([0-9\.]*).html`
	re := regexp.MustCompile(preg)
	test := re.FindStringSubmatch(url)

	//fmt.Println(test[1], test[2])
	originChapterId := test[2]
	imgUrl1 := fmt.Sprintf(imgUrl, bookName, t.id, chapterName, originChapterId) + "00%s.jpg"
	//fmt.Println(imgUrl1)

	images := t.model.Table.Images
	images.Bid = bookId
	images.Cid = chapterId
	images.ChapterId = chapterNum
	images.ImageUrl = ""
	images.IsRemote = 1
	images.CreatedAt = time.Now().Unix()

	var fix string
	for i := 1; i < counts; i++ {
		//fmt.Println(i)
		if i < 10 {
			fix = "0" + strconv.Itoa(i)
		} else {
			fix = strconv.Itoa(i)
		}

		images.OriginUrl = fmt.Sprintf(imgUrl1, fix)

		t.model.CreateImages(images)

		//fmt.Println(image)
		//break
	}

	//url := "https://www.mh160.com/kanmanhua/31512/658683.html"
	//https://mhpic5.lineinfo.cn/mh160tuku/w/%E4%B8%87%E7%95%8C%E4%BB%99%E8%B8%AA_31512/%E7%AC%AC97%E8%AF%9D_658683/0006.jpg
	//https://mhpic5.lineinfo.cn/mh160tuku/w/万界仙踪_31512/第97话_658683/0007.jpg

}
