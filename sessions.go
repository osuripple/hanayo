package main

import (
	"net/url"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kawatapw/api/common"
	"zxq.co/x/rs"
)

func sessionInitializer() func(c *gin.Context) {
	return func(c *gin.Context) {
		sess := sessions.Default(c)

		var ctx context

		var passwordChanged bool
		userid := sess.Get("userid")
		if userid, ok := userid.(int); ok {
			ctx.User.ID = userid
			var (
				pRaw     int64
				password string
			)
			err := db.QueryRow("SELECT username, privileges, flags, password_md5 FROM users WHERE id = ?", userid).
				Scan(&ctx.User.Username, &pRaw, &ctx.User.Flags, &password)
			if err != nil {
				c.Error(err)
			}
			if sess.Get("logout") == nil {
				sess.Set("logout", rs.String(15))
			}
			ctx.User.Privileges = common.UserPrivileges(pRaw)
			db.Exec("UPDATE users SET latest_activity = ? WHERE id = ?", time.Now().Unix(), userid)
			if s, ok := sess.Get("pw").(string); !ok || cmd5(password) != s {
				ctx = context{}
				sess.Clear()
				passwordChanged = true
			}
		}

		var addBannedMessage bool
		if ctx.User.ID != 0 && (ctx.User.Privileges&common.UserPrivilegeNormal == 0) {
			ctx = context{}
			sess.Clear()
			addBannedMessage = true
		}

		ctx.Language = getLanguageFromGin(c)

		c.Set("context", ctx)
		c.Set("session", sess)

		if addBannedMessage {
			addMessage(c, warningMessage{T(c, "You have been automatically logged out of your account because your account has either been banned or locked. Should you believe this is a mistake, you can contact our support team at support@ripple.moe.")})
		}
		if passwordChanged {
			addMessage(c, warningMessage{T(c, "You have been automatically logged out for security reasons. Please <a href='/login?redir=%s'>log back in</a>.", url.QueryEscape(c.Request.URL.Path))})
		}

		c.Next()
	}
}

func addMessage(c *gin.Context, m message) {
	sess := getSession(c)
	var messages []message
	messagesRaw := sess.Get("messages")
	if messagesRaw != nil {
		messages = messagesRaw.([]message)
	}
	messages = append(messages, m)
	sess.Set("messages", messages)
}

func getMessages(c *gin.Context) []message {
	sess := getSession(c)
	messagesRaw := sess.Get("messages")
	if messagesRaw == nil {
		return nil
	}
	sess.Delete("messages")
	return messagesRaw.([]message)
}

func getSession(c *gin.Context) sessions.Session {
	return c.MustGet("session").(sessions.Session)
}
