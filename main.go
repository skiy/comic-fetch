package main

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gcfg"
	"github.com/gogf/gf/g/os/glog"
	"github.com/skiy/comic-fetch/app/controller"
	"github.com/skiy/gf-utils/ucfg"
	"github.com/skiy/gf-utils/udb"
	"github.com/skiy/gf-utils/ulog"
	"os"
	"runtime"
)

var (
	cfg *gcfg.Config
	log *glog.Logger
)

var (
	// h 帮助
	h bool
	// v 版本号
	v bool
	// 添加漫画
	a bool
	// web or cli start
	s string
	// S 漫画网站标识 ( manhuaniu etc... )
	S string
	// I 漫画源 ID
	I int64
)

const version = "1.0.0"

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&a, "add", false, "add new comic")
	flag.StringVar(&s, "s", "cli", "comic start (web or cli)")
	flag.StringVar(&S, "S", "", "support sites: manhuaniu")
	flag.Int64Var(&I, "I", 0, "comic origin id")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `comic version: comic/%s
Usage: comic [-hvVtTq] [-web web run] [-S comic website] [-I comic id] [-add]

Options:
`, version)
	flag.PrintDefaults()
}

func main() {
	//全核性能启用
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if h { // 显示帮助
		flag.Usage()
	} else if v { // 显示版本号
		fmt.Printf("comic version: comic/%s\n", version)
	} else if s == "web" {
		fmt.Printf("comic web run\n")
		start("web")
	} else { // 默认采集
		if S != "" { // 更新 或 添加 指定的漫画
			supportSites := g.Map{
				"manhuaniu": true,
				"abc":       false,
			}

			if s, ok := supportSites[S]; !ok || !s.(bool) {
				fmt.Printf("unsupport this website: %s\n", S)
			} else {
				// 添加新漫画
				if a {
					if I == 0 { // 新漫画必须指定 ID
						fmt.Printf("must set comic origin id: (I)\n")
					} else { // 添加站点 S 的新漫 I

					}
				} else { // 更新漫画
					if I == 0 { // 更新此 S 站点的所有漫画
					} else { // 更新站点 S 的新漫 I

					}
				}
			}

		} else {
			start("command")
		}
	}
}

// start app / web start
func start(flag string) {
	// 判断 MYSQL 连接是否正常
	if err := checkConnectDB(); err != nil {
		ulog.Log.Fatalf("数据库连接失败: %s", err.Error())
	}

	var app controller.Controller
	if flag == "command" {
		app = controller.NewCommand()
	} else {
		app = controller.NewWeb()
	}

	// 启动
	if err := app.Start(); err != nil {
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

// init 初始化服务
func init() {
	//配置文件
	cfg = ucfg.InitCfg()

	//日志初始化
	ulog.InitLog()

	//日志配置
	log = ulog.ReadLog()
}
