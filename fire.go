package fire

import (
	"net/http"
	"reflect"
)

// the interface that all controllers must implement
type Controller interface {
	Name() string
}

type RegisteredController struct {
	ControllerValue reflect.Value
	Actions         []RegisteredAction
}

type RegisteredAction struct {
	Name  string
	Index int
}

type Request struct {
	http.Request
}

type Response struct {
	ResponseWriter http.ResponseWriter
}
