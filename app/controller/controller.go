package controller

import "github.com/gogf/gf/g"

// Controller interface
type Controller interface {
	// 获取数据
	ToFetch() (err error)
	// 获取章节 URL 列表
	ToFetchChapter() (chapterURLList g.SliceStr, err error)
}

const (
	pc = iota
	mobile
	api
)

var (
	webURL = map[string][]string{
		"manhuaniu": {
			pc:     "https://www.manhuaniu.com",
			mobile: "https://m.manhuaniu.com",
		},
	}
)
