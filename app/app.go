package app

import (
	"errors"
	"github.com/skiy/comic-fetch/app/config"
	"github.com/skiy/comic-fetch/app/controller"
	"github.com/skiy/comic-fetch/app/model"
	"github.com/skiy/gf-utils/udb"
)

// App App
type App struct{}

// NewApp App init
func NewApp() *App {
	t := &App{}
	return t
}

// Start App start
func (t *App) Start() (err error) {
	books := ([]model.TbBooks)(nil)

	db := udb.GetDatabase()

	if err := db.Table(config.TbNameBooks).Structs(&books); err != nil {
		return err
	}

	var ctrl controller.Controller

	// 遍历表
	for _, book := range books {
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
func (t *App) ctrl(name string, books *model.TbBooks) (ctrl controller.Controller, err error) {
	switch name {

	case "manhuaniu":
		ctrl = controller.NewManhuaniu(books)

	default:
		err = errors.New("can not fetch this comic website. ")

	}

	return
}
