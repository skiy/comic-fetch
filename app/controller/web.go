package controller

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/gf-utils/ucfg"

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
	if t.Port <= 0 || t.Port > 65535 {
		t.Port = 33001
		if port := ucfg.InitCfg().GetInt("server.http.port"); port != 0 {
			t.Port = port
		}
	}

	s := g.Server()
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
