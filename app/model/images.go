package model

import (
	"database/sql"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/skiy/comic-fetch/app/config"
)

// TbImages 图片表
type TbImages struct {
	ID        int64  `gconv:"id"`         // ID
	BookID    int64  `gconv:"book_id"`    // 漫画 ID
	ChapterID int64  `gconv:"chapter_id"` // 章节 ID
	EpisodeID int    `gconv:"episode_id"` // 话序 ID
	ImageURL  string `gconv:"image_url"`  // 图片地址
	OriginURL string `gconv:"origin_url"` // 漫画图片采集地址
	Size      int64  `gconv:"size"`       // 文件大小
	OrderID   int    `gconv:"order_id"`   // 图片排序
	IsRemote  int    `gconv:"is_remote"`  // 是否远程图片
	CreatedAt int64  `gconv:"created_at"` // 创建时间
}

// Images 图片
type Images struct {
	model
}

// NewImages Images init
func NewImages() *Images {
	t := &Images{}
	t.connect()
	return t
}

// GetDataOne 获取一条信息
func (t *Images) GetDataOne(where interface{}) (device gdb.Record, err error) {
	return t.DB.Table(config.TbNameImages).Where(where).One()
}

// AddData 添加一条信息
func (t *Images) AddData(data ...interface{}) (result sql.Result, err error) {
	return t.DB.Table(config.TbNameImages).Data(data).Insert()
}

// GetData 获取一组数据
func (t *Images) GetData(where interface{}) (result gdb.Result, err error) {
	return t.DB.Table(config.TbNameImages).Where(where).Select()
}

// AddDataBatch 批量添加一组信息
func (t *Images) AddDataBatch(data []TbImages, batch int) (result sql.Result, err error) {
	if batch == 0 {
		batch = len(data)
	}
	return t.DB.Table(config.TbNameImages).Data(data).Batch(batch).Insert()
}
