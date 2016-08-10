package main

import "github.com/gin-gonic/gin"

type simplePage struct {
	Handler, Template, TitleBar, KyutGrill string
}

var simplePages = [...]simplePage{
	{"/login", "login.html", "Log in", "login.png"},
	{"/", "homepage.html", "Home Page", "homepage.jpg"},
	{"/settings/avatar", "settings/avatar.html", "Change avatar", ""},
}

var additionalJS = map[string][]string{
	"/settings/avatar": []string{"/static/croppie.min.js"},
}

func loadSimplePages(r *gin.Engine) {
	for _, el := range simplePages {
		r.GET(el.Handler, simplePageFunc(el))
	}
}

func simplePageFunc(p simplePage) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp(c, 200, p.Template, &baseTemplateData{
			TitleBar:  p.TitleBar,
			KyutGrill: p.KyutGrill,
			Scripts:   additionalJS[p.Handler],
		})
	}
}
