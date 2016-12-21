package funcmap

import (
	"fmt"
	"html/template"
)

// HTML disables HTML-escaping on a certain element.
func HTML(value interface{}) template.HTML {
	return template.HTML(fmt.Sprint(value))
}
