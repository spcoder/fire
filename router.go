package fire

import (
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

func initRouter() {
	_responseType = reflect.TypeOf(new(Response))
	_requestType = reflect.TypeOf(new(Request))
	_actions = make(map[string]*actionInfo)
}

func getRootControllerName() string {
	return _rootControllerName
}

func getDefaultActionName() string {
	return "index"
}

func registerController(c Controller) {
	if _rootControllerName == "" {
		_rootControllerName = c.Name()
	}
	ctrlValue := reflect.ValueOf(c)
	ctrlType := ctrlValue.Type()
	for i := 0; i < ctrlType.NumMethod(); i++ {
		method := ctrlType.Method(i)
		if method.Type.Kind() == reflect.Func && method.Type.NumIn() == 3 {
			if method.Type.In(1) == _responseType && method.Type.In(2) == _requestType {
				_actions[c.Name()+"/"+strings.ToLower(method.Name)] = &actionInfo{controllerValue: ctrlValue, methodIndex: i}
			}
		}
	}
}

func findAction(controllerName, actionName string) *actionInfo {
	return _actions[controllerName+"/"+strings.ToLower(actionName)]
}
