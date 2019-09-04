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
		if err := r.Response.WriteJson(response); err != nil {
			r.Response.Status = 500
		}
		return
	}
	where["book_id"] = bookID

	// 漫画章节 ID
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
			}
		}

		response.Code = 0
		response.Message = "操作成功"
		response.Data = chapters
	}

	if err := r.Response.WriteJson(response); err != nil {
		r.Response.Status = 500
	}
}
