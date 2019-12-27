package api

import (
	"database/sql"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/skiy/comic-fetch/app/config/cerror"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gfutils/llog"
	"time"
)

// Chapter Chapter
type Chapter struct {
	Base
}

// NewChapter Chapter init
func NewChapter() *Chapter {
	t := &Chapter{}
	return t
}

// List 漫画列表
// sort=+id,-name 排序
// offset=10&limit=5 分页
// /api/books/:book_id/chapters[/:id]
// /api/books/:book_id/parts[/:chapter_num]
func (t *Chapter) List(r *ghttp.Request) {
	response := g.Map{}

	where := g.Map{}
	sort := "id desc"
	limit := 10

	// 漫画 ID
	bookID := r.GetInt("book_id")
	if bookID == 0 {
		r.Response.Status, response = lfunc.Response(cerror.ErrFailure, g.Map{"message": "漫画 ID 不存在"})
		_ = r.Response.WriteJson(response)
		return
	}
	where["book_id"] = bookID

	// 漫画章节 ID
	if i := r.GetInt("id"); i != 0 {
		where["id"] = i
	} else {
		llog.Log.Println(r.GetInt("chapter_num"))
		// 漫画章节序号
		if num := r.GetInt("chapter_num"); num != 0 {
			where["order_id"] = num
			limit = 1
		}
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

	chapters := ([]model.TbChapters)(nil)
	m := model.NewChapters()
	resp, err := m.GetDataExt(p)
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Warning(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrGetData, g.Map{"message": "获取漫画章节列表失败[Chapter.List]"})
	} else {
		if err != sql.ErrNoRows && resp != nil {
			if err := resp.ToStructs(&chapters); err != nil {
				llog.Log.Warning(err.Error())
			}
		}

		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess, g.Map{"data": chapters})
	}

	_ = r.Response.WriteJson(response)
}

// Update 更新漫画章节
func (t *Chapter) Update(r *ghttp.Request) {
	response := g.Map{}

	// 漫画 ID
	bookID := r.GetInt("book_id")
	if bookID == 0 {
		r.Response.Status, response = lfunc.Response(cerror.ErrBookIDNotExist)
		_ = r.Response.WriteJson(response)
		return
	}

	id := r.GetInt("id")
	if bookID == 0 {
		r.Response.Status, response = lfunc.Response(cerror.ErrFailure, g.Map{"message": ""})
		_ = r.Response.WriteJson(response)
		return
	}

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

	m := model.NewChapters()
	_, err := m.UpdateData(data, g.Map{"id": id, "book_id": bookID})
	if err != nil {
		llog.Log.Warning(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrUpdateData, g.Map{"message": "更新漫画章节失败[Chapter.Update]"})
	} else {
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess)
	}

	_ = r.Response.WriteJson(response)
}

// Delete 删除漫画章节
func (t *Chapter) Delete(r *ghttp.Request) {
	response := g.Map{}

	id := r.GetInt("id")

	// 深度删除 (删除关联的图库)
	isDeep := r.GetQueryBool("deep")

	m := model.NewChapters()
	_, err := m.DeleteData(g.Map{"id": id})
	if err != nil {
		llog.Log.Warningf(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrDeleteData, g.Map{"message": "删除漫画章节失败[Chapter.Delete]"})
	} else {
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess)

		// 深度删除
		if isDeep {

		}
	}

	_ = r.Response.WriteJson(response)
}
