package router

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/controller/web"
	"github.com/skiy/gf-utils/ucfg"
)

func init() {
	httpPort := 33001
	if p := ucfg.InitCfg().GetInt("server.http.port"); p != 0 {
		httpPort = p
	}

	s := g.Server()
	s.SetPort(httpPort)

	homeController := web.NewHomeController()
	s.BindHandler("/", homeController.Index)
}
