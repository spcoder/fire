package fire

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func Start() {
	StartWithPort(4000)
}

func StartWithPort(port int) {
	http.HandleFunc("/", frontController)

	// static files; in production these *should* be served by a web server
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("./www/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/css/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./www/img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./www/js/"))))

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
			serveNotFound(rw)
		}
	} else {
		if controllerName, actionName, err := parsePath(req); err == nil {
			if err := runAction(controllerName, actionName, rw, req); err != nil {
				serveNotFound(rw)
			}
		} else {
			serveBadRequest(rw)
		}
	}
}

func serveBadRequest(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, "Bad Request")
}

func serveNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, "Not Found")
}

func serveInternalError(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, "Internal Server Error")
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
	rr := &Request{req}
	in := []reflect.Value{reflect.ValueOf(resp), reflect.ValueOf(rr)}
	ai.controllerValue.Method(ai.methodIndex).Call(in)
	return nil
}
