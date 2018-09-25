package library

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"log"
	"net/http"
)

/**
获取源码
*/
func FetchSource(url string) (doc *goquery.Document) {
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
