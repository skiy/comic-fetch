package ldb

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// GetDB GetDB
func GetDB() (db gdb.DB) {
	db = g.DB()
	return db
}
