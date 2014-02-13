package fire

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var (
	rootControllerName string
	controllers        map[string]interface{}
)

func init() {
	controllers = make(map[string]interface{})
}

func AddRootController(c Controller) {
	rootControllerName = c.Name()
	controllers[c.Name()] = c
}

func AddController(c Controller) {
	controllers[c.Name()] = c
}

func AddControllers(cs ...Controller) {
	for _, c := range cs {
		controllers[c.Name()] = c
	}
}

func Start(port int) {
	http.HandleFunc("/", frontController)

	// static files; in production these *should* be served by a web server
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js/"))))

	fmt.Printf("Listening on port %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("An error occurred while starting the HTTP server.")
		panic(err)
	}
}

func frontController(rw http.ResponseWriter, req *http.Request) {
	if pathIsRoot(req) {
		if err := runAction(rootControllerName, "index", rw, req); err != nil {
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
	ctrl := controllers[controllerName]
	if ctrl == nil {
		return errors.New("Controller not found")
	}

	ctrlValue := reflect.ValueOf(ctrl)
	ctrlType := ctrlValue.Type()
	mIndex := -1

	for i := 0; i < ctrlType.NumMethod(); i++ {
		if actionName == strings.ToLower(ctrlType.Method(i).Name) {
			mIndex = i
			break
		}
	}

	if mIndex == -1 {
		return errors.New("Action not found")
	} else {
		// TODO: check to see if method implements http.Handle
	}

	in := []reflect.Value{reflect.ValueOf(rw), reflect.ValueOf(req)}
	ctrlValue.Method(mIndex).Call(in)
	return nil
}
