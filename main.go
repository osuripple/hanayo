package main

import (
	"database/sql"
	"fmt"

	"git.zxq.co/ripple/schiavolib"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thehowl/conf"
)

var (
	c  cf
	db *sql.DB
)

func main() {
	fmt.Println("hanayo v0.0.1")

	err := conf.Load(&c, "hanayo.conf")
	switch err {
	case nil:
		// carry on
	case conf.ErrNoFile:
		conf.Export(c, "hanayo.conf")
		fmt.Println("The configuration file was not found. We created one for you.")
		return
	default:
		panic(err)
	}

	db, err = sql.Open("mysql", c.DSN)
	if err != nil {
		panic(err)
	}

	if gin.Mode() == gin.DebugMode {
		fmt.Println("Development environment detected. Starting fsnotify on template folder...")
		err := reloader()
		if err != nil {
			fmt.Println(err)
		}
	}

	schiavo.Bunker.Send(fmt.Sprintf("**hanayo** STARTUATO, mode: %s", gin.Mode()))

	fmt.Println("Importing templates...")
	loadTemplates()

	fmt.Println("Starting webserver...")

	r := gin.Default()

	r.Static("/static", "static")

	r.GET("/", homePage)

	r.Run(":45221")
}
