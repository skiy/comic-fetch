package model

// TbImages 图片表
type TbImages struct {
	ID        int    `gconv:"id"`         // ID
	BookID    int    `gconv:"book_id"`    // 漫画 ID
	Cid       int    `gconv:"cid"`        // 章节编号
	ChapterID int    `gconv:"chapter_id"` // 章节 ID
	ImageURL  string `gconv:"image_url"`  // 图片地址
	OriginURL string `gconv:"origin_url"` // 漫画图片采集地址
	Size      int64  `gconv:"size"`       // 文件大小
	OrderID   int    `gconv:"order_id"`   // 图片排序
	IsRemote  int    `gconv:"is_remote"`  // 是否远程图片
	CreatedAt int64  `gconv:"created_at"` // 创建时间
}

// Images 图片
type Images struct {
}
