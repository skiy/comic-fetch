package controller

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/gf-utils/ucfg"

	// router
	_ "github.com/skiy/comic-fetch/app/router"
)

// Web Web
type Web struct{}

// NewWeb Web init
func NewWeb() *Web {
	t := &Web{}
	return t
}

// Start Web start
func (t *Web) Start() (err error) {
	httpPort := 33001
	if p := ucfg.InitCfg().GetInt("server.http.port"); p != 0 {
		httpPort = p
	}

	s := g.Server()
	s.SetPort(httpPort)

	err = s.Start()
	if err != nil {
		return
	}

	g.Wait()
	return nil
}
