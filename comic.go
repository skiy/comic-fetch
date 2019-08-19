package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gcfg"
	"github.com/gogf/gf/g/os/glog"
	"github.com/skiy/comic-fetch/app/controller"
	"github.com/skiy/gf-utils/ucfg"
	"github.com/skiy/gf-utils/udb"
	"github.com/skiy/gf-utils/ulog"
	"gopkg.in/urfave/cli.v2"
	"os"
	"runtime"
	"sort"
	"time"
)

type command struct {
	help bool   // h 帮助
	lang string // l 语言
	cmd  string // web / cli 启动方式
	port int    // --port WEB 端口
}

var (
	cfg *gcfg.Config
	log *glog.Logger
	cmd command
)

const (
	version = "1.0.0"
	author  = "Skiy Chan"
	email   = "dev@skiy.net"
)

func main() {
	//全核性能启用
	runtime.GOMAXPROCS(runtime.NumCPU())

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Show version and exit",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "This help",
	}

	cliFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "site",
			Usage:       "Fetch comic from this site (support: manhuaniu, mh1234)",
			Value:       "",
			DefaultText: "",
		},
		&cli.IntFlag{
			Name:        "id",
			Usage:       "Origin comic id",
			Value:       0,
			DefaultText: "0",
		},
	}

	app := &cli.App{
		Name:     "comic",
		Usage:    "Application for fetch comic",
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  author,
				Email: email,
			},
		},
		Copyright: "(c) 2019 Skiy Chan",

		Flags: []cli.Flag{
			//&cli.StringFlag{
			//	Name:  "lang",
			//	Value: "english",
			//	Aliases: []string{"l"},
			//	Usage: "language for the greeting",
			//	//EnvVars: []string{"LANG"},
			//},

			//&cli.IntFlag{
			//	Name:    "port",
			//	Usage:   "set website port",
			//	Value: 0,
			//	DefaultText: "random",
			//},
		},

		// default cli
		Action: func(c *cli.Context) error {
			start("command")
			return nil
		},

		// command cli / web
		Commands: []*cli.Command{
			// web
			{
				Name:  "web",
				Usage: "Comic website run",
				Action: func(c *cli.Context) error {
					cmd.port = c.Int("port")
					fmt.Println("Comic website run: ", cmd.port)
					start("web")
					return nil
				},
				// web port
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Usage:       "set website port",
						Value:       0,
						DefaultText: "0",
					},
				},
			},

			// cli
			{
				Name:  "cli",
				Usage: "Comic fetch run",
				Action: func(c *cli.Context) error {
					fmt.Println("Comic fetch run")
					start("command")
					return nil
				},
				Subcommands: []*cli.Command{
					// add comic
					{
						Name:  "add",
						Usage: "Add a new comic",
						Action: func(c *cli.Context) error {
							site := c.String("site")
							id := c.String("id")

							fmt.Println("Add a new comic: ", site, id)
							return nil
						},
						Flags: cliFlags,
					},
					// update comic
					{
						Name:  "update",
						Usage: "Update a comic",
						Action: func(c *cli.Context) error {
							site := c.String("site")
							id := c.String("id")

							fmt.Println("Update a comic: ", site, id)
							return nil
						},
						Flags: cliFlags,
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
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
		app.SetPort(cmd.port)
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
