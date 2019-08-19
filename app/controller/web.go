package controller

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/comic-fetch/app/library/lcfg"
	"github.com/skiy/comic-fetch/app/library/llog"

	// router
	_ "github.com/skiy/comic-fetch/app/router"
)

// Web Web
type Web struct {
	Port int
}

// NewWeb Web init
func NewWeb() *Web {
	t := &Web{}
	return t
}

// Start Web start
func (t *Web) Start() (err error) {
	// WEB 端口
	if t.Port <= 0 || t.Port > 65535 {
		t.Port = 33001
		if port := lcfg.InitCfg().GetInt("server.http.port"); port != 0 {
			t.Port = port
		}
	}

	s := g.Server()

	// 静态网站路径
	if distPath := lcfg.GetCfg().GetString("setting.template"); distPath != "" {
		llog.Log.Println(2, distPath)
		s.AddSearchPath(distPath)
	}

	s.SetPort(t.Port)

	// 关闭平滑重启功能
	g.SetServerGraceful(false)

	err = s.Start()
	if err != nil {
		return
	}

	g.Wait()
	return nil
}
