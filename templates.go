package main

import (
	"html/template"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

var templates map[string]*template.Template
var baseTemplates = [...]string{
	"templates/base.html",
	"templates/navbar.html",
}

func loadTemplates() {
	ts, err := ioutil.ReadDir("templates")
	if err != nil {
		panic(err)
	}

	templates = make(map[string]*template.Template, len(ts)-len(baseTemplates))

	for _, i := range ts {
		// make sure it's not a directory
		if i.IsDir() {
			continue
		}

		// do not compile base templates on their own
		var c bool
		for _, j := range baseTemplates {
			if i.Name() == j {
				c = true
				break
			}
		}
		if c {
			continue
		}

		// add new template to template slice
		templates[i.Name()] = template.Must(template.ParseFiles(
			append(baseTemplates[:], "templates/"+i.Name())...,
		))
	}
}

func resp(c *gin.Context, statusCode int, tpl string, data interface{}) {
	if c == nil {
		return
	}
	t := templates[tpl]
	if t == nil {
		c.String(500, "Template not found! Please tell this to a dev!")
		return
	}
	c.Status(statusCode)
	err := t.ExecuteTemplate(c.Writer, "base", data)
	if err != nil {
		c.Writer.WriteString("What on earth? Please tell this to a dev!")
	}
}

type baseTemplateData struct {
	TitleBar string
}
