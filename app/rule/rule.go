package rule

// Rule interface
type Rule interface {
	// 获取章节 URL 列表
	ToFetchChapterList() (chapterURLList []string, err error)

	// 获取章节数据
	ToFetchChapter(chapterURL string) (chapterName string, imageURLList []string, err error)
}
