package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PeppyMethod generates a method for the peppyapi
func PeppyMethod(a func(c *gin.Context, db *sqlx.DB)) gin.HandlerFunc {
	return func(c *gin.Context) {
		rateLimiter()
		perUserRequestLimiter(0, c.ClientIP())

		doggo.Incr("requests.peppy", nil, 1)

		// I have no idea how, but I manged to accidentally string the first 4
		// letters of the alphabet into a single function call.
		a(c, db)
	}
}
