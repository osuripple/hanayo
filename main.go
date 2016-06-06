package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hanayo v0.0.1")

	if gin.Mode() == gin.DebugMode {
		fmt.Println("Development environment detected. Starting fsnotify on template folder...")
		err := reloader()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Importing templates...")
	loadTemplates()

	fmt.Println("Starting webserver...")

	r := gin.Default()

	r.Static("/static", "static")

	r.GET("/", homePage)

	r.Run(":45221")
}
