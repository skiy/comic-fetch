package web

import "github.com/gogf/gf/g/net/ghttp"

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
