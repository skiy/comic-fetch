package router

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/controller/web"
	"github.com/skiy/comic-fetch/app/controller/web/api"
)

func init() {
	s := g.Server()

	// Static setting
	//s.SetIndexFolder(true)
	//s.SetServerRoot(".")

	//s.AddStaticPath("/", "home/index.html")
	//s.AddStaticPath("/static", "home/static")

	webHome := web.NewHomeController()
	s.BindHandler("/", webHome.Index)

	apiGroup := s.Group("/api")

	apiComic := api.NewComic()
	apiGroup.GET("/comics", apiComic.List)
}
