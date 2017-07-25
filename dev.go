package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"zxq.co/ripple/rippleapi/common"
)

func createAPIToken(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	sess := getSession(c)
	defer func() {
		sess.Save()
		c.Redirect(302, "/dev/tokens")
	}()

	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}

	privileges := common.Privileges(common.Int(c.PostForm("privileges"))).CanOnly(ctx.User.Privileges)
	description := c.PostForm("description")

	var (
		tokenStr string
		tokenMD5 string
	)

	for {
		tokenStr = common.RandomString(32)
		tokenMD5 = fmt.Sprintf("%x", md5.Sum([]byte(tokenStr)))

		var id int
		err := db.QueryRow("SELECT id FROM tokens WHERE token = ? LIMIT 1", tokenMD5).Scan(&id)
		if err == sql.ErrNoRows {
			break
		}
		if err != nil {
			c.Error(err)
			resp500(c)
			return
		}
	}

	_, err := db.Exec("INSERT INTO tokens(user, privileges, description, token, private, last_updated) VALUES (?, ?, ?, ?, '0', ?)",
		ctx.User.ID, privileges, description, tokenMD5, time.Now().Unix())
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	addMessage(c, successMessage{
		fmt.Sprintf("Your token has been created successfully! Your token is: <code>%s</code>.<br>Keep it safe, don't show it around, and store it now! We won't show it to you again.", tokenStr),
	})
}

func deleteAPIToken(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	sess := getSession(c)
	defer func() {
		sess.Save()
		c.Redirect(302, "/dev/tokens")
	}()

	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}

	db.Exec("DELETE FROM tokens WHERE id = ? AND user = ? AND private = 0", c.PostForm("id"), ctx.User.ID)
	addMessage(c, successMessage{"That token has been deleted successfully."})
}

func editAPIToken(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	sess := getSession(c)
	defer func() {
		sess.Save()
		c.Redirect(302, "/dev/tokens/edit?id="+c.PostForm("id"))
	}()

	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}

	privileges := common.Privileges(common.Int(c.PostForm("privileges"))).CanOnly(ctx.User.Privileges)
	description := c.PostForm("description")

	_, err := db.Exec("UPDATE tokens SET privileges = ?, description = ? WHERE user = ? AND id = ? AND private = 0",
		privileges, description, ctx.User.ID, c.PostForm("id"))
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	addMessage(c, successMessage{
		"Your token has been edited successfully!",
	})
}
