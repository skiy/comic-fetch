package controller

import (
	"errors"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/controller/command"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gf-utils/udb"
)

// Command Command
type Command struct{}

// NewCommand Command init
func NewCommand() *Command {
	t := &Command{}
	return t
}

// Start Command start
func (t *Command) Start() (err error) {
	books := ([]model.TbBooks)(nil)

	db := udb.GetDatabase()

	if err := db.Table(config.TbNameBooks).Structs(&books); err != nil {
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

		err = ctrl.ToFetch()
		if err != nil {
			return err
		}
	}

	return nil
}

// ctrl 返回控制器
func (t *Command) ctrl(name string, books *model.TbBooks) (ctrl command.Controller, err error) {
	switch name {

	case "manhuaniu":
		ctrl = command.NewManhuaniu(books)

	default:
		err = errors.New("can not fetch this comic website. ")

	}

	return
}
