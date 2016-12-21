package tpl

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"git.zxq.co/ripple/hanayo/tpl/funcmap"
)

// LoadTemplates loads all the templates, and returns simplepages.
func LoadTemplates() (map[string]*template.Template, []TemplateConfig, error) {
	l := &templateLoader{
		folderBase: "templates",
		templates:  make(map[string]*template.Template),
		baseTemplates: []string{
			"templates/base.html",
			"templates/simplepag.html",
			"templates/navbar.html",
		},
	}
	err := l.loadTemplates("/.")
	return l.templates, l.simplePages, err
}

type templateLoader struct {
	folderBase    string
	templates     map[string]*template.Template
	baseTemplates []string
	simplePages   []TemplateConfig
}

func (t *templateLoader) loadTemplates(subdir string) error {
	ts, err := ioutil.ReadDir(t.folderBase + subdir)
	if err != nil {
		return err
	}

TemplateFileLooper:
	for _, i := range ts {
		// if it's a directory, load recursively
		if i.IsDir() && i.Name() != ".." && i.Name() != "." {
			err := t.loadTemplates(subdir + "/" + i.Name())
			if err != nil {
				return err
			}
			continue
		}

		fullName := t.folderBase + subdir + "/" + i.Name()
		c, err := t.parseConfig(fullName)
		if err != nil {
			return err
		}

		if c.NoCompile {
			continue
		}

		var files = c.IncludedFiles(t.folderBase + subdir + "/")
		files = append(files, fullName)

		// do not compile base templates on their own
		for _, j := range t.baseTemplates {
			if fullName == j {
				continue TemplateFileLooper
			}
		}

		var inName string
		if subdir != "" && subdir[0] == '/' {
			inName = subdir[1:] + "/"
		}

		nt, err := template.New(i.Name()).Funcs(funcmap.FuncMap).ParseFiles(
			append(files, t.baseTemplates...)...,
		)
		if err != nil {
			return err
		}

		// add new template to template slice
		t.templates[inName+i.Name()] = nt

		if c.Handler != "" {
			t.simplePages = append(t.simplePages, c)
		}
	}
	return nil
}

func (t *templateLoader) parseConfig(file string) (TemplateConfig, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return TemplateConfig{}, err
	}
	i := bufio.NewScanner(f)
	var inConfig bool
	var buff string
	for i.Scan() {
		u := i.Text()
		switch u {
		case "{{/*###":
			inConfig = true
		case "*/}}":
			if !inConfig {
				continue
			}
			tr, err := LoadConf(buff)
			tr.Template = strings.TrimPrefix(file, t.folderBase+"/")
			return tr, err
		}
		if !inConfig {
			continue
		}
		buff += u + "\n"
	}
	return TemplateConfig{}, nil
}
