package api

import (
	"database/sql"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/library/ldb"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/library/llog"
	"github.com/skiy/comic-fetch/app/model"
	"net/http"
	"time"
)

// Chapter Chapter
type Chapter struct {
	core
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
	r.Response.Status = http.StatusBadRequest

	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	where := g.Map{}
	sort := "id desc"
	limit := 10

	// 漫画 ID
	bookID := r.GetInt("book_id")
	if bookID == 0 {
		response.Message = "漫画 ID 不存在"
		r.Response.WriteJson(response)
		return
	}
	where["book_id"] = bookID

	// 漫画章节 ID
	if i := r.GetInt("id"); i != 0 {
		where["id"] = i
	} else {
		// 漫画章节序号
		if num := r.GetInt("chapter_num"); num != 0 {
			where["order_id"] = num
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

	chapters := ([]model.TbChapters)(nil)
	resp, err := ldb.GetDB().Table(config.TbNameChapters).Where(where).OrderBy(sort).Offset(offset).Limit(limit).Select()
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Debug(err.Error())
		response.Message = "漫画章节列表获取失败[comic.List]"
	} else {
		if err != sql.ErrNoRows {
			if err := resp.ToStructs(&chapters); err != nil {
				response.Message = err.Error()
				r.Response.WriteJson(response)
				return
			}
		}

		r.Response.Status = http.StatusOK
		response.Code = 0
		response.Message = "操作成功"
		response.Data = chapters
	}

	r.Response.WriteJson(response)
}

// Update 更新漫画章节
func (t *Chapter) Update(r *ghttp.Request) {
	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"
	r.Response.Status = http.StatusBadRequest

	// 漫画 ID
	bookID := r.GetInt("book_id")
	if bookID == 0 {
		response.Message = "漫画 ID 不存在"
		r.Response.WriteJson(response)
		return
	}

	id := r.GetInt("id")
	if bookID == 0 {
		response.Message = "漫画章节 ID 不存在"
		r.Response.WriteJson(response)
		return
	}

	type form struct {
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
		"updated_at": time.Now().Unix(),
	}

	_, err := ldb.GetDB().Table(config.TbNameChapters).Where(g.Map{"id": id, "book_id": bookID}).Data(data).Update()
	if err != nil {
		llog.Log.Warningf(err.Error())
		response.Message = "漫画更新失败[Book.Update]"
	} else {
		r.Response.Status = http.StatusOK
		response.Code = 0
		response.Message = "操作成功"
	}

	r.Response.WriteJson(response)
}
