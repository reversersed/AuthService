package app

func New() (*app, error) {
	app := &app{}
	return app, nil
}
func (a *app) Run() error {
	return nil
}
