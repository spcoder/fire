package fire

func init() {
	initSettings()
	initRouter()
}

type Controller interface {
	Name() string
}

type Page struct {
	Title       string
	Description string
	Keywords    string
	Data        interface{}
}

func AddController(c Controller) {
	registerController(c)
}

func AddControllers(cs ...Controller) {
	for _, c := range cs {
		registerController(c)
	}
}
