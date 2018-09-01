package source

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
	"fmt"
	"github.com/axgle/mahonia"
)

type Mh160 struct {

}

func (this *Mh160) Init() {
	url := "https://www.mh160.com/kanmanhua/31512/658683.html"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	dec := mahonia.NewDecoder("gb18030")
	fmt.Println(res)
	rd := dec.NewReader(res.Body)

	doc, err := goquery.NewDocumentFromReader(rd)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".img-box").Each(func(i int, s *goquery.Selection) {
		title := s.Find("p").Text()

		//title,_ := s.Find("img").Attr("src")

		fmt.Printf("Review %d: - %s\n", i, title)
	})
}

func (this *Mh160) Init2() {
	url := "https://m.mh160.com/kanmanhua/31512/"

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
		title := s.Find("a").Text()

		fmt.Printf("Review %d: - %s\n", i, title)
	})
}

func (this *Mh160) Test() {
	url := "https://www.mh160.com/kanmanhua/31512/"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".plist ul li").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text()

		fmt.Printf("Review %d: %s\n", i, title)
	})
}