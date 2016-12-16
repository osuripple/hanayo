package http

import (
	"github.com/julienschmidt/httprouter"
)

// SetUpRoutes sets up the routes for the router.
func (s *Server) SetUpRoutes() error {
	s.Router = httprouter.New()
	s.GET("/", func(ctx *Context) error {
		ctx.Write([]byte(ctx.Request.URL.Query().Get("Test")))
		return nil
	})
	return nil
}
