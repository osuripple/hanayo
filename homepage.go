package main

import "github.com/gin-gonic/gin"

func homePage(c *gin.Context) {
	resp(c, 200, "homepage.html", baseTemplateData{
		TitleBar: "Home Page",
	})
}
