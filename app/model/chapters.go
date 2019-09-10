package model

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/skiy/comic-fetch/app/config"
)

// TbChapters 章节表
type TbChapters struct {
	ID        int64  `json:"id" gconv:"id"`                 // ID
	BookID    int64  `json:"book_id" gconv:"book_id"`       // 漫画 ID
	EpisodeID int    `json:"episode_id" gconv:"episode_id"` // 话序 ID
	Title     string `json:"title" gconv:"title"`           // 章节标题
	OrderID   int    `json:"order_id" gconv:"order_id"`     // 章节排序
	OriginID  int    `json:"origin_id" gconv:"origin_id"`   // 源章节ID
	Status    int    `json:"status" gconv:"status"`         // 状态 (0.采集成功, 1.采集失败, 2. 停止采集)
	OriginURL string `json:"origin_url" gconv:"origin_url"` // 采集地址
	CreatedAt int64  `json:"created_at" gconv:"created_at"` // 创建时间
	UpdatedAt int64  `json:"updated_at" gconv:"updated_at"` // 最后更新时间
}

// Chapters 漫画章节
type Chapters struct {
	model
}

// NewChapters Chapters init
func NewChapters() *Chapters {
	t := &Chapters{}
	t.connect()
	return t
}

// GetDataOne 获取一条信息
func (t *Chapters) GetDataOne(where interface{}) (device gdb.Record, err error) {
	return t.DB.Table(config.TbNameChapters).Where(where).One()
}

// AddData 添加一条信息
func (t *Chapters) AddData(data ...interface{}) (result sql.Result, err error) {
	return t.DB.Table(config.TbNameChapters).Data(data).Insert()
}

// GetData 获取一组数据
func (t *Chapters) GetData(where interface{}) (result gdb.Result, err error) {
	return t.DB.Table(config.TbNameChapters).Where(where).Select()
}

// UpdateData 获取一组数据
func (t *Chapters) UpdateData(data, where interface{}) (result sql.Result, err error) {
	return t.DB.Table(config.TbNameChapters).Data(data).Where(where).Update()
}
