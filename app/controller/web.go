package controller

import (
	"github.com/gogf/gf/g"
	"github.com/skiy/gf-utils/ucfg"

	// router
	_ "github.com/skiy/comic-fetch/app/router"
)

// Web Web
type Web struct {
	port int
}

// NewWeb Web init
func NewWeb() *Web {
	t := &Web{}
	return t
}

// SetPort Set web port with cli
func (t *Web) SetPort(port int) {
	t.port = port
}

// Start Web start
func (t *Web) Start() (err error) {
	if t.port <= 0 || t.port > 65535 {
		t.port = 33001
		if port := ucfg.InitCfg().GetInt("server.http.port"); port != 0 {
			t.port = port
		}
	}

	s := g.Server()
	s.SetPort(t.port)

	err = s.Start()
	if err != nil {
		return
	}

	g.Wait()
	return nil
}
