package api

import (
	"database/sql"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/library/ldb"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/library/llog"
	"github.com/skiy/comic-fetch/app/model"
	"net/http"
)

// Comic Comic
type Comic struct {
	core
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
	r.Response.Status = http.StatusBadRequest

	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	where := g.Map{}
	sort := "id asc"
	limit := 10

	// 漫画 ID
	bookID := r.GetInt("book_id")
	if bookID == 0 {
		response.Message = "漫画 ID 不存在"
		r.Response.WriteJson(response)
		return
	}

	chapterID := r.GetInt("chapter_id")
	if chapterID == 0 {
		response.Message = "漫画章节 ID 不存在"
		r.Response.WriteJson(response)
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

	Comics := ([]model.TbImages)(nil)
	resp, err := ldb.GetDB().Table(config.TbNameImages).Where(where).OrderBy(sort).Offset(offset).Limit(limit).Select()
	if err != nil && err != sql.ErrNoRows {
		llog.Log.Debug(err.Error())
		response.Message = "漫画章节图库列表获取失败[comic.List]"
	} else {
		if err != sql.ErrNoRows {
			if err := resp.ToStructs(&Comics); err != nil {
				response.Message = err.Error()
				r.Response.WriteJson(response)
				return
			}
		}
		r.Response.Status = http.StatusOK

		response.Code = 0
		response.Message = "操作成功"
		response.Data = Comics
	}

	r.Response.WriteJson(response)
}
