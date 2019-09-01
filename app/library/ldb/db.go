package ldb

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
)

// GetDB GetDB
func GetDB() (db gdb.DB) {
	db = g.DB()
	return db
}
