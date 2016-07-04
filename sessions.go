package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func sessionInitializer() func(c *gin.Context) {
	return func(c *gin.Context) {
		sess := sessions.Default(c)

		var ctx context
		tok := sess.Get("token")
		if tok, ok := tok.(string); ok {
			ctx.Token = tok
		}
		userid := sess.Get("userid")
		if userid, ok := userid.(int); ok {
			ctx.User.ID = userid
			db.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&ctx.User.Username)
		}

		// TODO: Add stay logged in
		// TODO: log out if banned

		c.Set("context", ctx)
		c.Set("session", sess)

		c.Next()
	}
}

func addMessage(c *gin.Context, m message) {
	sess := c.MustGet("session").(sessions.Session)
	var messages []message
	messagesRaw := sess.Get("messages")
	if messagesRaw != nil {
		messages = messagesRaw.([]message)
	}
	messages = append(messages, m)
	sess.Set("messages", messages)
	sess.Save()
}

func getMessages(c *gin.Context) []message {
	sess := c.MustGet("session").(sessions.Session)
	messagesRaw := sess.Get("messages")
	if messagesRaw == nil {
		return nil
	}
	sess.Delete("messages")
	sess.Save()
	return messagesRaw.([]message)
}
