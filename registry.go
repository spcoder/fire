package fire

import (
	"reflect"
)

type action struct {
	name  string
	index int
}

type controller struct {
	name    string
	value   reflect.Value
	actions []action
}

type registry struct {
	controllers []controller
}

var (
	_registry *registry
)

func init() {
	_registry = &registry{controllers: make([]controller, 0)}
}

func (r *registry) registerController(c Controller) {

}
