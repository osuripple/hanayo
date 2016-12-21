package http

import (
	"html/template"

	"git.zxq.co/ripple/hanayo/tpl"
	"github.com/julienschmidt/httprouter"
)

var (
	// Templates is a map containing the template files.
	Templates map[string]*template.Template
	// SimplePages are templates which automatically have an handler set up.
	SimplePages []tpl.TemplateConfig
)

// SetUpSimplePages sets up simplepages.
func (s *Server) SetUpSimplePages() error {
	if s.Router == nil {
		s.Router = httprouter.New()
	}
	var err error
	Templates, SimplePages, err = tpl.LoadTemplates()
	return err
}
