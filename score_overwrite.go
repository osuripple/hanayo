package main

import (
	"github.com/gin-gonic/gin"
)

func scoreOverwrite(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}
	var m message = successMessage{T(c, "Your score overwrite preferences have been saved.")}
	defer func() {
		addMessage(c, m)
		getSession(c).Save()
		c.Redirect(302, "/settings/score_overwrite")
	}()
	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
		m = errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")}
		return
	}
	arr := c.PostFormArray("overwrite")
	if len(arr) != 4 {
		m = errorMessage{T(c, "An error occurred.")}
		return
	}
}
