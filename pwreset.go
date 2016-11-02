package main

import (
	"database/sql"
	"fmt"
	"strings"

	"git.zxq.co/ripple/rippleapi/common"
	"git.zxq.co/x/rs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mailgun/mailgun-go.v1"
)

func passwordReset(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID != 0 {
		simpleReply(c, errorMessage{"You're already logged in!"})
		return
	}

	field := "username"
	if strings.Contains(c.PostForm("username"), "@") {
		field = "email"
	}

	var (
		id         int
		username   string
		email      string
		privileges uint64
	)

	err := db.QueryRow("SELECT id, username, email, privileges FROM users WHERE "+field+" = ?",
		c.PostForm("username")).
		Scan(&id, &username, &email, &privileges)

	switch err {
	case nil:
		// ignore
	case sql.ErrNoRows:
		simpleReply(c, errorMessage{"That user could not be found."})
		return
	default:
		c.Error(err)
		resp500(c)
		return
	}

	if common.UserPrivileges(privileges)&
		(common.UserPrivilegeNormal|common.UserPrivilegePendingVerification) == 0 {
		simpleReply(c, errorMessage{"You look pretty banned/locked here."})
		return
	}

	// generate key
	key := rs.String(50)

	// TODO: WHY THE FUCK DOES THIS USE USERNAME AND NOT ID PLEASE WRITE MIGRATION
	_, err = db.Exec("INSERT INTO password_recovery(k, u) VALUES (?, ?)", key, username)

	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	content := fmt.Sprintf(
		"Hey %s! Someone, which we really hope was you, requested a password reset for your account. "+
			"In case it was you, please <a href='%s'>click here</a> to reset your password on Ripple. "+
			"Otherwise, silently ignore this email.",
		username,
		config.BaseURL+"/pwreset/continue?k="+key,
	)
	msg := mailgun.NewMessage(
		config.MailgunFrom,
		"Ripple password recovery instructions",
		content,
		email,
	)
	msg.SetHtml(content)
	_, _, err = mg.Send(msg)

	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	addMessage(c, successMessage{"Done! You should shortly receive an email from us at the email you used to sign up on Ripple."})
	getSession(c).Save()
	c.Redirect(302, "/")
}

type passwordResetContinueTData struct {
	baseTemplateData
	Username string
	Key      string
}

func passwordResetContinue(c *gin.Context) {
	k := c.Query("k")

	if k == "" {
		respEmpty(c, "Password reset", errorMessage{"Nope."})
		return
	}

	var username string
	switch err := db.QueryRow("SELECT u FROM password_recovery WHERE k = ? LIMIT 1", k).
		Scan(&username); err {
	case nil:
		// move on
	case sql.ErrNoRows:
		respEmpty(c, "Reset password", errorMessage{"That key could not be found. Perhaps it expired?"})
		return
	default:
		c.Error(err)
		resp500(c)
		return
	}

	renderResetPassword(c, username, k)
}

func passwordResetContinueSubmit(c *gin.Context) {
	var username string
	switch err := db.QueryRow("SELECT u FROM password_recovery WHERE k = ? LIMIT 1", c.PostForm("k")).
		Scan(&username); err {
	case nil:
		// move on
	case sql.ErrNoRows:
		respEmpty(c, "Reset password", errorMessage{"That key could not be found. Perhaps it expired?"})
		return
	default:
		c.Error(err)
		resp500(c)
		return
	}

	p := c.PostForm("password")

	if s := validatePassword(p); s != "" {
		renderResetPassword(c, username, c.PostForm("k"), errorMessage{s})
		return
	}

	pass, err := generatePassword(p)
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	_, err = db.Exec("UPDATE users SET password_md5 = ?, salt = '', password_version = '2' WHERE username = ?",
		pass, username)
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	_, err = db.Exec("DELETE FROM password_recovery WHERE k = ? LIMIT 1", c.PostForm("k"))
	if err != nil {
		c.Error(err)
		resp500(c)
		return
	}

	addMessage(c, successMessage{"All right, we have changed your password and you should now be able to login! Have fun!"})
	getSession(c).Save()
	c.Redirect(302, "/login")
}

func renderResetPassword(c *gin.Context, username, k string, messages ...message) {
	resp(c, 200, "pwreset/continue.html", &passwordResetContinueTData{
		Username: username,
		Key:      k,
		baseTemplateData: baseTemplateData{
			TitleBar: "Reset password",
			Messages: messages,
		},
	})
}

func generatePassword(p string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(cmd5(p)), bcrypt.DefaultCost)
	return string(s), err
}
