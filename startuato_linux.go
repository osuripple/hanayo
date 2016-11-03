// +build !windows

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"git.zxq.co/ripple/schiavolib"
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/goagain"
)

func startuato(engine *gin.Engine) {
	engine.GET("/51/update", updateFromRemote)

	// Inherit a net.Listener from our parent process or listen anew.
	l, err := goagain.Listener()
	if err != nil {

		// Listen on a TCP or a UNIX domain socket (TCP here).
		if config.Unix {
			l, err = net.Listen("unix", config.ListenTo)
		} else {
			l, err = net.Listen("tcp", config.ListenTo)
		}
		if err != nil {
			schiavo.Bunker.Send(err.Error())
			log.Fatalln(err)
		}

		schiavo.Bunker.Send(fmt.Sprint("LISTENINGU STARTUATO ON ", l.Addr()))

		// Accept connections in a new goroutine.
		go http.Serve(l, engine)

	} else {

		// Resume accepting connections in a new goroutine.
		schiavo.Bunker.Send(fmt.Sprint("LISTENINGU RESUMINGU ON ", l.Addr()))
		go http.Serve(l, engine)

		// Kill the parent, now that the child has started successfully.
		if err := goagain.Kill(); err != nil {
			schiavo.Bunker.Send(err.Error())
			log.Fatalln(err)
		}

	}

	// Block the main goroutine awaiting signals.
	if _, err := goagain.Wait(l); err != nil {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}

	// Do whatever's necessary to ensure a graceful exit like waiting for
	// goroutines to terminate or a channel to become closed.
	//
	// In this case, we'll simply stop listening and wait one second.
	if err := l.Close(); err != nil {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}
	if err := db.Close(); err != nil {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}
	time.Sleep(time.Second * 1)
}

func updateFromRemote(c *gin.Context) {
	if c.Query("hanayokey") != config.APISecret {
		c.String(403, "nope")
		return
	}
	if f, err := os.Stat(".git"); err == os.ErrNotExist || !f.IsDir() {
		c.String(500, "not git ffs")
		return
	}
	c.String(200, "all right")
	go func() {
		if !execCommand("git", "pull", "origin", "master") {
			return
		}
		// go get
		//        -u: update all dependencies
		//        -d: stop after downloading deps
		if !execCommand("go", "get", "-v", "-u", "-d") {
			return
		}
		if !execCommand("bash", "-c", "go build -v") {
			return
		}

		proc, err := os.FindProcess(syscall.Getpid())
		if err != nil {
			log.Println(err)
			return
		}
		proc.Signal(syscall.SIGUSR2)
	}()
}

func execCommand(command string, args ...string) bool {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return false
	}
	data, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println(err)
		return false
	}
	// Bob. We got a problem.
	if len(data) != 0 {
		log.Println(string(data))
	}
	io.Copy(os.Stdout, stdout)
	cmd.Wait()
	stdout.Close()
	return true
}
