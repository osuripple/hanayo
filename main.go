package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"

	"git.zxq.co/ripple/schiavolib"
	"git.zxq.co/x/rs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thehowl/conf"
)

var (
	c struct {
		DSN string

		CookieSecret string

		RedisEnable         bool
		RedisMaxConnections int
		RedisNetwork        string
		RedisAddress        string
		RedisPassword       string

		AvatarURL string
		BaseURL   string
	}
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

	if c.CookieSecret == "" {
		c.CookieSecret = rs.String(46)
	}
	if c.AvatarURL == "" {
		c.AvatarURL = "https://a.ripple.moe"
	}
	if c.BaseURL == "" {
		c.BaseURL = "https://ripple.moe"
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

	schiavo.Prefix = "hanayo"
	schiavo.Bunker.Send(fmt.Sprintf("STARTUATO, mode: %s", gin.Mode()))

	fmt.Println("Starting session system...")
	var store sessions.Store
	if c.RedisMaxConnections != 0 {
		store, err = sessions.NewRedisStore(
			c.RedisMaxConnections,
			c.RedisNetwork,
			c.RedisAddress,
			c.RedisPassword,
			[]byte(c.CookieSecret),
		)
		if err != nil {
			fmt.Println(err)
			store = sessions.NewCookieStore([]byte(c.CookieSecret))
		}
	} else {
		store = sessions.NewCookieStore([]byte(c.CookieSecret))
	}
	gobRegisters := []interface{}{
		[]message{},
		errorMessage{},
		infoMessage{},
		neutralMessage{},
		warningMessage{},
		successMessage{},
	}
	for _, el := range gobRegisters {
		gob.Register(el)
	}

	fmt.Println("Importing templates...")
	loadTemplates()

	fmt.Println("Starting webserver...")

	r := gin.Default()

	r.Use(
		sessions.Sessions("session", store),
		sessionInitializer(),
	)

	r.Static("/static", "static")

	r.GET("/", homePage)

	r.GET("/login", login)
	r.POST("/login", loginSubmit)
	r.GET("/logout", logout)

	r.NoRoute(notFound)

	conf.Export(c, "hanayo.conf")

	r.Run(":45221")
}
