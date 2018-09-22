package library

import (
	"fmt"
	"github.com/go-ini/ini"
)

type Config struct {
	Setting *setting
	Mysql   *mysql
	Sqlite  *sqlite
	Redis   *redis
}

type setting struct {
	Datatype   string
	Imagefetch bool
}

type mysql struct {
	Host,
	User,
	Password,
	Name,
	Char string
}

type sqlite struct {
	Name string
}

type redis struct {
	Host,
	Port,
	Password string
	Db int
}

func (t *Config) ReadConfig() {
	cfg, err := ini.Load("conf.ini")
	if err != nil {
		fmt.Println(err)
	}

	s := new(setting)
	err = cfg.Section("Setting").MapTo(s)
	if err == nil {
		t.Setting = s
	}

	m := new(mysql)
	err = cfg.Section("Mysql").MapTo(m)
	if err == nil {
		t.Mysql = m
	}

	s2 := new(sqlite)
	err = cfg.Section("Mysql").MapTo(s2)
	if err == nil {
		t.Sqlite = s2
	}

	r := new(redis)
	err = cfg.Section("Redis").MapTo(r)
	if err == nil {
		t.Redis = r
	}
}
