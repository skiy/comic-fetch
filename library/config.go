package library

import (
	"github.com/go-ini/ini"
	"log"
)

type Config struct {
	Setting *setting
	Mysql   *mysql
	Sqlite  *sqlite
	Redis   *redis
	Image   *image
}

type setting struct {
	Datatype string
	ImageFetch,
	Cache bool
	WebPort int
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

type image struct {
	Path,
	Nametype string
}

func (t *Config) ReadConfig() {
	realpath, err := GetCurrentDirectory()
	if err != nil {
		log.Fatalln(err)
	}

	confpath := realpath + "/conf.ini"
	cfg, err := ini.Load(confpath)

	if err != nil {
		cfg, err = ini.Load("conf.ini")
	}

	if err != nil {
		log.Fatalln("Load conf.ini fail", err)
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
	err = cfg.Section("Sqlite").MapTo(s2)
	if err == nil {
		t.Sqlite = s2
	}

	r := new(redis)
	err = cfg.Section("Redis").MapTo(r)
	if err == nil {
		t.Redis = r
	}

	i := new(image)
	err = cfg.Section("Image").MapTo(i)
	if err == nil {
		t.Image = i
	}
}
