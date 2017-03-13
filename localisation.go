package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// T translates a string into the language specified by the request.
func T(c *gin.Context, s string, args ...interface{}) string {
	return fmt.Sprintf(s, args...)
}

func (b *baseTemplateData) T(s string, args ...interface{}) string {
	return T(b.Gin, s, args...)
}
