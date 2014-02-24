package fire

import (
	"io"
	"log"
	"net/http"
)

type Response struct {
	ControllerName string
	ActionName     string
	ResponseWriter http.ResponseWriter
}

func (r Response) Render(data interface{}) {
	r.RenderWithLayout("layout", data)
}

func (r Response) RenderWithLayout(layoutPath string, data interface{}) {
	resolvedLayoutPath := getLayoutPath(layoutPath)
	templatePath := getTemplatePath(r.ControllerName, r.ActionName)
	err := renderHTML(r, templatePath, resolvedLayoutPath, data)
	if err != nil {
		log.Println(err)
		r.RenderServerError()
	}
}

func (r Response) RenderFile(templatePath string, data interface{}) {
	r.RenderFileWithLayout(templatePath, "layout", data)
}

func (r Response) RenderFileWithLayout(templatePath, layoutPath string, data interface{}) {
	resolvedLayoutPath := getLayoutPath(layoutPath)
	err := renderHTML(r, templatePath, resolvedLayoutPath, data)
	if err != nil {
		log.Println(err)
		r.RenderServerError()
	}
}

func (r Response) Redirect(url string, req *Request) {
	r.RedirectWithStatus(url, http.StatusFound, req)
}

func (r Response) RedirectWithStatus(url string, code int, req *Request) {
	http.Redirect(r.ResponseWriter, req.HttpRequest, url, code)
}

func (r Response) RenderNotFound() {
	renderStatus(r.ResponseWriter, http.StatusNotFound)
}

func (r Response) RenderBadRequest() {
	renderStatus(r.ResponseWriter, http.StatusBadRequest)
}

func (r Response) RenderServerError() {
	renderStatus(r.ResponseWriter, http.StatusInternalServerError)
}

func renderStatus(rw http.ResponseWriter, code int) {
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, http.StatusText(code))
}
