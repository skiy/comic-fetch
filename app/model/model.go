package model

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
)

type model struct {
	DB gdb.DB
}

func (t *model) connect() {
	t.DB = g.DB()
}
