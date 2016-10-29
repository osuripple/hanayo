package main

import (
	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	if c.Query("stopsign") != "1" {
		u := tryBotnets(c)
		if u != "" {
			resp(c, 200, "elmo.html", passwordResetContinueTData{
				Username: u,
				baseTemplateData: baseTemplateData{
					TitleBar:       "Elmo! Stop!",
					HeadingTitle:   "Stop!",
					KyutGrill:      "stop_sign.png",
					HeadingOnRight: true,
				},
			})
			return
		}
	}
}

func tryBotnets(c *gin.Context) string {
	var username string

	err := db.QueryRow("SELECT u.username FROM ip_user i LEFT JOIN users u ON u.id = i.userid WHERE i.ip = ?", clientIP(c)).Scan(&username)
	if err != nil {
		c.Error(err)
		return ""
	}
	if username != "" {
		return username
	}

	cook, _ := c.Cookie("y")
	err = db.QueryRow("SELECT u.username FROM identity_tokens i LEFT JOIN users u ON u.id = i.userid WHERE i.token = ?",
		cook).Scan(&username)
	if err != nil {
		c.Error(err)
		return ""
	}
	return username
}
