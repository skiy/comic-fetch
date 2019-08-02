package fetch

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/g/encoding/gcharset"
	"github.com/gogf/gf/g/net/ghttp"
	"strings"
)

// PageSource 获取网页源码
//
// @sourceUrl string 请求网址
// @pageChartset string 网页字符集编码 (GBK,GB2312,UTF-8,GB18030)
// @return @doc DOM, @err error
func PageSource(sourceURL, pageChartset string) (doc *goquery.Document, err error) {
	c := ghttp.NewClient()
	response, err := c.Get(sourceURL)
	if err != nil {
		return nil, err
	}

	readCloser := response.Body
	//ulog.ReadLog().Println(2, response.ReadAllString())

	// 非 UTF-8 字符集则需要转换
	if strings.EqualFold(pageChartset, "UTF-8") {
		doc, err = goquery.NewDocumentFromReader(readCloser)
	} else {
		source, err := gcharset.Convert("UTF-8", pageChartset, response.ReadAllString())
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader([]byte(source))
		doc, err = goquery.NewDocumentFromReader(reader)
	}

	return
}
