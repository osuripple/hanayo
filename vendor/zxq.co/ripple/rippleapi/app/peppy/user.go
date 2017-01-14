// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"database/sql"
	"fmt"

	"zxq.co/ripple/ocl"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thehowl/go-osuapi"
)

// GetUser retrieves general user information.
func GetUser(c *gin.Context, db *sqlx.DB) {
	if c.Query("u") == "" {
		c.JSON(200, defaultResponse)
		return
	}
	var user osuapi.User
	whereClause, p := genUser(c, db)
	whereClause = "WHERE " + whereClause

	mode := genmode(c.Query("m"))

	var lbpos *int
	err := db.QueryRow(fmt.Sprintf(
		`SELECT
			users.id, users.username,
			users_stats.playcount_%s, users_stats.ranked_score_%s, users_stats.total_score_%s,
			leaderboard_%s.position, users_stats.pp_%s, users_stats.avg_accuracy_%s,
			users_stats.country
		FROM users
		LEFT JOIN users_stats ON users_stats.id = users.id
		INNER JOIN leaderboard_%s ON leaderboard_%s.user = users.id
		%s
		LIMIT 1`,
		mode, mode, mode, mode, mode, mode, mode, mode, whereClause,
	), p).Scan(
		&user.UserID, &user.Username,
		&user.Playcount, &user.RankedScore, &user.TotalScore,
		&lbpos, &user.PP, &user.Accuracy,
		&user.Country,
	)
	if err != nil {
		c.JSON(200, defaultResponse)
		if err != sql.ErrNoRows {
			c.Error(err)
		}
		return
	}
	if lbpos != nil {
		user.Rank = *lbpos
	}
	user.Level = ocl.GetLevelPrecise(user.TotalScore)

	c.JSON(200, []osuapi.User{user})
}
