package model

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

type model struct {
	DB gdb.DB
}

func (t *model) connect() {
	t.DB = g.DB()
}
