package fire

import (
	"net/http"
)

type Controller interface {
	Name() string
}

type Request struct {
	*http.Request
}

type Response struct {
	ControllerName string
	ActionName     string
	http.ResponseWriter
}

func (r Response) Render(data interface{}) {
	err := renderTemplate(r, r.ControllerName, r.ActionName, "layout", "", "", data)
	if err != nil {
		serveInternalError(r)
	}
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
