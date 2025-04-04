package utils

import (
	"bytes"
	"html/template"
	"path/filepath"
)

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func RenderTemplateToString(templatePath string, datasoal, datakunci any) (string, error) {
	html := template.Must(template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"safeHTML": safeHTML,
	}).ParseFiles(templatePath))

	var buf bytes.Buffer
	err := html.Execute(&buf, map[string]any{
		"datasoal":  datasoal,
		"datakunci": datakunci,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
