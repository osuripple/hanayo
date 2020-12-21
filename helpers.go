package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kawatapw/hanayo/modules/bbcode"
	tp "github.com/kawatapw/hanayo/modules/top-passwords"
)

//go:generate go run scripts/generate_mappings.go -g
//go:generate go run scripts/top_passwords.go

func cmd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func validatePassword(p string) string {
	if len(p) < 8 {
		return "Your password is too short! It must be at least 8 characters long."
	}

	for _, k := range tp.TopPasswords {
		if k == p {
			return "Your password is one of the most common passwords on the entire internet. No way we're letting you use that!"
		}
	}

	return ""
}

func recaptchaCheck(c *gin.Context) bool {
	f := make(url.Values)
	f.Add("secret", config.RecaptchaPrivate)
	f.Add("response", c.PostForm("g-recaptcha-response"))
	f.Add("remoteip", clientIP(c))

	req, err := http.Post("https://www.google.com/recaptcha/api/siteverify",
		"application/x-www-form-urlencoded", strings.NewReader(f.Encode()))
	if err != nil {
		c.Error(err)
		return false
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		c.Error(err)
		return false
	}

	var e struct {
		Success bool `json:"success"`
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		c.Error(err)
		return false
	}

	return e.Success
}

func parseBBCode(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.String(200, "Error")
		return
	}
	d := bbcode.Compile(string(body))
	c.String(200, d)
}

func mustCSRFGenerate(u int) string {
	v, err := CSRF.Generate(u)
	if err != nil {
		panic(err)
	}
	return v
}
