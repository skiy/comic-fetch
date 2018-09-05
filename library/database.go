package library

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)

type Database struct {
	DSN,
	Dbhost,
	Dbuser,
	Dbpwd,
	Dbname,
	Dbchar string
}

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

func (t *Database) Connect() (db *gorm.DB, err error) {

	if t.DSN != "" {
		return gorm.Open("mysql", t.DSN)
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
