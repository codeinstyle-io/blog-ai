package utils

import "text/template"

// GetTemplateFuncs returns the common template functions used across the application
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}
}
