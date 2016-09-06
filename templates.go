package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rjeczalik/notify"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pariz/gountries"
)

var templates = make(map[string]*template.Template)
var baseTemplates = [...]string{
	"templates/base.html",
	"templates/navbar.html",
}

var gdb = gountries.New()

func countryReadable(s string) string {
	if s == "XX" || s == "" {
		return ""
	}
	reg, err := gdb.FindCountryByAlpha(s)
	if err != nil {
		return ""
	}
	return reg.Name.Common
}

func loadTemplates(subdir string) {
	ts, err := ioutil.ReadDir("templates" + subdir)
	if err != nil {
		panic(err)
	}

	for _, i := range ts {
		// if it's a directory, load recursively
		if i.IsDir() && i.Name() != ".." && i.Name() != "." {
			loadTemplates(subdir + "/" + i.Name())
			continue
		}

		// do not compile base templates on their own
		var comp bool
		for _, j := range baseTemplates {
			if "templates"+subdir+"/"+i.Name() == j {
				comp = true
				break
			}
		}
		if comp {
			continue
		}

		var inName string
		if subdir != "" && subdir[0] == '/' {
			inName = subdir[1:] + "/"
		}

		// add new template to template slice
		templates[inName+i.Name()] = template.Must(template.New(i.Name()).Funcs(funcMap).ParseFiles(
			append([]string{"templates" + subdir + "/" + i.Name()}, baseTemplates[:]...)...,
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
	sess := c.MustGet("session").(sessions.Session)
	if corrected, ok := data.(page); ok {
		corrected.SetMessages(getMessages(c))
		corrected.SetPath(c.Request.URL.Path)
		corrected.SetContext(c.MustGet("context").(context))
		corrected.SetGinContext(c)
		corrected.SetSession(sess)
	}
	sess.Save()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		c.Writer.WriteString(
			"oooops! A brit monkey stumbled upon a banana while trying to process your request. " +
				"This doesn't make much sense, but in a few words: we fucked up something while processing your " +
				"request. We are sorry for this, but don't worry: we have been notified and are on it!",
		)
		c.Error(err)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(statusCode)
	_, err = io.Copy(c.Writer, buf)
	if err != nil {
		c.Writer.WriteString("We don't know what's happening now.")
		c.Error(err)
		return
	}
}

type baseTemplateData struct {
	TitleBar       string
	HeadingTitle   string
	HeadingOnRight bool
	Scripts        []string
	KyutGrill      string
	DisableHH      bool // HH = Huge Heading
	Context        context
	Path           string
	Messages       []message
	FormData       map[string]string
	Gin            *gin.Context
	Session        sessions.Session
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
func (b *baseTemplateData) SetGinContext(c *gin.Context) {
	b.Gin = c
}
func (b *baseTemplateData) SetSession(sess sessions.Session) {
	b.Session = sess
}
func (b baseTemplateData) Get(s string, params ...interface{}) map[string]interface{} {
	s = fmt.Sprintf(s, params...)
	req, err := http.NewRequest("GET", config.API+s, nil)
	if err != nil {
		b.Gin.Error(err)
		return nil
	}
	req.Header.Set("User-Agent", "hanayo")
	req.Header.Set("H-Key", config.APISecret)
	req.Header.Set("X-Ripple-Token", b.Context.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.Gin.Error(err)
		return nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		b.Gin.Error(err)
		return nil
	}
	x := make(map[string]interface{})
	err = json.Unmarshal(data, &x)
	if err != nil {
		b.Gin.Error(err)
		return nil
	}
	return x
}
func (b baseTemplateData) Has(privs uint64) bool {
	return uint64(b.Context.User.Privileges)&privs == privs
}

type page interface {
	SetMessages([]message)
	SetPath(string)
	SetContext(context)
	SetGinContext(*gin.Context)
	SetSession(sessions.Session)
}

func reloader() error {
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch("./templates/...", c, notify.All); err != nil {
		return err
	}
	go func() {
		for range c {
			fmt.Println("Change detected! Refreshing templates")
			loadTemplates("")
		}
		defer notify.Stop(c)
	}()
	return nil
}
