package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"zxq.co/ripple/rippleapi/common"
	"zxq.co/x/rs"
)

func startEmailVerification(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	var m message
	defer func() {
		if m != nil {
			respEmpty(c, "Verify email", m)
		}
	}()

	if ok, _ := CSRF.Validate(ctx.User.ID, c.Query("csrf")); !ok {
		m = errorMessage{"CSRF token expired. Please try again."}
		return
	}

	if ctx.User.Flags&common.FlagEmailVerified > 0 {
		m = errorMessage{"Your email has already been verified!"}
		return
	}

	var key string
	for i := 0; i < 10; i++ {
		key = rs.String(50)
		_, err := db.Exec(
			"INSERT INTO verification_emails(`key`, user, `time`) VALUES (?, ?, ?)",
			key, ctx.User.ID, time.Now().Unix(),
		)
		if err == nil {
			break
		}
		fmt.Println("verification email:", err)
	}

	var email string
	db.Get(&email, "SELECT email FROM users WHERE id = ?", ctx.User.ID)

	content := fmt.Sprintf(`Howdy, %s! Someone, which we hope was you, requested to have their email verified.
In case it was you indeed, <a href="%s">click on this link.</a>`,
		ctx.User.Username, config.BaseURL+"/email_verify/finish?k="+key)
	msg := mg.NewMessage(config.MailgunFrom, "Ripple email verification", content, email)
	msg.SetHtml(content)

	_, _, err := mg.Send(msg)
	if err != nil {
		m = errorMessage{"An error occurred."}
		c.Error(err)
		return
	}

	addMessage(c, successMessage{"Success! You should have received an email with instructions on how to verify your email address."})
	getSession(c).Save()
	c.Redirect(302, "/")
}

func finishEmailVerification(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		resp403(c)
		return
	}

	if ctx.User.Flags&common.FlagEmailVerified > 0 {
		respEmpty(c, "Verify email", errorMessage{"Your account has already been verified once!"})
		return
	}

	k := c.Query("k")
	var u int
	db.Get(
		&u,
		"SELECT user FROM verification_emails WHERE `key` = ? AND user = ? AND `time` > ?",
		k, ctx.User.ID, time.Now().Add(-time.Hour*24).Unix(),
	)

	if u != ctx.User.ID {
		respEmpty(c, "Verify email", errorMessage{"The email verification you were looking for could not be found. Perhaps it's another user's, or it expired?"})
		return
	}

	db.Exec("DELETE FROM verification_emails WHERE user = ?", ctx.User.ID)
	db.Exec("UPDATE users SET flags = flags | 3 WHERE id = ?", ctx.User.ID)

	addMessage(c, successMessage{"Your email has been verified. Thanks!"})
	getSession(c).Save()
	c.Redirect(302, "/")
}
