package api

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/config/cerror"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/library/llog"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/comic-fetch/app/service/command"
	"time"
)

// Book Book
type Book struct {
	Base
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
	response := g.Map{}

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

	p := model.Params{
		Where:  where,
		Sort:   sort,
		Offset: offset,
		Limit:  limit,
	}

	books := ([]model.TbBooks)(nil)
	m := model.NewBooks()
	resp, err := m.GetDataExt(p)
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Warning(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrGetData, g.Map{"message": "漫画列表获取失败[Book.List]"})

	} else {
		if err != sql.ErrNoRows {
			if resp != nil {
				if err := resp.ToStructs(&books); err != nil {
					llog.Log.Warning(err.Error())
				}
			}
		}
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess, g.Map{"data": books})
	}

	_ = r.Response.WriteJson(response)
}

// Search 搜索
func (t *Book) Search(r *ghttp.Request) {
	response := g.Map{}

	// 漫画 名
	name := r.GetString("name")
	if name == "" {
		r.Response.Status, response = lfunc.Response(cerror.ErrBookNameNotExist)
		_ = r.Response.WriteJson(response)
		return
	}

	like1 := "name like ?"
	like2 := "%" + name + "%"

	p := model.Params{
		Where: g.Map{
			like1: like2,
		},
		Sort: "created_at desc",
	}

	books := ([]model.TbBooks)(nil)
	bookModel := model.NewBooks()
	resp, err := bookModel.GetDataExt(p)
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Warning(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrGetData, g.Map{"message": "漫画搜索失败[Book.Search]"})
	} else {
		if err != sql.ErrNoRows && resp != nil {
			if err := resp.ToStructs(&books); err != nil {
				llog.Log.Warning(err.Error())
			}
		}
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess, g.Map{"data": books})
	}

	_ = r.Response.WriteJson(response)
}

// Add 添加新漫画
func (t *Book) Add(r *ghttp.Request) {
	response := g.Map{}

	type form struct {
		ID   int    `params:"id" gvalid:"id@required"`
		Site string `params:"site" gvalid:"site@required"`
	}

	formData := new(form)
	if err := r.GetToStruct(formData); err != nil {
		r.Response.Status, response = lfunc.Response(cerror.ErrInvalidParameter)
		_ = r.Response.WriteJson(response)
		return
	}

	if err := gvalid.CheckStruct(formData, nil); err != nil {
		llog.Log.Warning(err.String())
		r.Response.Status, response = lfunc.Response(cerror.ErrInvalidParameter)
		_ = r.Response.WriteJson(response)
		return
	}

	if _, ok := config.WebURL[formData.Site]; ok {
		cliApp := command.NewCommand()

		if err := cliApp.Add(formData.Site, formData.ID); err != nil {
			r.Response.Status, response = lfunc.Response(cerror.ErrAddData, g.Map{"message": "添加新漫画失败[Book.Add]"})
		} else {
			r.Response.Status, response = lfunc.Response(cerror.ErrAddSuccess)
		}
	} else {
		param := g.Map{
			"message": fmt.Sprintf("不支持此网站 (%v) 添加新漫画", formData.Site),
		}
		r.Response.Status, response = lfunc.Response(cerror.ErrInvalidParameter, param)
	}

	_ = r.Response.WriteJson(response)
}

// Update 更新漫画信息
func (t *Book) Update(r *ghttp.Request) {
	response := g.Map{}

	id := r.GetInt("id")
	type form struct {
		Status int `params:"status" gvalid:"status@in:0,1,2"`
	}

	formData := new(form)
	if err := r.GetToStruct(formData); err != nil {
		r.Response.Status, response = lfunc.Response(cerror.ErrInvalidParameter)
		_ = r.Response.WriteJson(response)
		return
	}

	if err := gvalid.CheckStruct(formData, nil); err != nil {
		llog.Log.Warning(err.String())
		r.Response.Status, response = lfunc.Response(cerror.ErrInvalidParameter)
		_ = r.Response.WriteJson(response)
		return
	}

	data := g.Map{
		"status":     formData.Status,
		"updated_at": time.Now().Unix(),
	}

	m := model.NewBooks()
	_, err := m.UpdateData(data, g.Map{"id": id})
	if err != nil {
		llog.Log.Warningf(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrUpdateData, g.Map{"message": "更新漫画失败[Book.Update]"})
	} else {
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess)
	}

	_ = r.Response.WriteJson(response)
}

// Delete 删除漫画
func (t *Book) Delete(r *ghttp.Request) {
	response := g.Map{}

	id := r.GetInt("id")

	// 深度删除 (删除关联的章节及图库)
	isDeep := r.GetQueryBool("deep")

	m := model.NewBooks()
	_, err := m.DeleteData(g.Map{"id": id})
	if err != nil {
		llog.Log.Warningf(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrUpdateData, g.Map{"message": "更新漫画失败[Book.Update]"})
	} else {
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess)

		// 深度删除
		if isDeep {

		}
	}

	_ = r.Response.WriteJson(response)
}
