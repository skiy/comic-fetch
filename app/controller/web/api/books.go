package api

import (
	"database/sql"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/library/ldb"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/library/llog"
	"github.com/skiy/comic-fetch/app/model"
)

// Book Book
type Book struct {
	core
}

// NewBook Book init
func NewBook() *Book {
	t := &Book{}
	return t
}

// List 漫画列表
// sort=+id,-name 排序
// offset=10&limit=5 分页
// /api/books[/:id]
func (t *Book) List(r *ghttp.Request) {
	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	where := g.Map{}
	sort := "id desc"
	limit := 10

	// 漫画 ID
	if i := r.GetInt("id"); i != 0 {
		where["id"] = i
	}

	// 排序
	if s := r.GetQueryString("sort"); s != "" {
		if s2 := t.sort(s, nil); s2 != "" {
			sort = s2
		}
	}

	// 翻页
	offset := r.GetQueryInt("offset")

	// 每页显示数量
	if l := r.GetQueryInt("limit"); l > 0 {
		limit = l
	}

	books := ([]model.TbBooks)(nil)
	resp, err := ldb.GetDB().Table(config.TbNameBooks).Where(where).OrderBy(sort).Offset(offset).Limit(limit).Select()
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Debug(err.Error())
		response.Message = "漫画列表获取失败[Book.List]"
	} else {
		if err != sql.ErrNoRows {
			if err := resp.ToStructs(&books); err != nil {
				response.Message = err.Error()
			}
		}

		response.Code = 0
		response.Message = "操作成功"
		response.Data = books
	}

	if err := r.Response.WriteJson(response); err != nil {
		r.Response.Status = 500
	}
}

func (t *Book) Search(r *ghttp.Request) {
	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	// 漫画 ID
	name := r.GetString("name")
	if name == "" {
		response.Message = "漫画名不存在"
		if err := r.Response.WriteJson(response); err != nil {
			r.Response.Status = 500
		}
		return
	}

	like := "%" + name + "%"

	books := ([]model.TbBooks)(nil)
	resp, err := ldb.GetDB().Table(config.TbNameBooks).Where("name like ?", like).OrderBy("created_at desc").Select()

	if err != nil && err != sql.ErrNoRows {
		llog.Log.Debug(err.Error())
		response.Message = "漫画搜索失败[Book.Search]"
	} else {
		if err != sql.ErrNoRows {
			if err := resp.ToStructs(&books); err != nil {
				response.Message = err.Error()
			}
		}

		response.Code = 0
		response.Message = "操作成功"
		response.Data = books
	}

	if err := r.Response.WriteJson(response); err != nil {
		r.Response.Status = 500
	}
}
