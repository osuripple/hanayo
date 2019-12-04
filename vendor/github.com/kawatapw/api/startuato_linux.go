// +build !windows

package main

import (
	"log"
	"net"
	"time"

	"github.com/rcrowley/goagain"
	"github.com/valyala/fasthttp"
	"github.com/kawatapw/api/common"
)

func startuato(hn fasthttp.RequestHandler) {
	conf, _ := common.Load()
	// Inherit a net.Listener from our parent process or listen anew.
	l, err := goagain.Listener()
	if nil != err {

		// Listen on a TCP or a UNIX domain socket (TCP here).
		if conf.Unix {
			l, err = net.Listen("unix", conf.ListenTo)
		} else {
			l, err = net.Listen("tcp", conf.ListenTo)
		}
		if nil != err {
			log.Fatalln(err)
		}

		// Accept connections in a new goroutine.
		go fasthttp.Serve(l, hn)
	} else {

		// Resume accepting connections in a new goroutine.
		go fasthttp.Serve(l, hn)

		// Kill the parent, now that the child has started successfully.
		if err := goagain.Kill(); nil != err {
			log.Fatalln(err)
		}

	}

	// Block the main goroutine awaiting signals.
	if _, err := goagain.Wait(l); nil != err {
		log.Fatalln(err)
	}

	// Do whatever's necessary to ensure a graceful exit like waiting for
	// goroutines to terminate or a channel to become closed.
	//
	// In this case, we'll simply stop listening and wait one second.
	if err := l.Close(); nil != err {
		log.Fatalln(err)
	}
	if err := db.Close(); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second * 1)
}
