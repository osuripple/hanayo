// Package http handles HTTP requests, middlewares and check on Hanayo.
package http

import (
	"fmt"
	"net/http"
	"time"

	"git.zxq.co/ripple/hanayo"
	"github.com/julienschmidt/httprouter"
)

// Server is an HTTP server on Hanayo.
type Server struct {
	UserService          hanayo.UserService
	TFAService           hanayo.TFAService
	SystemSettingService hanayo.SystemSettingService
	Router               *httprouter.Router
}

// ServeHTTP hands over the request to the router, after logging the request
// and setting up a panic handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := time.Now()

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Panic:", err)
		}

		var ip string
		switch {
		case r.Header.Get("CF-Connecting-IP") != "":
			ip = r.Header.Get("CF-Connecting-IP")
		case r.Header.Get("X-Forwarded-For") != "":
			ip = r.Header.Get("X-Forwarded-For")
		default:
			ip = r.RemoteAddr
		}

		fmt.Printf(
			"%-39s | %-13s | %7s %-13s\n",
			ip, time.Since(b), r.Method, r.URL.Path,
		)
	}()

	s.Router.ServeHTTP(w, r)
}

// GET registers a handler for a GET request.
func (s *Server) GET(p string, h HandlerFunc) { s.Router.Handle("GET", p, Wrapper(h, s)) }

// POST registers a handler for a POST request.
func (s *Server) POST(p string, h HandlerFunc) { s.Router.Handle("POST", p, Wrapper(h, s)) }

// DELETE registers a handler for a DELETE request.
func (s *Server) DELETE(p string, h HandlerFunc) { s.Router.Handle("DELETE", p, Wrapper(h, s)) }

// OPTIONS registers a handler for a OPTIONS request.
func (s *Server) OPTIONS(p string, h HandlerFunc) { s.Router.Handle("OPTIONS", p, Wrapper(h, s)) }

// HEAD registers a handler for a HEAD request.
func (s *Server) HEAD(p string, h HandlerFunc) { s.Router.Handle("HEAD", p, Wrapper(h, s)) }

// PATCH registers a handler for a PATCH request.
func (s *Server) PATCH(p string, h HandlerFunc) { s.Router.Handle("PATCH", p, Wrapper(h, s)) }

// PUT registers a handler for a PUT request.
func (s *Server) PUT(p string, h HandlerFunc) { s.Router.Handle("PUT", p, Wrapper(h, s)) }
