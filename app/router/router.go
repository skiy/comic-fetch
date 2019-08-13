package router

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/controller/web"
)

func init() {
	s := g.Server()

	// Static setting
	//s.SetIndexFolder(true)
	//s.SetServerRoot(".")
	s.AddSearchPath("dist")

	//s.AddStaticPath("/", "home/index.html")
	//s.AddStaticPath("/static", "home/static")

	homeController := web.NewHomeController()
	s.BindHandler("/abc", homeController.Index)
}
