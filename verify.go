package main

import (
	"fmt"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
	"git.zxq.co/x/rs"
	"github.com/gin-gonic/gin"
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

	if !csrfExist(ctx.User.ID, c.Query("csrf")) {
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
