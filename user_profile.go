package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type profileData struct {
	baseTemplateData
	UserID int
}

func userProfile(c *gin.Context) {
	var (
		userID   int
		username string
	)

	u := c.Param("user")
	if _, err := strconv.Atoi(u); err != nil {
		err := db.QueryRow("SELECT id, username FROM users WHERE username = ? LIMIT 1", u).Scan(&userID, &username)
		if err != nil && err != sql.ErrNoRows {
			c.Error(err)
		}
	} else {
		err := db.QueryRow(`SELECT id, username FROM users WHERE id = ? LIMIT 1`, u).Scan(&userID, &username)
		switch {
		case err == nil:
		case err == sql.ErrNoRows:
			err := db.QueryRow(`SELECT id, username FROM users WHERE username = ? LIMIT 1`, u).Scan(&userID, &username)
			if err != nil && err != sql.ErrNoRows {
				c.Error(err)
			}
		default:
			c.Error(err)
		}
	}

	data := new(profileData)
	data.UserID = userID

	defer resp(c, 200, "profile.html", data)

	if data.UserID == 0 {
		data.TitleBar = "User not found"
		data.Messages = append(data.Messages, warningMessage{"That user could not be found!"})
		return
	}

	data.TitleBar = username + "'s profile"
	data.HeadingTitle = fmt.Sprintf("<div class='user profile heading'><img src='%s/%d' class='avatar'><span>%s</span>", config.AvatarURL, userID, username)
}
