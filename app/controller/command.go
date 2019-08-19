package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/controller/command"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gf-utils/ulog"
)

// Command Command
type Command struct {
}

// NewCommand Command init
func NewCommand() *Command {
	t := &Command{}
	return t
}

// Update Update comics
func (t *Command) Update(where interface{}) (err error) {
	books := ([]model.TbBooks)(nil)

	bookModel := model.NewBooks()
	resp, err := bookModel.GetData(where)
	if err != nil {
		return err
	}

	if err := resp.ToStructs(&books); err != nil {
		return err
	}

	var ctrl command.Controller

	// 遍历表
	for _, book := range books {
		// 更新的状态, 非正在更新
		if book.Status != 0 {
			continue
		}

		ctrl, err = t.ctrl(book.OriginFlag, &book)
		if err != nil {
			return err
		}

		ulog.Log.Infof("\n正在采集漫画: %s\n源站: %s\n源站漫画URL: %s\n", book.Name, book.OriginWeb, book.OriginURL)

		err = ctrl.ToFetch()
		if err != nil {
			return err
		}
	}

	return nil
}

// Add 添加新漫画
func (t *Command) Add(flag string, bookID int) (err error) {
	where := g.Map{
		"origin_flag":    flag,
		"origin_book_id": bookID,
	}

	bookModel := model.NewBooks()
	bookRes, err := bookModel.GetDataOne(where)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if bookRes != nil {
		return fmt.Errorf("漫画已存在: %v", where)
	}

	var siteURL, originWeb, originWebType string

	if site, ok := config.WebURL[flag]; ok {
		originWebType = gconv.String(site["flag"])
		originWeb = gconv.String(site["name"])

		if sURL, ok := site[originWebType]; ok {
			siteURL = sURL
		} else {
			return fmt.Errorf("此网站 (%v) 添加新漫画方式有误: %v", flag, site["flag"])
		}

	} else {
		return fmt.Errorf("不支持此网站 (%v) 添加新漫画", flag)
	}

	book := &model.TbBooks{}
	book.OriginWeb = originWeb
	book.OriginWebType = originWebType
	book.OriginFlag = flag
	book.OriginBookID = bookID

	var ctrl command.Controller

	ctrl, err = t.ctrl(book.OriginFlag, book)
	if err != nil {
		return err
	}

	ulog.Log.Infof("\n正在新增漫画 \n源站: %s\n源站漫画ID: %d\n", book.OriginWeb, book.OriginBookID)

	err = ctrl.AddBook(siteURL)
	if err != nil {
		return err
	}

	return
}

// ctrl 返回控制器
func (t *Command) ctrl(name string, books *model.TbBooks) (ctrl command.Controller, err error) {
	switch name {

	case "manhuaniu":
		ctrl = command.NewManhuaniu(books)

	case "mh1234":
		ctrl = command.NewMh1234(books)

	default:
		err = errors.New("can not fetch this comic website. ")

	}

	return
}
