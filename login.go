package main

import "github.com/gin-gonic/gin"

func login(c *gin.Context) {
	resp(c, 200, "login.html", &baseTemplateData{
		TitleBar:  "Log in",
		KyutGrill: "login.png",
		Path:      c.Request.URL.Path,
	})
}
