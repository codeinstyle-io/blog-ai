package utils

import (
	"encoding/json"
	"html/template"
)

// GetTemplateFuncs returns the common template functions used across the application
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
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
				return "[]"
			}
			return string(b)
		},
	}
}
