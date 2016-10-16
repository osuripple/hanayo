package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func loginSubmit(c *gin.Context) {
	if c.MustGet("context").(context).User.ID != 0 {
		loginSubmitReplyError(c, "You are already logged in!")
		return
	}

	if c.PostForm("username") == "" || c.PostForm("password") == "" {
		loginSubmitReplyError(c, "Username or password not set.")
		return
	}

	param := "username"
	if strings.Contains(c.PostForm("username"), "@") {
		param = "email"
	}

	var data struct {
		ID              int
		Username        string
		Password        string
		PasswordVersion int
		Country         string
		pRaw            int64
		Privileges      common.UserPrivileges
	}
	err := db.QueryRow(`
	SELECT 
		u.id, u.password_md5,
		u.username, u.password_version,
		s.country, u.privileges
	FROM users u
	LEFT JOIN users_stats s ON s.id = u.id
	WHERE u.`+param+` = ? LIMIT 1`, strings.TrimSpace(c.PostForm("username"))).Scan(
		&data.ID, &data.Password,
		&data.Username, &data.PasswordVersion,
		&data.Country, &data.pRaw,
	)
	data.Privileges = common.UserPrivileges(data.pRaw)

	switch {
	case err == sql.ErrNoRows:
		loginSubmitReplyError(c, "No user with such "+param+"!")
		return
	case err != nil:
		c.Error(err)
		resp500(c)
		return
	}

	if data.PasswordVersion == 1 {
		addMessage(c, warningMessage{"Your password is sooooooo old, that we don't even know how to deal with it anymore. Could you please change it?"})
		c.Redirect(302, "/forgot_password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(fmt.Sprintf("%x", md5.Sum([]byte(c.PostForm("password")))))); err != nil {
		loginSubmitReplyError(c, "Wrong password.")
		return
	}

	if data.Privileges&common.UserPrivilegeNormal == 0 {
		loginSubmitReplyError(c, "You are not allowed to login. This means your account is either banned or locked.")
		return
	}

	if data.Country == "XX" {
		// TODO
	}

	setYCookie(data.ID, c)

	sess := c.MustGet("session").(sessions.Session)
	sess.Set("userid", data.ID)

	tfaEnabled := is2faEnabled(data.ID)
	if !tfaEnabled {
		s, err := generateToken(data.ID, c)
		if err != nil {
			resp500(c)
			c.Error(err)
			return
		}
		sess.Set("token", s)
		logIP(c, data.ID)
	} else {
		sess.Set("2fa_must_validate", true)
	}

	sess.Save()

	if tfaEnabled {
		c.Redirect(302, "/2fa_gateway/generate")
	} else {
		addMessage(c, successMessage{fmt.Sprintf("Hey %s! You are now logged in.", c.PostForm("username"))})
		c.Redirect(302, "/")
	}
	return
}

func loginSubmitReplyError(c *gin.Context, msg string) {
	resp(c, 200, "login.html", &baseTemplateData{
		TitleBar:  "Log in",
		KyutGrill: "login.png",
		Messages: []message{
			errorMessage{msg},
		},
		FormData: map[string]string{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		},
	})
}

func logout(c *gin.Context) {
	ctx := c.MustGet("context").(context)
	if ctx.User.ID == 0 {
		resp(c, 200, "empty.html", &baseTemplateData{TitleBar: "Log out", Messages: []message{warningMessage{"You're already logged out!"}}})
		return
	}
	sess := c.MustGet("session").(sessions.Session)
	sess.Clear()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "rt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
	addMessage(c, successMessage{"Successfully logged out."})
	sess.Save()
	c.Redirect(302, "/")
}

func generateToken(id int, c *gin.Context) (string, error) {
	tok := common.RandomString(32)
	_, err := db.Exec(
		`INSERT INTO tokens(user, privileges, description, token, private)
					VALUES (   ?,        '0',           ?,     ?,     '1');`,
		id, c.Request.Header.Get("X-Real-IP"), fmt.Sprintf("%x", md5.Sum([]byte(tok))))
	if err != nil {
		return "", err
	}
	return tok, nil
}
