package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gcfg"
	"github.com/gogf/gf/g/os/glog"
	"github.com/skiy/comic-fetch/app"
	"github.com/skiy/gf-utils/ucfg"
	"github.com/skiy/gf-utils/udb"
	"github.com/skiy/gf-utils/ulog"
	"runtime"
)

var (
	cfg *gcfg.Config
	log *glog.Logger
)

func main() {
	//全核性能启用
	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化服务
	initialize()

	// 判断 MYSQL 连接是否正常
	if err := checkConnectDB(); err != nil {
		ulog.Log.Fatalf("数据库连接失败: %s", err.Error())
	}

	// 启动
	if err := app.NewApp().Start(); err != nil {
		ulog.Log.Fatalf("程序启动失败: %s", err.Error())
	}
}

// checkConnectDB 检测数据库连接是否正常
func checkConnectDB() (err error) {
	if err = udb.GetDatabase().PingMaster(); err != nil {
		return fmt.Errorf("%s(Database)", err.Error())
	}
	return err
}

// initialize 初始化服务
func initialize() {
	//配置文件
	cfg = ucfg.InitCfg()

	//日志初始化
	ulog.InitLog()

	//日志配置
	log = ulog.ReadLog()
}
