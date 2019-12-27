package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/glog"
	"github.com/skiy/comic-fetch/app/config"
	command2 "github.com/skiy/comic-fetch/app/service/command"
	"github.com/skiy/comic-fetch/app/service/web"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/ldb"
	"github.com/skiy/gfutils/llog"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"time"
)

type command struct {
	help   bool   // h 帮助
	lang   string // l 语言
	cmd    string // web / cli 启动方式
	port   int    // --port WEB 端口
	config string // --config 配置文件路径
}

var (
	cfg *gcfg.Config
	log *glog.Logger
	err error
	cmd command
)

const (
	version = "2.0.0"
	author  = "Skiy Chan"
	email   = "dev@skiy.net"
)

// load 加载配置信息
func load() {
	if cmd.config != "" {
		lcfg.SetCfgName(cmd.config)
	} else {
		lcfg.SetCfgName("config.toml")
	}

	cfg, err = lcfg.Init()
	if err != nil {
		return
	}

	err = llog.Init()
	if err != nil {
		return
	}
	log = llog.Log

	err = ldb.Init()
	if err != nil {
		return
	}
}

func main() {
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
			Aliases:     []string{"s"},
			Value:       "",
			DefaultText: "",
		},
		&cli.IntFlag{
			Name:        "id",
			Usage:       "Origin comic id Or comic id",
			Aliases:     []string{"i"},
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
			{
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

			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},

		// default cli
		Action: func(c *cli.Context) error {
			cmd.config = c.String("config")
			cliStart()
			return nil
		},

		// command cli / web
		Commands: []*cli.Command{
			// web
			{
				Name:  "web",
				Usage: "Comic website run",
				Action: func(c *cli.Context) error {
					cmd.config = c.String("config")
					cmd.port = c.Int("port")

					webStart()
					return nil
				},
				// web port
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Usage:       "set website port",
						Value:       0,
						Aliases:     []string{"p"},
						DefaultText: "0",
					},
				},
			},

			// cli
			{
				Name:  "cli",
				Usage: "Comic fetch run",
				Action: func(c *cli.Context) error {
					cmd.config = c.String("config")

					cliStart()
					return nil
				},
				Subcommands: []*cli.Command{
					// add comic
					{
						Name:  "add",
						Usage: "Add a new comic",
						Action: func(c *cli.Context) error {
							cmd.config = c.String("config")
							load()

							if err != nil {
								log.Fatalf("%s\n", err.Error())
							}

							site := c.String("site")
							id := c.Int("id")

							if id == 0 {
								log.Warningf("漫画 (%s) 参数 id 缺失", site)
								return nil
							}

							if _, ok := config.WebURL[site]; ok {
								cliApp := command2.NewCommand()

								if err := cliApp.Add(site, id); err != nil {
									log.Fatalf("添加新漫画失败: %s", err.Error())
									return nil
								}
							} else {
								log.Warningf("不支持此网站 (%v) 添加新漫画", site)
							}

							return nil
						},
						Flags: cliFlags,
					},
					// update comic
					{
						Name:  "update",
						Usage: "Update a comic",
						Action: func(c *cli.Context) error {
							cmd.config = c.String("config")
							load()

							if err != nil {
								log.Fatalf("%s\n", err.Error())
							}

							site := c.String("site")
							id := c.Int("id")

							cliApp := command2.NewCommand()
							where := g.Map{}

							if site != "" {
								where["origin_flag"] = site

								if id != 0 {
									where["origin_book_id"] = id
								}
							} else if id != 0 {
								where["id"] = id
							}

							if err := cliApp.Update(where); err != nil {
								log.Fatalf("更新漫画失败: %s", err.Error())
								return nil
							}

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

// webStart web run
func webStart() {
	load()

	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	app := web.NewWeb()
	app.Port = cmd.port

	// 启动
	if err := app.Start(); err != nil {
		log.Fatalf("WEB 程序启动失败: %s", err.Error())
	}
}

// cliStart cli run
func cliStart() {
	load()

	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	app := command2.NewCommand()

	// 启动
	if err := app.Update(g.Map{}); err != nil {
		log.Fatalf("CLI 程序启动失败: %s", err.Error())
	}
}
