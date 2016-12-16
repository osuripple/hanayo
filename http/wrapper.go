package http

import (
	"net/http"

	"git.zxq.co/ripple/hanayo/fail"
	"github.com/julienschmidt/httprouter"
)

// HandlerFunc is a function capable of handling requests.
type HandlerFunc func(*Context) error

// Wrapper wraps handler functions so that they can be handled by httprouter.
func Wrapper(f HandlerFunc, s *Server) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		dispatcher(f, w, r, p, s)
	}
}

func dispatcher(f HandlerFunc, w http.ResponseWriter, r *http.Request, p httprouter.Params, s *Server) {
	ctx := &Context{
		Request:     r,
		Writer:      w,
		UserService: s.UserService,
	}

	err := f(ctx)
	switch err.(type) {
	case fail.Error:
		// user should get notified about it
	default:
		// error should get logged
	}
}
