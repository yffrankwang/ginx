package ginhtml

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error. In that case, if the
// second (error) argument evaluates to non-nil during execution, execution
// terminates and Execute returns that error. FuncMap has the same base type
// as FuncMap in "text/template", copied here so clients need not import
// "text/template".
type FuncMap map[string]interface{}

// Delims delims for template
type Delims struct {
	Left  string
	Right string
}

// HTMLTemplates html template engine
type HTMLTemplates struct {
	extensions map[string]bool // template extensions
	funcs      FuncMap         // template functions
	delims     Delims          // delimeters

	template *template.Template
}

// NewHTMLTemplates new template engine
func NewHTMLTemplates(extensions ...string) *HTMLTemplates {
	ht := &HTMLTemplates{
		delims: Delims{Left: "{{", Right: "}}"},
	}

	ht.Extensions(extensions...)
	return ht
}

// Extensions sets template entensions.
func (ht *HTMLTemplates) Extensions(extensions ...string) {
	if len(extensions) == 0 {
		extensions = []string{".html", ".gohtml"}
	}

	he := map[string]bool{}
	for _, s := range extensions {
		he[s] = true
	}
	ht.extensions = he
}

// Delims sets template left and right delims and returns a Engine instance.
func (ht *HTMLTemplates) Delims(left, right string) {
	ht.delims = Delims{Left: left, Right: right}
}

// Funcs sets the FuncMap used for template.FuncMap.
func (ht *HTMLTemplates) Funcs(funcMap FuncMap) {
	ht.funcs = funcMap
}

func (ht *HTMLTemplates) init() {
	if ht.template == nil {
		tpl := template.New("")
		tpl.Delims(ht.delims.Left, ht.delims.Right)
		tpl.Funcs(template.FuncMap(ht.funcs))
		ht.template = tpl
	}
}

// Load glob and parse template files under the root path
func (ht *HTMLTemplates) Load(root string) (err error) {
	ht.init()

	root, err = filepath.Abs(root)
	if err != nil {
		return
	}

	root = filepath.ToSlash(root)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return ht.loadFile(nil, root, path)
	})
	if err != nil {
		return
	}

	return
}

// LoadFS glob and parse template files from FS
func (ht *HTMLTemplates) LoadFS(fsys fs.FS, root string) (err error) {
	ht.init()

	err = fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return ht.loadFile(fsys, root, path)
	})
	if err != nil {
		return
	}

	return
}

// loadFile load template file
func (ht *HTMLTemplates) loadFile(fsys fs.FS, root, path string) error {
	ext := filepath.Ext(path)
	if _, ok := ht.extensions[ext]; !ok {
		return nil
	}

	text, err := readFile(fsys, path)
	if err != nil {
		return fmt.Errorf("HTMLTemplates load template %q error: %w", path, err)
	}

	path = toTemplateName(root, path, ext)

	tpl := ht.template.New(path)
	_, err = tpl.Parse(text)
	if err != nil {
		return fmt.Errorf("HTMLTemplates parse template %q error: %w", path, err)
	}
	return nil
}

// Render render template with io.Writer
func (ht *HTMLTemplates) Render(w io.Writer, name string, data interface{}) error {
	err := ht.template.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("HTMLTemplates execute template %q error: %w", name, err)
	}

	return nil
}

// readFile read file content to string
func readFile(fsys fs.FS, path string) (text string, err error) {
	var data []byte
	if fsys == nil {
		data, err = os.ReadFile(path)
	} else {
		data, err = fs.ReadFile(fsys, path)
	}

	if err != nil {
		return "", fmt.Errorf("Failed to read template %v, error: %w", path, err)
	}
	return string(data), nil
}

func toTemplateName(root, path, ext string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, root)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ext)
	return path
}
