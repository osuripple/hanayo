package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type discordLinkResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type discordLinkRequestMethod int

const (
	getRequest  discordLinkRequestMethod = 0
	postRequest                          = 1
)

func discordLinkURLSimple(handler string) string {
	return fmt.Sprintf(
		config.OldFrontend+"/discord/%s?k=%s",
		handler, config.DonorBotSecret,
	)
}

func discordLinkURL(handler string, ept string, qs ...interface{}) string {
	return fmt.Sprintf(
		"%s&%s",
		discordLinkURLSimple(handler),
		fmt.Sprintf(ept, qs...),
	)
}

func discordLinkRequest(
	url string, m discordLinkRequestMethod, postBody *map[string]interface{},
) (*discordLinkResponse, error) {
	var resp *http.Response
	var err error

	if m == getRequest {
		resp, err = http.Get(url)
	} else {
		var jsonB []byte
		jsonB, err = json.Marshal(postBody)
		if err != nil {
			return nil, err
		}
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonB))
	}
	if err != nil {
		return nil, err
	}
	var o discordLinkResponse
	err = json.NewDecoder(resp.Body).Decode(&o)
	if err != nil {
		return nil, fmt.Errorf("o9d decode: %v", err)
	}
	if o.Status != 200 {
		return nil, fmt.Errorf("o9d: %d (%s)", resp.StatusCode, o.Message)
	}
	return &o, nil
}

func discordLinkGet(url string) (*discordLinkResponse, error) {
	return discordLinkRequest(url, getRequest, nil)
}

func discordLinkPost(
	url string, postBody *map[string]interface{},
) (*discordLinkResponse, error) {
	return discordLinkRequest(url, postRequest, postBody)
}

func discordLinkFinish(c *gin.Context) {
	defer func() {
		getSession(c).Save()
		c.Redirect(302, "/settings/discord")
	}()
	_, err := discordLinkGet(
		discordLinkURL(
			"oauth.php",
			"state=%s&code=%s",
			c.Query("state"), c.Query("code"),
		),
	)
	if err != nil {
		c.Error(err)
		addMessage(c, errorMessage{T(c, "An error occurred.")})
		return
	}
	addMessage(c, successMessage{T(c, "Your account has been linked successfully!")})
}

func discordUnlink(c *gin.Context) {
	defer func() {
		getSession(c).Save()
		c.Redirect(302, "/settings/discord")
	}()
	ctx := getContext(c)
	if ok, _ := CSRF.Validate(ctx.User.ID, c.Query("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}
	_, err := discordLinkGet(
		discordLinkURL("unlink.php", "uid=%d", ctx.User.ID),
	)
	if err != nil {
		c.Error(err)
		addMessage(c, errorMessage{T(c, "An error occurred.")})
		return
	}
	addMessage(c, successMessage{T(c, "Your account has been unlinked successfully!")})
}

func discordSubmit(c *gin.Context) {
	defer func() {
		getSession(c).Save()
		c.Redirect(302, "/settings/discord")
	}()
	ctx := getContext(c)
	if ok, _ := CSRF.Validate(ctx.User.ID, c.PostForm("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}
	_, err := discordLinkPost(
		discordLinkURLSimple("role.php"),
		&map[string]interface{}{
			"uid":    ctx.User.ID,
			"colour": c.PostForm("colour"),
			"name":   c.PostForm("name"),
		},
	)
	if err != nil {
		c.Error(err)
		addMessage(c, errorMessage{T(c, "An error occurred.")})
		return
	}
	addMessage(c, successMessage{T(c, "Your custom role has been edited successfully!")})
}
