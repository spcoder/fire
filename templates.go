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

func init() {
	_templates = make(map[string]*template.Template)
	cacheTemplates()
}

func cacheTemplates() {
	if filenames, err := filepath.Glob(filepath.Join("templates", "*", "*.html")); err != nil {
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

func renderTemplate(resp Response, controllerName, templateName, layoutName, keywords, description string, data interface{}) (err error) {
	layoutPath := getLayoutPath(layoutName)
	var layoutBuffer bytes.Buffer
	if tmpl := _templates[layoutPath]; tmpl != nil {
		err = _templates[layoutPath].Execute(&layoutBuffer, nil)
		if err != nil {
			return
		}
	}

	templatePath := getTemplatePath(controllerName, templateName)
	var templateBuffer bytes.Buffer
	if tmpl := _templates[templatePath]; tmpl != nil {
		err = _templates[templatePath].Execute(&templateBuffer, data)
		if err != nil {
			return
		}
	}

	html := replaceTokens(layoutBuffer, templateBuffer, keywords, description)
	io.WriteString(resp, html)
	return
}

func replaceTokens(layout, template bytes.Buffer, keywords, description string) (html string) {
	html = strings.Replace(layout.String(), "%content%", template.String(), -1)
	html = strings.Replace(html, "%keywords%", keywords, -1)
	html = strings.Replace(html, "%description%", description, -1)
	return
}

func getByteBuffer(path string, b *bytes.Buffer, data interface{}) (err error) {
	return _templates[path].Execute(b, data)
}

func getLayoutPath(name string) string {
	return filepath.Join("templates", "layouts", addHTMLExtension(name))
}

func getTemplatePath(controllerName, name string) string {
	return filepath.Join("templates", controllerName, addHTMLExtension(name))
}

func addHTMLExtension(filename string) string {
	if strings.HasSuffix(filename, ".html") == false {
		return filename + ".html"
	} else {
		return filename
	}
}
