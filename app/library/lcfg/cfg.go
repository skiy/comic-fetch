package lcfg

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
)

var cfg *gcfg.Config

// SetCfgName config path set
func SetCfgName(path string) {
	g.Config().SetFileName(path)
}

// InitCfg config init
func InitCfg() *gcfg.Config {
	cfg = g.Config()
	return cfg
}

// GetCfg get config
func GetCfg() *gcfg.Config {
	return cfg
}
