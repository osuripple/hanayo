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
	UserService hanayo.UserService
	Router      *httprouter.Router
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
