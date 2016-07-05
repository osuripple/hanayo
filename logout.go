package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func logout(c *gin.Context) {
	ctx := c.MustGet("context").(context)
	if ctx.User.ID == 0 {
		resp(c, 200, "empty.html", &baseTemplateData{TitleBar: "Log out", Messages: []message{warningMessage{"You're already logged out!"}}})
		return
	}
	sess := c.MustGet("session").(sessions.Session)
	sess.Clear()
	sess.Save()
	addMessage(c, successMessage{"Successfully logged out."})
	c.Redirect(302, "/")
}
