package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type profileData struct {
	baseTemplateData
	UserID   int
	Username string
	Email    string
}

func userProfile(c *gin.Context) {
	var field = "id"
	u := c.Param("user")
	if _, err := strconv.Atoi(u); err != nil {
		field = "username"
	}

	data := new(profileData)
	err := db.QueryRow(`SELECT id, username, email FROM users WHERE `+field+` = ?`, u).Scan(
		&data.UserID, &data.Username, &data.Email,
	)
	if err != nil && err != sql.ErrNoRows {
		c.Error(err)
		resp500(c)
		return
	}

	defer resp(c, 200, "profile.html", data)

	if data.UserID == 0 {
		data.TitleBar = "User not found"
		return
	}

	data.TitleBar = data.Username + "'s profile"
	data.HeadingTitle = fmt.Sprintf("<div class='user profile heading'><img src='%s/%d' class='avatar'><span>%s</span>", config.AvatarURL, data.UserID, data.Username)
}
