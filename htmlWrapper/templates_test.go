package htmlWrapper

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Borrowed from text/template
// https://golang.org/src/text/template/examplefiles_test.go

// templateFile defines the contents of a template to be stored in a file, for testing.
type templateFile struct {
	name     string
	contents string
}

var testingTemplateFiles = []templateFile{
	// We have a lot of pages almost each using a different parent layout
	// One of the pages is using a subfolder
	{"pages/home.html", `{{define "content"}}home {{.Name}}{{end}} {{/* use one */}}`},
	{"pages/about.html", `{{define "content"}}about {{.Name}}{{end}}{{/* use two */}}`},
	{"pages/third.html", `{{define "content"}}about {{.Name}}{{end}}{{/* use three */}}`},
	{"pages/fourth.html", `{{define "content"}}about {{.Name}}<test-three a="test"><test-two b="a"></test-two></test-three>{{end}} {{/* use one */}}`},
	{"pages/lol.html", `{{define "content"}}about {{.Name}}<header-button></header-button>{{end}} {{/* use one */}}`},
	// We have two different layouts (using two different styles)
	{"layouts/one.html", `Layout 1: {{.Name}} {{block "content" .}}{{end}} {{block "includes/sidebar" .}}{{end}}`},
	{"layouts/two.html", `Layout 2: {{.Name}} {{template "content" .}} {{template "includes/sidebar" .}} <test-one></test-one>`},
	{"layouts/three.html", `Layout 2: {{.Name}} {{template "content" .}} {{template "includes/sidebar" .}} <test-two a="test"></test-two> <test-two a="test" b="{{.Name}}"></test-two >`},
	// We have two includes shared among the pages
	{"includes/header.html", `header`},
	{"includes/sidebar.html", `sidebar {{.Name}}`},
	//For testing html elements
	{"elements/test.html", `<!-- test-one: --> testText`},
	{"elements/testTwo.html", `<!-- test-two: a,b --> testText#b# und #b#`},
	{"elements/testThree.html", `<!-- test-three: a --> test#a##content#`},
	{"elements/testlast.html", `<!-- header-button: --> <fincal-test a="12" b="abcde"> <test2></test2> </fincal-test>
<!-- test2: --> bazinga
<!-- fincal-test: a,b --> button hidden#a#trolling#b# #content#`},
}

func createTestDir(files []templateFile) (dir string, err error) {
	dir, err = os.MkdirTemp("", "template")
	if err != nil {
		return
	}
	for _, file := range files {

		// Create sub directory of file (if needed)
		fd := filepath.Dir(filepath.Join(dir, file.name))
		err = os.MkdirAll(fd, os.ModePerm)
		if err != nil {
			return
		}

		var f *os.File
		f, err = os.Create(filepath.Join(dir, file.name))
		if err != nil {
			return
		}
		_, err = io.WriteString(f, file.contents)
		if err != nil {
			return
		}
		err = f.Close()
		if err != nil {
			return
		}
	}
	return
}

//
// Tests
//

func TestTemplates(t *testing.T) {
	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)
	assert.Nil(t, err)

	// Clean up after the test; another quirk of running as an example.
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatalln(err)
		}
	}(dir)

	templates, err := New(dir, ".html", nil)
	assert.Nil(t, err)

	//
	// 1: Test page + include + layout
	//
	data := struct{ Name string }{"John"}
	b, err := templates.Compile("home", data)
	assert.Nil(t, err)

	want := "Layout 1: John home John sidebar John"
	assert.Equal(t, want, string(b.Bytes()))

	//
	// 2: Test layout isolation
	//
	data = struct{ Name string }{"Jane"}
	b, err = templates.Compile("about", data)
	assert.Nil(t, err)

	want = "Layout 2: Jane about Jane sidebar Jane testText"
	assert.Equal(t, want, string(b.Bytes()))

	//
	// 3: New ReplacingEngine
	//

	data = struct{ Name string }{"Jane"}
	b, err = templates.Compile("third", data)
	assert.Nil(t, err)

	want = "Layout 2: Jane about Jane sidebar Jane testText und  testTextJane und Jane"
	assert.Equal(t, want, string(b.Bytes()))

	//
	// 3.1: New ReplacingEngine
	// no attributes and tags in tests
	//

	data = struct{ Name string }{"Jane"}
	b, err = templates.Compile("fourth", data)
	assert.Nil(t, err)

	want = "Layout 1: Jane about JanetesttesttestTexta und a sidebar Jane"
	assert.Equal(t, want, string(b.Bytes()))

	//
	// 3.2: Special subreplacing Engine behavoir
	//

	data = struct{ Name string }{"Jane"}
	b, err = templates.Compile("lol", data)
	assert.Nil(t, err)

	want = "Layout 1: Jane about Janebutton hidden12trollingabcde bazinga sidebar Jane"
	assert.Equal(t, want, string(b.Bytes()))

	//
	// 4: Test HTTP handler
	//
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data = struct{ Name string }{"Bob"}
		err := templates.Render(w, "home", data, http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	router.ServeHTTP(rr, req)

	want = "Layout 1: Bob home Bob sidebar Bob"
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, want, rr.Body.String())
}

func BenchmarkCompile(b *testing.B) {

	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)
	assert.Nil(b, err)

	// Clean up after the test; another quirk of running as an example.
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatalln(err)
		}
	}(dir)

	templates, err := New(dir, ".html", DefaultFunctions)
	assert.Nil(b, err)

	b.ResetTimer()

	data := struct{ Name string }{"John"}

	for i := 0; i < b.N; i++ {

		body, err := templates.Compile("home", data)
		assert.Nil(b, err)

		want := "Layout 1: John home John sidebar John"
		assert.Equal(b, want, string(body.Bytes()))
	}
}

/*
func BenchmarkCompileBuffer(b *testing.B) {

	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)

	if err != nil {
		b.Error(err)
	}

	// Clean up after the test; another quirk of running as an example.
	defer os.RemoveAll(dir)

	templates, err := New(dir, ".html")

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	data := struct{ Name string }{"John"}

	for i := 0; i < b.N; i++ {

		buf := bufpool.GetAll()

		err := templates.CompileWithBuffer("home", data, buf)
		if err != nil {
			b.Error(err)
		}

		got := string(buf.Bytes())
		want := "Layout 1: John home John sidebar John"

		bufpool.Put(buf)

		if got != want {
			b.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
		}
	}
}
*/

func BenchmarkNativeTemplates(b *testing.B) {

	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)

	if err != nil {
		b.Error(err)
	}

	// Clean up after the test; another quirk of running as an example.
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatalln(err)
		}
	}(dir)

	t := template.New("")

	var by []byte
	for _, name := range []string{"pages/home/home", "layouts/one", "includes/sidebar"} {
		by, err = os.ReadFile(filepath.Join(dir, name+".html"))
		if err != nil {
			b.Error(err)
		}
		_, err = t.New(name).Parse(string(by))
		if err != nil {
			b.Error(err)
		}
	}

	b.ResetTimer()

	data := struct{ Name string }{"John"}

	for i := 0; i < b.N; i++ {

		by = nil
		buf := bytes.NewBuffer(by)

		if err := t.ExecuteTemplate(buf, "layouts/one", data); err != nil {
			b.Error(err)
		}

		got := string(buf.Bytes())
		want := "Layout 1: John home John sidebar John"

		if got != want {
			b.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
		}
	}
}
