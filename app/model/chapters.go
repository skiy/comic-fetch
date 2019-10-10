package model

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/skiy/comic-fetch/app/config/ctable"
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
	base
}

// NewChapters Chapters init
func NewChapters() *Chapters {
	t := &Chapters{}
	t.connect()
	return t
}

// GetDataOne 获取一条信息
func (t *Chapters) GetDataOne(where interface{}) (record gdb.Record, err error) {
	return t.getDataOne(ctable.TbNameChapters, where)
}

// AddData 添加一条信息
func (t *Chapters) AddData(data ...interface{}) (result sql.Result, err error) {
	return t.addData(ctable.TbNameChapters, data...)
}

// UpdateData 更新数据
func (t *Chapters) UpdateData(data, where interface{}) (result sql.Result, err error) {
	return t.updateData(ctable.TbNameChapters, data, where)
}

// DeleteData 删除数据
func (t *Chapters) DeleteData(where interface{}) (result sql.Result, err error) {
	return t.deleteData(ctable.TbNameChapters, where)
}

// GetData 获取一组数据
func (t *Chapters) GetData(where interface{}, sort string) (result gdb.Result, err error) {
	return t.getData(ctable.TbNameChapters, where, "")
}

// GetDataExt 获取一组数据 (扩展型)
func (t *Chapters) GetDataExt(params Params) (result gdb.Result, err error) {
	return t.getDataExt(ctable.TbNameChapters, params)
}
