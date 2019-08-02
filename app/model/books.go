package model

// TbBooks 漫画表
type TbBooks struct {
	ID             int    `gconv:"id"`               // ID
	Name           string `gconv:"name"`             // 漫画名
	ImageURL       string `gconv:"image_url"`        // 漫画图标地址
	Status         int    `gconv:"status"`           // 状态 (0正在更新,1暂停更新,2完结)
	OriginURL      string `gconv:"origin_url"`       // 漫画采集地址
	OriginWeb      string `gconv:"origin_web"`       // 源站名
	OriginWebType  int    `gconv:"origin_web_type"`  // 采集源类型 (0.pc, 1.mobile, 3.api)
	OriginFlag     string `gconv:"origin_flag"`      // 源站标识
	OriginImageURL string `gconv:"origin_image_url"` // 源站漫画图标地址
	OriginPathURL  string `gconv:"origin_path_url"`  // 上次采集图片实际路径
	OriginBookID   int    `gconv:"origin_book_id"`   // 本书ID
	UpdatedAt      int64  `gconv:"updated_at"`       // 更新时间
	CreatedAt      int64  `gconv:"created_at"`       // 创建时间
}

// Books 漫画
type Books struct {
}
