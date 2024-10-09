package config

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateManager struct {
	templates *template.Template
}

func NewTemplateManager() error {
	tmpl, err := template.ParseGlob(filepath.Join(TEMPLATE_DIR, "*.html"))
	if err != nil {
		return err
	}
	TMPL = &TemplateManager{templates: tmpl}
	return nil
}

func (tm *TemplateManager) Render(w http.ResponseWriter, tmpl string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return tm.templates.ExecuteTemplate(w, tmpl, data)
}
