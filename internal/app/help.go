package app

type app struct {
	storage map[string]string
	cfg     netAdress
}

func NewApp() *app {
	a := app{}
	a.storage = make(map[string]string)
	return &a
}
