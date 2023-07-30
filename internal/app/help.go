package app

type app struct {
	storage map[string]string
}

func NewApp() *app {
	a := app{}
	a.storage = make(map[string]string)
	return &a
}
