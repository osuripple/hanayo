package main

import (
	"errors"
	"net/url"

	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/gin"
)

type simplePage struct {
	Handler, Template, TitleBar, KyutGrill string
	MinPrivileges                          common.UserPrivileges
}

var simplePages = [...]simplePage{
	{"/", "homepage.html", "Home Page", "homepage.jpg", 0},
	{"/login", "login.html", "Log in", "login.png", 0},
	{"/settings/avatar", "settings/avatar.html", "Change avatar", "settings.png", 2},
	{"/dev/tokens", "dev/tokens.html", "Your API tokens", "dev.png", 2},
	{"/beatmaps/rank_request", "beatmaps/rank_request.html", "Request beatmap ranking", "request_beatmap_ranking.jpg", 2},
	{"/donate", "support.html", "Support Ripple", "donate.jpg", 0},
	{"/doc", "doc.html", "Documentation", "documentation.jpg", 0},
	{"/doc/:id", "doc_content.html", "View document", "documentation.jpg", 0},
	{"/help", "help.html", "Contact support", "help.jpg", 0},
	{"/leaderboard", "leaderboard.html", "Leaderboard", "leaderboard.jpg", 0},
	{"/friends", "friends.html", "Friends", "", 2},
	{"/changelog", "changelog.html", "Changelog", "", 0},
	{"/team", "team.html", "Team", "", 0},
	{"/pwreset", "pwreset.html", "Reset password", "", 0},
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
		s := getContext(c)
		if s.User.Privileges&p.MinPrivileges != p.MinPrivileges {
			resp(c, 200, "empty.html", &baseTemplateData{TitleBar: "Forbidden", Messages: []message{warningMessage{"You should not be 'round here."}}})
			return
		}
		resp(c, 200, p.Template, &baseTemplateData{
			TitleBar:       p.TitleBar,
			KyutGrill:      p.KyutGrill,
			Scripts:        additionalJS[p.Handler],
			HeadingOnRight: hhr,
		})
	}
}

func simpleReply(c *gin.Context, errs ...message) error {
	var chosen simplePage
	for _, s := range simplePages {
		if s.Handler == c.Request.URL.Path {
			chosen = s
		}
	}
	if chosen.Handler == "" {
		return errors.New("simpleReply: simplepage not found")
	}
	resp(c, 200, chosen.Template, &baseTemplateData{
		TitleBar:  chosen.TitleBar,
		KyutGrill: chosen.KyutGrill,
		Scripts:   additionalJS[chosen.Handler],
		FormData:  normaliseURLValues(c.Request.PostForm),
		Messages:  errs,
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
