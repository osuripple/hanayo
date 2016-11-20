package main

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO: replace with simple ResponseInfo containing userid
type profileData struct {
	baseTemplateData
	UserID int
}

func userProfile(c *gin.Context) {
	var (
		userID   int
		username string
	)

	ctx := getContext(c)

	u := c.Param("user")
	if _, err := strconv.Atoi(u); err != nil {
		err := db.QueryRow("SELECT id, username FROM users WHERE username = ? AND "+ctx.OnlyUserPublic()+" LIMIT 1", u).Scan(&userID, &username)
		if err != nil && err != sql.ErrNoRows {
			c.Error(err)
		}
	} else {
		err := db.QueryRow(`SELECT id, username FROM users WHERE id = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username)
		switch {
		case err == nil:
		case err == sql.ErrNoRows:
			err := db.QueryRow(`SELECT id, username FROM users WHERE username = ? AND `+ctx.OnlyUserPublic()+` LIMIT 1`, u).Scan(&userID, &username)
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
	data.DisableHH = true
	data.Scripts = append(data.Scripts, "/static/jquery.nanoscroller.min.js", "/static/profile.js")
}
