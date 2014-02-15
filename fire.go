package fire

type Controller interface {
	Name() string
}

type Page struct {
	Title       string
	Description string
	Keywords    string
	Context     interface{}
}

func AddRootController(c Controller) {
	registerRootControllerName(c.Name())
	registerController(c)
}

func AddController(c Controller) {
	registerController(c)
}

func AddControllers(cs ...Controller) {
	for _, c := range cs {
		registerController(c)
	}
}
