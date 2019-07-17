package scfg

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gcfg"
)

var cfg *gcfg.Config

func InitCfg() *gcfg.Config {
	cfg = g.Config()
	return cfg
}

func GetCfg() *gcfg.Config {
	return cfg
}
