package main

import "github.com/gin-gonic/gin"

func testHandler(c *gin.Context) {
	resp(c, 200, "test.html", baseTemplateData{
		TitleBar: "Home Page",
	})
}
