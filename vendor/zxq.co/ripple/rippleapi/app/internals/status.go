// Package internals has methods that suit none of the API packages.
package internals

import (
	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	Status int `json:"status"`
}

// Status is used for checking the API is up by the ripple website, on the status page.
func Status(c *gin.Context) {
	c.JSON(200, statusResponse{
		Status: 1,
	})
}
