package main

import "github.com/gin-gonic/gin"

type simplePage struct {
	Handler, Template, TitleBar, KyutGrill string
}

var simplePages = [...]simplePage{
	{"/", "homepage.html", "Home Page", "homepage.jpg"},
	{"/login", "login.html", "Log in", "login.png"},
	{"/settings/avatar", "settings/avatar.html", "Change avatar", ""},
	{"/dev/tokens", "dev/tokens.html", "Your API tokens", "dev.png"},
	{"/beatmaps/rank_request", "beatmaps/rank_request.html", "Request beatmap ranking", ""},
}

// indexes of pages in simplePages that have huge heading on the right
var hugeHeadingRight = [...]int{
	3,
}

var additionalJS = map[string][]string{
	"/settings/avatar": []string{"/static/croppie.min.js"},
}

func loadSimplePages(r *gin.Engine) {
	for i, el := range simplePages {
		// if the page has hugeheading on the right, tell it to the
		// simplePageFunc.
		var right bool
		for _, hhr := range hugeHeadingRight {
			if hhr == i {
				right = true
			}
		}
		r.GET(el.Handler, simplePageFunc(el, right))
	}
}

func simplePageFunc(p simplePage, hhr bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp(c, 200, p.Template, &baseTemplateData{
			TitleBar:       p.TitleBar,
			KyutGrill:      p.KyutGrill,
			Scripts:        additionalJS[p.Handler],
			HeadingOnRight: hhr,
		})
	}
}
