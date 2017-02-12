// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/thehowl/go-osuapi"
	"github.com/valyala/fasthttp"
	"zxq.co/ripple/ocl"
	"zxq.co/ripple/rippleapi/common"
)

// GetUser retrieves general user information.
func GetUser(c *fasthttp.RequestCtx, db *sqlx.DB) {
	if query(c, "u") == "" {
		json(c, 200, defaultResponse)
		return
	}
	var user osuapi.User
	whereClause, p := genUser(c, db)
	whereClause = "WHERE " + whereClause

	mode := genmode(query(c, "m"))

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
		json(c, 200, defaultResponse)
		if err != sql.ErrNoRows {
			common.Err(c, err)
		}
		return
	}
	if lbpos != nil {
		user.Rank = *lbpos
	}
	user.Level = ocl.GetLevelPrecise(user.TotalScore)

	json(c, 200, []osuapi.User{user})
}
