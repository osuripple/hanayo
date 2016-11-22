package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"net/url"

	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func loginSubmit(c *gin.Context) {
	if getContext(c).User.ID != 0 {
		simpleReply(c, errorMessage{"You are already logged in!"})
		return
	}

	if c.PostForm("username") == "" || c.PostForm("password") == "" {
		simpleReply(c, errorMessage{"Username or password not set."})
		return
	}

	param := "username_safe"
	u := c.PostForm("username")
	if strings.Contains(u, "@") {
		param = "email"
	} else {
		u = safeUsername(u)
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
	WHERE u.`+param+` = ? LIMIT 1`, strings.TrimSpace(u)).Scan(
		&data.ID, &data.Password,
		&data.Username, &data.PasswordVersion,
		&data.Country, &data.pRaw,
	)
	data.Privileges = common.UserPrivileges(data.pRaw)

	switch {
	case err == sql.ErrNoRows:
		if param == "username_safe" {
			param = "username"
		}
		simpleReply(c, errorMessage{"No user with such " + param + "!"})
		return
	case err != nil:
		c.Error(err)
		resp500(c)
		return
	}

	if data.PasswordVersion == 1 {
		addMessage(c, warningMessage{"Your password is sooooooo old, that we don't even know how to deal with it anymore. Could you please change it?"})
		c.Redirect(302, "/pwreset")
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(data.Password),
		[]byte(cmd5(c.PostForm("password"))),
	); err != nil {
		simpleReply(c, errorMessage{"Wrong password."})
		return
	}

	// TODO: if bcrypt.Cost < bcrypt.DefaultCost, regenerate password with
	// default cost.

	sess := getSession(c)

	if data.Privileges&common.UserPrivilegePendingVerification > 0 {
		setYCookie(data.ID, c)
		addMessage(c, warningMessage{"You will need to verify your account first."})
		sess.Save()
		c.Redirect(302, "/register/verify?u="+strconv.Itoa(data.ID))
		return
	}

	if data.Privileges&common.UserPrivilegeNormal == 0 {
		simpleReply(c, errorMessage{"You are not allowed to login. This means your account is either banned or locked."})
		return
	}

	setYCookie(data.ID, c)

	sess.Set("userid", data.ID)
	sess.Set("pw", cmd5(data.Password))

	tfaEnabled := is2faEnabled(data.ID)
	if !tfaEnabled {
		s, err := generateToken(data.ID, c)
		if err != nil {
			resp500(c)
			c.Error(err)
			return
		}
		sess.Set("token", s)
		if data.Country == "XX" {
			setCountry(c, data.ID)
		}
		logIP(c, data.ID)
	} else {
		sess.Set("2fa_must_validate", true)
	}

	redir := c.PostForm("redir")
	if len(redir) > 0 && redir[0] != '/' {
		redir = ""
	}

	if tfaEnabled {
		sess.Save()
		c.Redirect(302, "/2fa_gateway?redir="+url.QueryEscape(redir))
	} else {
		addMessage(c, successMessage{fmt.Sprintf("Hey %s! You are now logged in.", template.HTMLEscapeString(data.Username))})
		sess.Save()
		if redir == "" {
			redir = "/"
		}
		c.Redirect(302, redir)
	}
	return
}

func safeUsername(u string) string {
	return strings.TrimSpace(strings.Replace(strings.ToLower(u), " ", "_", -1))
}

func logout(c *gin.Context) {
	ctx := getContext(c)
	if ctx.User.ID == 0 {
		respEmpty(c, "Log out", warningMessage{"You're already logged out!"})
		return
	}
	sess := getSession(c)
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
		id, clientIP(c), cmd5(tok))
	if err != nil {
		return "", err
	}
	return tok, nil
}

func checkToken(s string, id int, c *gin.Context) (string, error) {
	if err := db.QueryRow("SELECT 1 FROM tokens WHERE token = ?", cmd5(s)).Scan(new(int)); err == sql.ErrNoRows {
		return generateToken(id, c)
	} else if err != nil {
		return "", err
	}
	return s, nil
}
