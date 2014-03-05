package fire

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	indexFilePath string
)

func Start() {
	StartWithPort(4000)
}

func StartWithPort(port int) {
	initTemplates()

	var err error
	if indexFilePath, err = filepath.Abs(filepath.Join(StaticFileBase(), "index.html")); err != nil {
		panic(err)
	}

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", frontController)

	fmt.Printf("Listening on port %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("An error occurred while starting the HTTP server.")
		panic(err)
	}
}

func frontController(rr http.ResponseWriter, req *http.Request) {
	rw := NewFireResponseWriter(rr)
	logMsg := fmt.Sprintf("%s - %s", req.Method, req.URL.Path)

	if pathIsRoot(req) {
		if indexFileExists() {
			http.ServeFile(rw, req, indexFilePath)
		} else {
			if err := runAction(getRootControllerName(), getDefaultActionName(), rw, req); err != nil {
				renderStatus(rw, http.StatusNotFound)
			}
		}
	} else {
		if filePath, err := staticFileExists(req); err == nil {
			http.ServeFile(rw, req, filePath)
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

	logMsg += fmt.Sprintf(" [%d|%s]", rw.Code, http.StatusText(rw.Code))
	log.Println(logMsg)
}

func indexFileExists() bool {
	_, err := os.Stat(indexFilePath)
	return err == nil
}

func staticFileExists(req *http.Request) (filename string, err error) {
	filename, err = filepath.Abs(filepath.Join(StaticFileBase(), req.URL.Path))
	if err != nil {
		return
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("Static file is a directory")
	}

	return
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
