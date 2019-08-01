package app

// App App
type App struct{}

// NewApp App init
func NewApp() *App {
	t := &App{}
	return t
}

// Start App start
func (t *App) Start() (err error) {
	return err
}
