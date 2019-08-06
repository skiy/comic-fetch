package model

// TbChapters 章节表
type TbChapters struct {
	ID        int    `gconv:"id"`         // ID
	BookID    int    `gconv:"book_id"`    // 漫画 ID
	ChapterID int    `gconv:"chapter_id"` // 章节 ID
	Title     string `gconv:"title"`      // 章节标题
	OrderID   int    `gconv:"order_id"`   // 章节排序
	OriginID  int    `gconv:"origin_id"`  // 源章节ID
	OriginURL string `gconv:"origin_url"` // 采集地址
	CreatedAt int64  `gconv:"created_at"` // 创建时间
}

// Chapters 漫画章节
type Chapters struct {
}
