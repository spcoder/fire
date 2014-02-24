package fire

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func Start() {
	StartWithPort(4000)
}

func StartWithPort(port int) {
	initTemplates()

	http.HandleFunc("/", frontController)

	staticFileBase := fmt.Sprintf("/%s/", StaticFileBase())

	// static files; in production these *should* be served by a web server
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle(staticFileBase, http.StripPrefix(staticFileBase, http.FileServer(http.Dir("."+staticFileBase))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("."+staticFileBase+"css/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("."+staticFileBase+"img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("."+staticFileBase+"js/"))))

	fmt.Printf("Listening on port %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("An error occurred while starting the HTTP server.")
		panic(err)
	}
}

func frontController(rw http.ResponseWriter, req *http.Request) {
	if pathIsRoot(req) {
		if err := runAction(getRootControllerName(), getDefaultActionName(), rw, req); err != nil {
			renderStatus(rw, http.StatusNotFound)
		}
	} else {
		if controllerName, actionName, err := parsePath(req); err == nil {
			if err := runAction(controllerName, actionName, rw, req); err != nil {
				renderStatus(rw, http.StatusNotFound)
			}
		} else {
			renderStatus(rw, http.StatusBadRequest)
		}
	}
}

func pathIsRoot(req *http.Request) bool {
	return req.URL.Path == "/"
}

func parsePath(req *http.Request) (controllerName, actionName string, err error) {
	splitPath := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")
	if len(splitPath) < 2 {
		return "", "", errors.New("Bad Request")
	} else {
		return splitPath[0], splitPath[1], nil
	}
}

func runAction(controllerName, actionName string, rw http.ResponseWriter, req *http.Request) (err error) {
	ai := findAction(controllerName, actionName)
	if ai == nil {
		return errors.New("Action not found")
	}

	resp := &Response{ControllerName: controllerName, ActionName: actionName, ResponseWriter: rw}
	r := &Request{HttpRequest: req}
	in := []reflect.Value{reflect.ValueOf(resp), reflect.ValueOf(r)}
	ai.controllerValue.Method(ai.methodIndex).Call(in)
	return nil
}
