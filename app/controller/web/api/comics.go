package api

import (
	"database/sql"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/skiy/comic-fetch/app/config/cerror"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/library/llog"
	"github.com/skiy/comic-fetch/app/model"
)

// Comic Comic
type Comic struct {
	Base
}

// NewComic Comic init
func NewComic() *Comic {
	t := &Comic{}
	return t
}

// List 漫画列表
// sort=+id,-name 排序
// offset=10&limit=5 分页
// /api/books/:book_id/chapters/:chapter_id/comics[/:id]
// /api/books/:book_id/chapters/:chapter_id/parts[/:comic_num]
func (t *Comic) List(r *ghttp.Request) {
	response := g.Map{}

	where := g.Map{}
	sort := "id asc"
	limit := 0

	// 漫画 ID
	bookID := r.GetInt("book_id")
	if bookID == 0 {
		r.Response.Status, response = lfunc.Response(cerror.ErrBookIDNotExist)
		_ = r.Response.WriteJson(response)
		return
	}

	chapterID := r.GetInt("chapter_id")
	if chapterID == 0 {
		r.Response.Status, response = lfunc.Response(cerror.ErrChapterIDNotExist)
		_ = r.Response.WriteJson(response)
		return
	}

	where["book_id"] = bookID
	where["chapter_id"] = chapterID

	// 漫画章节图 ID
	if i := r.GetInt("id"); i != 0 {
		where["id"] = i
	} else {
		// 漫画章节序号
		if num := r.GetInt("comic_num"); num != 0 {
			where["order_id"] = num
			limit = 1
		}
	}

	// 排序
	if s := r.GetQueryString("sort"); s != "" {
		filterFields := map[string]bool{
			"id": true,
		}
		if s2 := t.sort(s, filterFields); s2 != "" {
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

	comics := ([]model.TbImages)(nil)
	m := model.NewImages()
	resp, err := m.GetDataExt(p)
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Warning(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrGetData, g.Map{"message": "漫画章节图库列表获取失败[Comic.List]"})
	} else {
		if err != sql.ErrNoRows && resp != nil {
			if err := resp.ToStructs(&comics); err != nil {
				llog.Log.Warning(err.Error())
				return
			}
		}
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess, g.Map{"data": comics})
	}

	_ = r.Response.WriteJson(response)
}

// Delete 删除漫画图库
func (t *Comic) Delete(r *ghttp.Request) {
	response := g.Map{}

	id := r.GetInt("id")

	m := model.NewImages()
	_, err := m.DeleteData(g.Map{"id": id})
	if err != nil {
		llog.Log.Warningf(err.Error())
		r.Response.Status, response = lfunc.Response(cerror.ErrDeleteData, g.Map{"message": "删除漫画图库失败[Image.Delete]"})
	} else {
		r.Response.Status, response = lfunc.Response(cerror.ErrSuccess)
	}

	_ = r.Response.WriteJson(response)
}
