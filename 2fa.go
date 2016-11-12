package main

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
	"git.zxq.co/x/rs"

	"github.com/gin-gonic/gin"
)

var allowedPaths = [...]string{
	"/logout",
	"/2fa_gateway",
	"/2fa_gateway/verify",
	"/2fa_gateway/clear",
}

// middleware to deny all requests to non-allowed pages
func twoFALock(c *gin.Context) {
	// check it's not a static file
	if len(c.Request.URL.Path) >= 8 && c.Request.URL.Path[:8] == "/static/" {
		c.Next()
		return
	}

	ctx := getContext(c)
	if ctx.User.ID == 0 {
		c.Next()
		return
	}

	sess := getSession(c)
	if v, _ := sess.Get("2fa_must_validate").(bool); !v {
		// * check 2fa is enabled.
		//   if it is,
		//   * check whether the current ip is found in the database.
		//     if it is, move on and show the page.
		//     if it isn't, set 2fa_must_validate
		//   if it isn't, move on.
		enabled := is2faEnabled(ctx.User.ID)
		if enabled {
			err := db.QueryRow("SELECT 1 FROM ip_user WHERE userid = ? AND ip = ? LIMIT 1", ctx.User.ID, clientIP(c)).Scan(new(int))
			if err != sql.ErrNoRows {
				c.Next()
				return
			}
			sess.Set("2fa_must_validate", true)
		} else {
			c.Next()
			return
		}
	}

	// check it's one of the few approved paths
	for _, a := range allowedPaths {
		if a == c.Request.URL.Path {
			sess.Save()
			c.Next()
			return
		}
	}
	addMessage(c, warningMessage{"You need to complete the 2fa challenge first."})
	sess.Save()
	c.Redirect(302, "/2fa_gateway")
	c.Abort()
}

// is2faEnabled checks 2fa is enabled for an user.
func is2faEnabled(user int) bool {
	return db.QueryRow("SELECT 1 FROM 2fa_telegram WHERE userid = ?", user).Scan(new(int)) != sql.ErrNoRows
}

func tfaGateway(c *gin.Context) {
	sess := getSession(c)

	// check 2fa hasn't been disabled
	i, _ := sess.Get("userid").(int)
	if i == 0 {
		c.Redirect(302, "/")
	}
	if !is2faEnabled(i) {
		sess.Delete("2fa_must_validate")
		c.Redirect(302, "/")
		return
	}
	// check previous 2fa thing is still valid
	err := db.QueryRow("SELECT 1 FROM 2fa WHERE userid = ? AND ip = ? AND expire > ?",
		i, clientIP(c), time.Now().Unix()).Scan(new(int))
	if err != nil {
		db.Exec("INSERT INTO 2fa(userid, token, ip, expire, sent) VALUES (?, ?, ?, ?, 0);",
			i, strings.ToUpper(rs.String(8)), clientIP(c), time.Now().Add(time.Hour).Unix())
		http.Get("http://127.0.0.1:8888/update")
	}

	resp(c, 200, "2fa_gateway.html", &baseTemplateData{
		TitleBar:  "Two Factor Authentication",
		KyutGrill: "2fa.jpg",
	})
}

func clientIP(c *gin.Context) string {
	ff := c.Request.Header.Get("CF-Connecting-IP")
	if ff != "" {
		return ff
	}
	return c.ClientIP()
}

func clear2fa(c *gin.Context) {
	// basically deletes from db 2fa tokens, so that it gets regenerated when user hits gateway page
	sess := getSession(c)
	i, _ := sess.Get("userid").(int)
	if i == 0 {
		c.Redirect(302, "/")
	}
	db.Exec("DELETE FROM 2fa WHERE userid = ? AND ip = ?", i, clientIP(c))
	addMessage(c, successMessage{"A new code has been generated and sent to you through Telegram."})
	sess.Save()
	c.Redirect(302, "/2fa_gateway")
}

func verify2fa(c *gin.Context) {
	sess := getSession(c)
	i, _ := sess.Get("userid").(int)
	if i == 0 {
		c.Redirect(302, "/")
	}
	var id int
	var expire common.UnixTimestamp
	err := db.QueryRow("SELECT id, expire FROM 2fa WHERE userid = ? AND ip = ? AND token = ?", i, clientIP(c), strings.ToUpper(c.Query("token"))).Scan(&id, &expire)
	if err == sql.ErrNoRows {
		c.String(200, "1")
		return
	}
	if time.Now().After(time.Time(expire)) {
		c.String(200, "1")
		db.Exec("INSERT INTO 2fa(userid, token, ip, expire, sent) VALUES (?, ?, ?, ?, 0);",
			i, strings.ToUpper(rs.String(8)), clientIP(c), time.Now().Add(time.Hour).Unix())
		http.Get("http://127.0.0.1:8888/update")
		return
	}
	s, err := generateToken(i, c)
	if err != nil {
		resp500(c)
		c.Error(err)
		return
	}
	sess.Set("token", s)
	logIP(c, i)
	addMessage(c, successMessage{"You've been successfully logged in."})
	sess.Delete("2fa_must_validate")
	sess.Save()
	db.Exec("DELETE FROM 2fa WHERE id = ?", id)
	c.String(200, "0")
}
