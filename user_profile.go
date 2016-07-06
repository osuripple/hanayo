package main

import (
	"database/sql"
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

	resp(c, 200, "profile.html", data)
}
