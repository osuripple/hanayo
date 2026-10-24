package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"zxq.co/ripple/rippleapi/common"
)

func ircGenToken(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	var s, m string
	for {
		var err error
		s, err = common.RandomString(32)
		if err != nil {
			c.Error(err)
			resp500(c)
			return
		}
		m = cmd5(s)
		if db.QueryRow("SELECT 1 FROM irc_tokens WHERE token = ? LIMIT 1", m).
			Scan(new(int)) == sql.ErrNoRows {
			break
		}
	}

	db.Exec("DELETE FROM irc_tokens WHERE userid = ?", ctx.User.ID)
	db.Exec("INSERT INTO irc_tokens(userid, token) VALUES (?, ?)", ctx.User.ID, m)
	simple(c, getSimple("/irc"), []message{successMessage{
		T(c, "Your new IRC token is <code>%s</code>. The old IRC token is not valid anymore.<br>Keep it safe, don't show it around, and store it now! We won't show it to you again.", s),
	}}, nil)
}
