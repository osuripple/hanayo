package main

import "github.com/gin-gonic/gin"

func notFound(c *gin.Context) {
	resp(c, 200, "not_found.html", baseTemplateData{
		TitleBar:  "Not Found",
		KyutGrill: "not_found.jpg",
	})
}
