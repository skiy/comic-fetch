package model

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/skiy/comic-fetch/app/config/ctable"
)

// TbImages 图片表
type TbImages struct {
	ID        int64  `json:"id" gconv:"id"`                 // ID
	BookID    int64  `json:"book_id" gconv:"book_id"`       // 漫画 ID
	ChapterID int64  `json:"chapter_id" gconv:"chapter_id"` // 章节 ID
	EpisodeID int    `json:"episode_id" gconv:"episode_id"` // 话序 ID
	ImageURL  string `json:"image_url" gconv:"image_url"`   // 图片地址
	OriginURL string `json:"origin_url" gconv:"origin_url"` // 漫画图片采集地址
	Size      int64  `json:"size" gconv:"size"`             // 文件大小
	OrderID   int    `json:"order_id" gconv:"order_id"`     // 图片排序
	IsRemote  int    `json:"is_remote" gconv:"is_remote"`   // 是否远程图片
	CreatedAt int64  `json:"created_at" gconv:"created_at"` // 创建时间
}

// Images 图片
type Images struct {
	base
}

// NewImages Images init
func NewImages() *Images {
	t := &Images{}
	t.connect()
	return t
}

// GetDataOne 获取一条信息
func (t *Images) GetDataOne(where interface{}) (record gdb.Record, err error) {
	return t.getDataOne(ctable.TbNameImages, where)
}

// AddData 添加一条信息
func (t *Images) AddData(data ...interface{}) (result sql.Result, err error) {
	return t.addData(ctable.TbNameImages, data...)
}

// UpdateData 更新数据
func (t *Images) UpdateData(data, where interface{}) (result sql.Result, err error) {
	return t.updateData(ctable.TbNameImages, data, where)
}

// DeleteData 删除数据
func (t *Images) DeleteData(where interface{}) (result sql.Result, err error) {
	return t.deleteData(ctable.TbNameImages, where)
}

// GetData 获取一组数据
func (t *Images) GetData(where interface{}) (result gdb.Result, err error) {
	return t.getData(ctable.TbNameImages, where, "id DESC")
}

// GetDataExt 获取一组数据 (扩展型)
func (t *Images) GetDataExt(params Params) (result gdb.Result, err error) {
	return t.getDataExt(ctable.TbNameImages, params)
}

// AddDataBatch 批量添加一组信息
func (t *Images) AddDataBatch(data []TbImages, batch int) (result sql.Result, err error) {
	return t.addDataBatch(ctable.TbNameImages, data, batch)
}
