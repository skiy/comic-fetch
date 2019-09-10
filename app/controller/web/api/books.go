package api

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/library/ldb"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/library/llog"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/comic-fetch/app/service/command"
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
	r.Response.Status = 500

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
		llog.Log.Warning(err.Error())
		response.Message = "漫画列表获取失败[Book.List]"
	} else {
		if err != sql.ErrNoRows {
			if err := resp.ToStructs(&books); err != nil {
				llog.Log.Warning(err.Error())
			}
		}

		response.Code = 0
		response.Message = "操作成功"
		response.Data = books
	}

	r.Response.WriteJson(response)
}

// Search 搜索
func (t *Book) Search(r *ghttp.Request) {
	r.Response.Status = 500

	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	// 漫画 名
	name := r.GetString("name")
	if name == "" {
		response.Message = "漫画名不存在"
		r.Response.WriteJson(response)
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
				llog.Log.Warning(err.Error())
			}
		}

		r.Response.Status = 200
		response.Code = 0
		response.Message = "操作成功"
		response.Data = books
	}

	r.Response.WriteJson(response)
}

// Add 添加新漫画
func (t *Book) Add(r *ghttp.Request) {
	r.Response.Status = 500

	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	type form struct {
		ID   int    `params:"id" gvalid:"id@required"`
		Site string `params:"site" gvalid:"site@required"`
	}

	formData := new(form)
	r.GetToStruct(formData)

	if err := gvalid.CheckStruct(formData, nil); err != nil {
		response.Message = err.FirstString()
		r.Response.WriteJson(response)
		return
	}

	if _, ok := config.WebURL[formData.Site]; ok {
		cliApp := command.NewCommand()

		if err := cliApp.Add(formData.Site, formData.ID); err != nil {
			response.Message = fmt.Sprintf("添加新漫画失败: %s", err.Error())
		} else {
			r.Response.Status = 200
			response.Code = 0
			response.Message = "添加新漫画成功"
		}
	} else {
		response.Message = fmt.Sprintf("不支持此网站 (%v) 添加新漫画", formData.Site)
	}

	r.Response.WriteJson(response)
}

// Update 更新漫画信息
func (t *Book) Update(r *ghttp.Request) {
	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"
	r.Response.Status = 500

	type form struct {
		ID     int `params:"id" gvalid:"id@required"`
		Status int `params:"status" gvalid:"status@in:0,1,2"`
	}

	formData := new(form)
	r.GetToStruct(formData)

	if err := r.GetToStruct(formData); err != nil {
		response.Message = err.Error()
		r.Response.WriteJson(response)
		return
	}

	if err := gvalid.CheckStruct(formData, nil); err != nil {
		response.Message = err.FirstString()
		r.Response.WriteJson(response)
		return
	}

	data := g.Map{
		"status": formData.Status,
	}

	_, err := ldb.GetDB().Table(config.TbNameBooks).Where(g.Map{"id": formData.ID}).Data(data).Update()
	if err != nil {
		llog.Log.Warningf(err.Error())
		response.Message = "漫画更新失败[Book.Update]"
	} else {
		r.Response.Status = 200
		response.Code = 0
		response.Message = "操作成功"
	}

	r.Response.WriteJson(response)
}
