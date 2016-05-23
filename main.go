package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hanayo v0.0.1")

	fmt.Println("Importing templates...")
	loadTemplates()

	fmt.Println("Starting webserver...")

	r := gin.Default()

	r.Static("/static", "static")

	r.GET("/", testHandler)

	r.Run(":45221")
}
