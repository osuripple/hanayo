package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type discordLinkResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func discordLinkURL(handler string, ept string, qs ...interface{}) string {
	return fmt.Sprintf(
		config.OldFrontend+"/discord/%s?k=%s&%s",
		handler, config.DonorBotSecret, fmt.Sprintf(ept, qs...),
	)
}

func discordLinkRequest(url string) (*discordLinkResponse, error) {
	resp, err := http.Get(url)
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

func discordFinish(c *gin.Context) {
	defer c.Redirect(302, "/settings/discord")
	_, err := discordLinkRequest(
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
	defer c.Redirect(302, "/settings/discord")
	ctx := getContext(c)
	if ok, _ := CSRF.Validate(ctx.User.ID, c.Query("csrf")); !ok {
		addMessage(c, errorMessage{T(c, "Your session has expired. Please try redoing what you were trying to do.")})
		return
	}
	_, err := discordLinkRequest(
		discordLinkURL("unlink.php", "uid=%d", ctx.User.ID),
	)
	if err != nil {
		c.Error(err)
		addMessage(c, errorMessage{T(c, "An error occurred.")})
		return
	}
	addMessage(c, successMessage{T(c, "Your account has been unlinked successfully!")})
}
