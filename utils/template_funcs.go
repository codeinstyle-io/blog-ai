package utils

import (
	"encoding/json"
	"html/template"
	"strings"
	"time"
)

// GetTemplateFuncs returns the map of template functions used in templates
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("January 2, 2006 15:04")
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"trim":  strings.TrimSpace,
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"raw": func(s string) template.HTML {
			return template.HTML(s)
		},
		"json": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return ""
			}
			return string(b)
		},
	}
}
