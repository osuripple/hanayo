package main

import "github.com/gin-gonic/gin"

type homePageData struct {
	baseTemplateData
	Posts []BlogPost
}

func homePage(c *gin.Context) {
	posts := getBlogPosts(5)
	resp(c, 200, "homepage.html", homePageData{
		baseTemplateData: baseTemplateData{
			TitleBar:  "Home Page",
			KyutGrill: "homepage.jpg",
		},
		Posts: posts,
	})
}
