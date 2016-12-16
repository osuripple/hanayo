package http

import (
	"net/http"

	"git.zxq.co/ripple/hanayo"
)

// Context is the information about the request that is passed to all handlers.
type Context struct {
	http.Request
	Writer http.ResponseWriter

	User   *hanayo.User
	Token  string
	Errors []error

	// Services
	UserService hanayo.UserService
}

// Err adds an error to errors, that will be then reported to sentry
// (if enabled)
func (c *Context) Err(err error) {
	c.Errors = append(c.Errors, err)
}
