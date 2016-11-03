package main

import (
	"errors"
	"net/url"

	"github.com/gin-gonic/gin"
)

func loadSimplePages(r *gin.Engine) {
	for _, el := range simplePages {
		r.GET(el.Handler, simplePageFunc(el))
	}
}

func simplePageFunc(p templateConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := getContext(c)
		if s.User.Privileges&p.mp() != p.mp() {
			resp(c, 200, "empty.html", &baseTemplateData{TitleBar: "Forbidden", Messages: []message{warningMessage{"You should not be 'round here."}}})
			return
		}
		resp(c, 200, p.Template, &baseTemplateData{
			TitleBar:       p.TitleBar,
			KyutGrill:      p.KyutGrill,
			Scripts:        p.additionalJS(),
			HeadingOnRight: p.HugeHeadingRight,
		})
	}
}

func simpleReply(c *gin.Context, errs ...message) error {
	var chosen templateConfig
	for _, s := range simplePages {
		if s.Handler == c.Request.URL.Path {
			chosen = s
		}
	}
	if chosen.Handler == "" {
		return errors.New("simpleReply: simplepage not found")
	}
	resp(c, 200, chosen.Template, &baseTemplateData{
		TitleBar:       chosen.TitleBar,
		KyutGrill:      chosen.KyutGrill,
		Scripts:        chosen.additionalJS(),
		HeadingOnRight: chosen.HugeHeadingRight,
		FormData:       normaliseURLValues(c.Request.PostForm),
		Messages:       errs,
	})
	return nil
}

func normaliseURLValues(uv url.Values) map[string]string {
	m := make(map[string]string, len(uv))
	for k, v := range uv {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}
