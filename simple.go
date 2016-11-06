package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

func loadSimplePages(r *gin.Engine) {
	for _, el := range simplePages {
		if el.Handler == "" {
			continue
		}
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
		simple(c, p, nil, nil)
	}
}

func simpleReply(c *gin.Context, errs ...message) error {
	t := getSimple(c.Request.URL.Path)
	if t.Handler == "" {
		return errors.New("simpleReply: simplepage not found")
	}
	simple(c, t, errs, nil)
	return nil
}

func getSimple(h string) templateConfig {
	for _, s := range simplePages {
		if s.Handler == h {
			return s
		}
	}
	fmt.Println("oh handler shit not found", h)
	return templateConfig{}
}

func getSimpleByFilename(f string) templateConfig {
	for _, s := range simplePages {
		if s.Template == f {
			return s
		}
	}
	fmt.Println("oh shit not found", f)
	return templateConfig{}
}

func simple(c *gin.Context, p templateConfig, errs []message, requestInfo map[string]interface{}) {
	resp(c, 200, p.Template, &baseTemplateData{
		TitleBar:       p.TitleBar,
		KyutGrill:      p.KyutGrill,
		Scripts:        p.additionalJS(),
		HeadingOnRight: p.HugeHeadingRight,
		FormData:       normaliseURLValues(c.Request.PostForm),
		RequestInfo:    requestInfo,
		Messages:       errs,
	})
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
