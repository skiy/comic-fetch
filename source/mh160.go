package source

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
	"fmt"
	"github.com/axgle/mahonia"
	regexp "regexp"
	"strconv"
)

type Mh160 struct {
}

func (this  *Mh160) Init() {
	this.chapter();
}

func (this *Mh160) chapter() {
	bookId := "31512"
	bookName := "万界仙踪"
	url := "https://m.mh160.com/kanmanhua/" + bookId + "/"

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

	doc, err := goquery.NewDocumentFromReader(rd)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".chapter-list ul li").Each(func(i int, s *goquery.Selection) {
		chapterName := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href");

		fmt.Printf("Review %d: - %s - %s \n", i, chapterName, url)

		this.detail(url, bookId, bookName, chapterName)
	})
}

func (this *Mh160) detail(url, bookId, bookName, chapterName string) {
	//preg := `2[0-9-\s:]*`
	preg := `/([0-9\/]*)/([0-9\.]*).html`
	re := regexp.MustCompile(preg)
	test := re.FindStringSubmatch(url)


	fmt.Println(test[1], test[2])

	imgUrl := "https://mhpic5.lineinfo.cn/mh160tuku/w/%s_%s/%s_%s/00%s.jpg"

	var fix string
	for i:=1; i<=10; i++ {
		if i < 10 {
			fix = "0" + strconv.Itoa(i)
		} else {
			fix = strconv.Itoa(i)
		}

		imgUrl2 := fmt.Sprintf(imgUrl, bookName, bookId, chapterName, test[2], fix)

		fmt.Println(imgUrl2)
		//break
	}

	//url := "https://www.mh160.com/kanmanhua/31512/658683.html"
	//https://mhpic5.lineinfo.cn/mh160tuku/w/%E4%B8%87%E7%95%8C%E4%BB%99%E8%B8%AA_31512/%E7%AC%AC97%E8%AF%9D_658683/0006.jpg
	//https://mhpic5.lineinfo.cn/mh160tuku/w/万界仙踪_31512/第97话_658683/0007.jpg


}