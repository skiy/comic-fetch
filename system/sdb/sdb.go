package sdb

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
)

// GetDatabase GetDatabase
func GetDatabase() (db gdb.DB) {
	db = g.Database()
	return db
}
