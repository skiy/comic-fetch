package command

// Controller interface
type Controller interface {
	// 新增漫画
	AddBook(siteURL string) (err error)

	// 获取数据
	ToFetch() (err error)

	// 获取章节 URL 列表
	ToFetchChapterList() (chapterURLList []string, err error)

	// 获取章节数据
	ToFetchChapter(chapterURL string) (chapterName string, imageURLList []string, err error)
}
