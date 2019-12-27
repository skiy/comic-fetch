package web

import (
	"github.com/gogf/gf/frame/g"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"

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
		if port := lcfg.Get().GetInt("server.http.port"); port != 0 {
			t.Port = port
		}
	}

	s := g.Server()

	// 静态网站路径
	if distPath := lcfg.Get().GetString("setting.template"); distPath != "" {
		llog.Log.Println(2, distPath)
		s.AddSearchPath(distPath)
	}

	// 设置服务日志路径
	if logPath := lcfg.Get().GetString("log.path"); logPath != "" {
		s.SetLogPath(logPath)
		s.SetAccessLogEnabled(true)
		s.SetErrorLogEnabled(true)
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
