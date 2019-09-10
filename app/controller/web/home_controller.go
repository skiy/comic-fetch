package web

import "github.com/gogf/gf/net/ghttp"

// HomeController HomeController
type HomeController struct {
}

// NewHomeController HomeController init
func NewHomeController() *HomeController {
	t := &HomeController{}
	return t
}

// Index Home Page
func (t *HomeController) Index(r *ghttp.Request) {
	r.Response.Write("string")
}
