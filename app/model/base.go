package model

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// Params 查询参数
type Params struct {
	Where   interface{} // 条件
	Sort    string      // 排序
	Offset, // 位移
	Limit int // 查询数量
	Rows bool // 查询总数
}

type base struct {
	DB gdb.DB
}

func (t *base) connect() {
	t.DB = g.DB()
}

// GetDataOne 获取一条信息
func (t *base) getDataOne(name string, where interface{}) (record gdb.Record, err error) {
	return t.DB.Table(name).Where(where).One()
}

// AddData 添加一条信息
func (t *base) addData(name string, data ...interface{}) (result sql.Result, err error) {
	return t.DB.Table(name).Data(data).Insert()
}

// UpdateData 更新数据
func (t *base) updateData(name string, data, where interface{}) (result sql.Result, err error) {
	return t.DB.Table(name).Where(where).Data(data).Update()
}

// DeleteData 删除数据
func (t *base) deleteData(name string, where interface{}) (result sql.Result, err error) {
	return t.DB.Table(name).Where(where).Delete()
}

// GetDataExt 获取数据扩展方式
func (t *base) getDataExt(name string, params Params) (result gdb.Result, err error) {
	m := t.DB.Table(name).Where(params.Where)

	if params.Sort != "" {
		m = m.OrderBy(params.Sort)
	}

	if params.Limit > 0 {
		m = m.Offset(params.Offset).Limit(params.Limit)
	}

	result, err = m.Select()
	return
}

// AddDataBatch 批量添加一组信息
func (t *base) addDataBatch(name string, data []TbImages, batch int) (result sql.Result, err error) {
	if batch == 0 {
		batch = len(data)
	}
	return t.DB.Table(name).Data(data).Batch(batch).Insert()
}
