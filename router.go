package fire

import (
	"fmt"
	"reflect"
	"strings"
)

type actionInfo struct {
	controllerValue reflect.Value
	methodIndex     int
}

var (
	_responseType       reflect.Type
	_requestType        reflect.Type
	_actions            map[string]*actionInfo
	_rootControllerName string
)

func init() {
	_responseType = reflect.TypeOf(new(Response))
	_requestType = reflect.TypeOf(new(Request))
	_actions = make(map[string]*actionInfo)
}

func registerRootControllerName(name string) {
	_rootControllerName = name
}

func getRootControllerName() string {
	if _rootControllerName == "" {
		return "root"
	} else {
		return _rootControllerName
	}
}

func getDefaultActionName() string {
	return "index"
}

func registerController(c Controller) {
	ctrlValue := reflect.ValueOf(c)
	ctrlType := ctrlValue.Type()
	for i := 0; i < ctrlType.NumMethod(); i++ {
		method := ctrlType.Method(i)
		if method.Type.Kind() == reflect.Func && method.Type.NumIn() == 3 {
			if method.Type.In(1) == _responseType && method.Type.In(2) == _requestType {

				fmt.Println(ctrlValue.Elem().Type().Name(), "#", method.Name, "is an action")

				_actions[c.Name()+"/"+strings.ToLower(method.Name)] = &actionInfo{controllerValue: ctrlValue, methodIndex: i}
			}
		}
	}
}

func findAction(controllerName, actionName string) *actionInfo {
	return _actions[controllerName+"/"+strings.ToLower(actionName)]
}
