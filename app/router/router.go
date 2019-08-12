package router

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/controller/web"
)

func init() {
	s := g.Server()

	homeController := web.NewHomeController()
	s.BindHandler("/", homeController.Index)

	s.SetIndexFolder(true)
	//s.SetServerRoot(".")
	//s.AddSearchPath("dist")
	s.AddStaticPath("/static", "dist")
}
