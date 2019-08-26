package api

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/skiy/comic-fetch/app/library/lfunc"
	"github.com/skiy/comic-fetch/app/model"
)

// Comic Comic
type Comic struct {
}

// NewComic Comic init
func NewComic() *Comic {
	t := &Comic{}
	return t
}

// List 漫画列表
func (t *Comic) List(r *ghttp.Request) {
	var response lfunc.Response
	response.Code = 1
	response.Message = "操作失败"

	books := ([]model.TbBooks)(nil)

	bookModel := model.NewBooks()
	resp, err := bookModel.GetData(g.Map{})
	if err != nil {
		response.Message = err.Error()
	} else {
		if err := resp.ToStructs(&books); err != nil {
			response.Message = err.Error()
		} else {
			response.Code = 0
			response.Message = "操作成功"
			response.Data = books
		}
	}

	r.Response.WriteJson(response)
}
