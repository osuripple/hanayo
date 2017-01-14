// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetMatch retrieves general match information.
func GetMatch(c *gin.Context, db *sqlx.DB) {
	c.JSON(200, defaultResponse)
}
