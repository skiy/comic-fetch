package app

import (
	"errors"
	"github.com/skiy/comic-fetch/app/controller"
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

	var ctrl controller.Controller

	if ctrl, err = t.ctrl("mh160"); err != nil {
		return err
	}

	ctrl.ToFetch()

	return nil
}

// ctrl 返回控制器
func (t *App) ctrl(name string) (ctrl controller.Controller, err error) {
	switch name {

	case "manhuaniu":
		ctrl = controller.NewManhuaniu()

	default:
		err = errors.New("can not fetch this comic website. ")

	}

	return
}
