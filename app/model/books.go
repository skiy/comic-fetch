package model

import (
	"database/sql"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/skiy/comic-fetch/app/config"
)

// TbBooks 漫画表
type TbBooks struct {
	ID             int64  `json:"id" gconv:"id"`                             // ID
	Name           string `json:"name" gconv:"name"`                         // 漫画名
	ImageURL       string `json:"image_url" gconv:"image_url"`               // 漫画图标地址
	Status         int    `json:"status" gconv:"status"`                     // 状态 (0正在更新,1暂停更新,2完结)
	OriginURL      string `json:"origin_url" gconv:"origin_url"`             // 漫画采集地址
	OriginWeb      string `json:"origin_web" gconv:"origin_web"`             // 源站名
	OriginWebType  string `json:"origin_web_type" gconv:"origin_web_type"`   // 采集源类型 (pc, mobile, api)
	OriginFlag     string `json:"origin_flag" gconv:"origin_flag"`           // 源站标识
	OriginImageURL string `json:"origin_image_url" gconv:"origin_image_url"` // 源站漫画图标地址
	OriginPathURL  string `json:"origin_path_url" gconv:"origin_path_url"`   // 上次采集图片实际路径
	OriginBookID   int    `json:"origin_book_id" gconv:"origin_book_id"`     // 本书ID
	UpdatedAt      int64  `json:"updated_at" gconv:"updated_at"`             // 更新时间
	CreatedAt      int64  `json:"created_at" gconv:"created_at"`             // 创建时间
}

// Books 漫画
type Books struct {
	model
}

// NewBooks Books init
func NewBooks() *Books {
	t := &Books{}
	t.connect()
	return t
}

// GetDataOne 获取一条信息
func (t *Books) GetDataOne(where interface{}) (device gdb.Record, err error) {
	return t.DB.Table(config.TbNameBooks).Where(where).One()
}

// AddData 添加一条信息
func (t *Books) AddData(data ...interface{}) (result sql.Result, err error) {
	return t.DB.Table(config.TbNameBooks).Data(data).Insert()
}

// GetData 获取一组数据
func (t *Books) GetData(where interface{}) (result gdb.Result, err error) {
	return t.DB.Table(config.TbNameBooks).Where(where).OrderBy("id DESC").Select()
}
