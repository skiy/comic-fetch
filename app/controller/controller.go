package controller

// Controller interface
type Controller interface {
	// 获取章节 URL 列表
	ToFetchChapter() (err error)
}
