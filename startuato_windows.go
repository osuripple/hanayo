// +build windows

package main

import (
	"net"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func startuato(engine *gin.Engine) {
	var (
		l net.Listener
		err error
	)
	// Listen on a TCP or a UNIX domain socket (TCP here).
	if config.Unix {
		l, err = net.Listen("unix", config.ListenTo)
	} else {
		l, err = net.Listen("tcp", config.ListenTo)
	}
	if nil != err {
		log.Fatalln(err)
	}

	http.Serve(l, engine)
}
