package fire

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type pageData struct {
	Page
	Content template.HTML
}

var (
	_basePath  string
	_templates map[string]*template.Template
)

func init() {
	_basePath = filepath.Join("app", "templates")
	_templates = make(map[string]*template.Template)
	cacheTemplates()
}

func cacheTemplates() {
	if filenames, err := filepath.Glob(filepath.Join(_basePath, "*", "*.html")); err != nil {
		panic(err)
	} else {
		for _, filename := range filenames {
			if b, err := ioutil.ReadFile(filename); err != nil {
				panic(err)
			} else {
				t := string(b)
				if tmpl, err := template.New(filename).Parse(t); err != nil {
					panic(err)
				} else {
					_templates[filename] = tmpl
				}
			}
		}
	}
}

func renderHTML(resp Response, templatePath, layoutPath string, p Page) (err error) {
	var templateBuffer bytes.Buffer
	if tmpl := _templates[filepath.Join(_basePath, templatePath)]; tmpl != nil {
		err = tmpl.Execute(&templateBuffer, p)
		if err != nil {
			return
		}
	}

	pdata := pageData{Page: p, Content: template.HTML(templateBuffer.String())}

	var layoutBuffer bytes.Buffer
	if tmpl := _templates[filepath.Join(_basePath, layoutPath)]; tmpl != nil {
		err = tmpl.Execute(&layoutBuffer, pdata)
		if err != nil {
			return
		}
	}

	io.WriteString(resp, layoutBuffer.String())
	return
}

func getLayoutPath(name string) string {
	return filepath.Join("layouts", addHTMLExtension(name))
}

func getTemplatePath(controllerName, name string) string {
	return filepath.Join(controllerName, addHTMLExtension(name))
}

func addHTMLExtension(filename string) string {
	if strings.HasSuffix(filename, ".html") == false {
		return filename + ".html"
	} else {
		return filename
	}
}
