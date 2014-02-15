package fire

import (
	"io"
	"net/http"
)

type Response struct {
	ControllerName string
	ActionName     string
	http.ResponseWriter
}

func (r Response) Render(p Page) {
	layoutPath := getLayoutPath("layout")
	templatePath := getTemplatePath(r.ControllerName, r.ActionName)
	err := renderHTML(r, templatePath, layoutPath, p)
	if err != nil {
		serveInternalError(r)
	}
}

func (r Response) RenderFile(templatePath string, p Page) {
	layoutPath := getLayoutPath("layout")
	err := renderHTML(r, templatePath, layoutPath, p)
	if err != nil {
		serveInternalError(r)
	}
}

func (r Response) RenderNotFound() {
	r.WriteHeader(http.StatusNotFound)
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(r, "Not Found")
}
