package htmlWrapper

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Children define the base template using comments: { /* use basetemplate */ }
var parentRegex = regexp.MustCompile(`\{\{\s*/\*\s*use\s(\w+)\s*\*/\s*}}`)

// NotFoundError for type assertions in the caller while still providing context
type NotFoundError struct {
	Name string
}

func (t *NotFoundError) Error() string {
	return fmt.Sprintf("template %q not found", t.Name)
}

// Templates Collection
type Templates struct {
	Extension string
	Dir       string
	Templates map[string]*template.Template
}

// New templates collection
func New(templatesDir, extension string, functions template.FuncMap) (*Templates, error) {
	t := &Templates{
		Extension: extension,
		Dir:       templatesDir,
		Templates: make(map[string]*template.Template),
	}

	return t, t.load(functions)
}

// Funcs function map for templates
func (t *Templates) Funcs(functions template.FuncMap) *Templates {
	for _, tmpl := range t.Templates {
		tmpl.Funcs(functions)
	}

	return t
}

// Handles loading required templates
func (t *Templates) load(functions template.FuncMap) (err error) {
	//load Elements into memory
	htmlMap := loadHTMLParts(t.Dir, "elements", t.Extension)
	htmlMap.cleanMap()

	var pages map[string][]byte
	pages, err = loadTemplateFiles(t.Dir, "pages", t.Extension)
	if err != nil {
		return
	}

	// Shared templates across multiple pages (sidebars, scripts, footers, etc...)
	var includes map[string][]byte
	includes, err = loadTemplateFiles(t.Dir, "includes", t.Extension)
	if err != nil {
		return
	}

	// Layouts used by pages
	var layouts map[string][]byte
	layouts, err = loadTemplateFiles(t.Dir, "layouts", t.Extension)
	if err != nil {
		return
	}

	htmlMap.replaceAllElements(&pages, "pages")
	htmlMap.replaceAllElements(&includes, "includes")
	htmlMap.replaceAllElements(&layouts, "layouts")

	var tmpl *template.Template
	for name, b := range pages {

		matches := parentRegex.FindSubmatch(b)
		basename := filepath.Base(name)

		tmpl, err = template.New(basename).Funcs(functions).Parse(string(b))
		if err != nil {
			return
		}

		// Uses a layout
		if len(matches) == 2 {

			l, ok := layouts[filepath.Join("layouts", string(matches[1]))]
			if !ok {
				err = fmt.Errorf("unknown file: layouts/%s%s", matches[1], t.Extension)
				return
			}

			_, err = tmpl.New("layout").Parse(string(l))
			if err != nil {
				return
			}
		}

		if len(includes) > 0 {
			for name, src := range includes {
				name = strings.ReplaceAll(name, string(os.PathSeparator), "/")
				_, err = tmpl.New(name).Parse(string(src))
				if err != nil {
					return
				}
			}
		}

		t.Templates[basename] = tmpl
	}

	return
}

// DefinedTemplates loaded by got (for debugging)
func (t *Templates) DefinedTemplates() (out string) {
	for _, tmpl := range t.Templates {
		out += fmt.Sprintf("%s%s\n", tmpl.Name(), tmpl.DefinedTemplates())
	}
	return
}

// Render the template & data to the ResponseWriter safely
func (t *Templates) Render(w http.ResponseWriter, template string, data interface{}, status int) error {

	buf, err := t.Compile(template, data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

// Compile the template and return the buffer containing the rendered bytes
func (t *Templates) Compile(template string, data interface{}) (*bytes.Buffer, error) {

	// Look for the template
	tmpl, ok := t.Templates[template]

	if !ok {
		err := &NotFoundError{template}
		return nil, err
	}

	// Create a buffer so syntax errors don't return a half-rendered response body
	var b []byte
	buf := bytes.NewBuffer(b)

	if err := tmpl.ExecuteTemplate(buf, "layout", data); err != nil {
		return buf, err
	}

	return buf, nil
}
