package source

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
	"fmt"
	"github.com/axgle/mahonia"
	"regexp"
	"strconv"
)

type Mh160 struct {
	id  int
	url string
}

func (this *Mh160) Init() {
	this.id = 31512
	this.url = "https://m.mh160.com" //手机版
	//this.url = "https://www.mh160.com" //PC版
	this.mobileChapter()
}

/**
	获取源码
 */
func (this *Mh160) fetchSource(url string) (doc *goquery.Document) {
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
func (this *Mh160) pcChapter() {
	bookUrl := fmt.Sprintf("/kanmanhua/%d/", this.id)

	doc := this.fetchSource(this.url + bookUrl)

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

		counts := this.countImage(url)
		this.detail(url, bookName, chapterName, counts)
	})
}

/**
	移动端获取章节
 */
func (this *Mh160) mobileChapter() {
	bookUrl := fmt.Sprintf("/kanmanhua/%d/", this.id)

	doc := this.fetchSource(this.url + bookUrl)

	var bookName string
	doc.Find(".main-bar").Each(func(i int, s *goquery.Selection) {
		bookName = s.Find("h1").Text()
		//fmt.Println(bookName)
	})

	doc.Find(".chapter-list ul li").Each(func(i int, s *goquery.Selection) { //手机版
		chapterName := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")

		//fmt.Printf("Review %d: - %s - %s \n", i, chapterName, url)

		counts := this.countImage(url)
		this.detail(url, bookName, chapterName, counts)

		/*
		if i == 0 {
			log.Fatalln(1)
		}
		*/
	})
}

/**
	获取共几话
 */
func (this *Mh160) countImage(url string) (counts int) {
	fetchUrl := this.url + url
	//fmt.Println(fetchUrl)

	var err error
	doc := this.fetchSource(fetchUrl)

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
func (this *Mh160) detail(url, bookName, chapterName string, counts int) {
	imgUrl := "https://mhpic5.lineinfo.cn/mh160tuku/w/%s_%d/%s_%s/"

	//preg := `2[0-9-\s:]*`
	preg := `/([0-9\/]*)/([0-9\.]*).html`
	re := regexp.MustCompile(preg)
	test := re.FindStringSubmatch(url)

	//fmt.Println(test[1], test[2])
	chapterId := test[2]
	imgUrl1 := fmt.Sprintf(imgUrl, bookName, this.id, chapterName, chapterId) + "00%s.jpg"
	//fmt.Println(imgUrl1)

	var fix string
	for i := 1; i < counts; i++ {
		//fmt.Println(i)
		if i < 10 {
			fix = "0" + strconv.Itoa(i)
		} else {
			fix = strconv.Itoa(i)
		}

		imgUrl2 := fmt.Sprintf(imgUrl1, fix)
		//imgUrl2 := fmt.Sprintf(imgUrl, bookName, this.id, chapterName, test[2], fix)

		fmt.Println(imgUrl2)
		//break
	}

	//url := "https://www.mh160.com/kanmanhua/31512/658683.html"
	//https://mhpic5.lineinfo.cn/mh160tuku/w/%E4%B8%87%E7%95%8C%E4%BB%99%E8%B8%AA_31512/%E7%AC%AC97%E8%AF%9D_658683/0006.jpg
	//https://mhpic5.lineinfo.cn/mh160tuku/w/万界仙踪_31512/第97话_658683/0007.jpg

}
