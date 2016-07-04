package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path"

	"gopkg.in/fsnotify.v1"

	"git.zxq.co/ripple/schiavolib"
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
			if i.Name() == path.Base(j) {
				c = true
				break
			}
		}
		if c {
			continue
		}

		fm := template.FuncMap{
			"html": func(value interface{}) template.HTML {
				return template.HTML(fmt.Sprint(value))
			},
		}

		// add new template to template slice
		templates[i.Name()] = template.Must(template.ParseFiles(
			append(baseTemplates[:], "templates/"+i.Name())...,
		)).Funcs(fm)
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
	if corrected, ok := data.(page); ok {
		corrected.SetMessages(getMessages(c))
		corrected.SetPath(c.Request.URL.Path)
		corrected.SetContext(c.MustGet("context").(context))
	}
	c.Status(statusCode)
	err := t.ExecuteTemplate(c.Writer, "base", data)
	if err != nil {
		c.Writer.WriteString("What on earth? Please tell this to a dev!")
		fmt.Println(err)
		schiavo.Bunker.Send(err.Error())
	}
}

type baseTemplateData struct {
	TitleBar  string
	Scripts   []string
	KyutGrill string
	Context   context
	Path      string
	Messages  []message
	FormData  map[string]string
}

func (b *baseTemplateData) SetMessages(m []message) {
	b.Messages = append(b.Messages, m...)
}
func (b *baseTemplateData) SetPath(path string) {
	b.Path = path
}
func (b *baseTemplateData) SetContext(c context) {
	b.Context = c
}

type page interface {
	SetMessages([]message)
	SetPath(string)
	SetContext(context)
}

func reloader() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = w.Add("templates")
	if err != nil {
		w.Close()
		return err
	}
	go func() {
		for {
			select {
			case ev := <-w.Events:
				if ev.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Change detected! Refreshing templates")
					loadTemplates()
				}
			case err := <-w.Errors:
				fmt.Println(err)
			}
		}
	}()
	return nil
}
