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
}

type base struct {
	db         gdb.DB
	rowsStyle, // 获取数量的风格 (0,不获取数量, 1.仅获取数量, 2.数量和数据都获取)
	count int
}

func (t *base) connect() {
	t.db = g.DB()
}

// GetDataOne 获取一条信息
func (t *base) getDataOne(name string, where interface{}) (record gdb.Record, err error) {
	return t.db.Table(name).Where(where).One()
}

// AddData 添加一条信息
func (t *base) addData(name string, data ...interface{}) (result sql.Result, err error) {
	return t.db.Table(name).Data(data).Insert()
}

// UpdateData 更新数据
func (t *base) updateData(name string, data, where interface{}) (result sql.Result, err error) {
	return t.db.Table(name).Where(where).Data(data).Update()
}

// DeleteData 删除数据
func (t *base) deleteData(name string, where interface{}) (result sql.Result, err error) {
	return t.db.Table(name).Where(where).Delete()
}

// 获取数据
func (t *base) getData(name string, where interface{}, sort string) (result gdb.Result, err error) {
	if sort == "" {
		sort = "id DESC"
	}
	m := t.db.Table(name).Where(where)

	// 获取数量
	if t.rowsStyle != 0 {
		t.count, err = m.Count()

		// 重置获取数量为否
		t.rowsStyle = 0

		// 仅获取数量
		if t.rowsStyle == 1 {
			return
		}
	}

	return m.OrderBy(sort).Select()
}

// GetDataExt 获取数据扩展方式
func (t *base) getDataExt(name string, params Params) (result gdb.Result, err error) {
	m := t.db.Table(name).Where(params.Where)

	// 获取数量
	if t.rowsStyle != 0 {
		t.count, err = m.Count()

		// 重置获取数量为否
		t.rowsStyle = 0

		// 仅获取数量
		if t.rowsStyle == 1 {
			return
		}
	}

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
func (t *base) addDataBatch(name string, data []interface{}, batch int) (result sql.Result, err error) {
	if batch == 0 {
		batch = len(data)
	}
	return t.db.Table(name).Data(data).Batch(batch).Insert()
}

// Rows 是否获取数量
func (t *base) Rows(flag int) *base {
	t.rowsStyle = flag
	return t
}

// Count 获取数据总数
func (t *base) Count() int {
	return t.count
}
