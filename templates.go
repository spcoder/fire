package fire

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	_templates map[string]*template.Template
)

type layoutData struct {
	Content template.HTML
	Context interface{}
}

func initTemplates() {
	cacheTemplates()
}

func getBasePath() string {
	return TemplateFileBase()
}

func cacheTemplates() {
	_templates = make(map[string]*template.Template)

	if filenames, err := filepath.Glob(filepath.Join(getBasePath(), "*", "*.html")); err != nil {
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

func renderHTML(resp Response, templatePath, layoutPath string, context interface{}) (err error) {
	if IsDevelopment() {
		cacheTemplates()
	}

	var templateBuffer bytes.Buffer
	if tmpl := _templates[filepath.Join(getBasePath(), templatePath)]; tmpl != nil {
		err = tmpl.Execute(&templateBuffer, context)
		if err != nil {
			return
		}
	}

	context = layoutData{Context: context, Content: template.HTML(templateBuffer.String())}

	var layoutBuffer bytes.Buffer
	if layoutPath == "" {
		layoutBuffer = templateBuffer
	} else {
		if tmpl := _templates[filepath.Join(getBasePath(), layoutPath)]; tmpl != nil {
			err = tmpl.Execute(&layoutBuffer, context)
			if err != nil {
				return
			}
		}
	}

	html := layoutBuffer.String()
	io.WriteString(resp.ResponseWriter, html)
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
