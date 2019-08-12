package controller

import (
	"github.com/gogf/gf/g"
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
	err = g.Server().Start()
	if err != nil {
		return
	}

	g.Wait()
	return nil
}
