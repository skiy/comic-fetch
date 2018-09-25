package library

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

type Database struct {
	Datatype string
	DSN,
	Dbhost,
	Dbuser,
	Dbpwd,
	Dbname,
	Dbchar string
}

/**
初始化数据库
*/
func (t *Database) Init(host, user, pwd, name, char string) {
	t.Dbhost = host
	t.Dbuser = user
	t.Dbpwd = pwd
	t.Dbname = name
	t.Dbchar = char

	t.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		t.Dbuser,
		t.Dbpwd,
		t.Dbhost,
		t.Dbname,
		t.Dbchar,
	)
}

/**
连接数据库
*/
func (t *Database) Connect() (db *gorm.DB, err error) {

	if t.Datatype == "mysql" {
		if t.DSN != "" {
			return gorm.Open("mysql", t.DSN)
		}
	} else if t.Datatype == "sqlite" {
		return gorm.Open("sqlite3", t.Dbname+".db")
	}

	return nil, errors.New("数据库连接失败")
}

/**
  数据库字段批量加前缀
*/
func (t *Database) FieldAddPrev(prev, fieldStr string) string {
	fieldArr := strings.Split(fieldStr, ",")

	prev = prev + "."
	var newFieldArr []string
	for _, v := range fieldArr {
		newFieldArr = append(newFieldArr, prev+v)
	}
	newFieldStr := strings.Join(newFieldArr, ",")

	return newFieldStr
}

/**
将
*/
func (t *Database) FieldMakeQmark(str string, symbol string) string {
	strArr := strings.Split(str, ",")
	strLen := len(strArr)

	if symbol == "" {
		symbol = "?"
	}

	arr := make([]string, strLen)
	for i := 0; i < strLen; i++ {
		arr[i] = symbol
	}
	return strings.Join(arr, ",")
}

/**
将
*/
func (t *Database) FieldMakeValue(str string) string {
	strArr := strings.Split(str, ",")
	strLen := len(strArr)

	arr := make([]string, strLen)
	for i := 0; i < strLen; i++ {
		arr[i] = "'%s'"
	}
	return strings.Join(arr, ",")
}

/**
将 a,b,c 处理成 'a','b','c'
*/
func (t *Database) ValueMakeData(str string) string {
	strArr := strings.Split(str, ",")
	strLen := len(strArr)

	if strLen == 0 {
		return ""
	}

	var arr []string
	for _, v := range strArr {
		arr = append(arr, "'"+v+"'")
	}

	return strings.Join(arr, ",")
}
