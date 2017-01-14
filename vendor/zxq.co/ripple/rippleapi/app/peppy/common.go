package peppy

import (
	"database/sql"
	"strconv"

	"zxq.co/ripple/rippleapi/common"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var modes = []string{"std", "taiko", "ctb", "mania"}

var defaultResponse = []struct{}{}

func genmode(m string) string {
	i := genmodei(m)
	return modes[i]
}
func genmodei(m string) int {
	v := common.Int(m)
	if v > 3 || v < 0 {
		v = 0
	}
	return v
}
func rankable(m string) bool {
	x := genmodei(m)
	return x == 0 || x == 3
}

func genUser(c *gin.Context, db *sqlx.DB) (string, string) {
	var whereClause string
	var p string

	// used in second case of switch
	s, err := strconv.Atoi(c.Query("u"))

	switch {
	// We know for sure that it's an username.
	case c.Query("type") == "string":
		whereClause = "users.username_safe = ?"
		p = common.SafeUsername(c.Query("u"))
	// It could be an user ID, so we look for an user with that username first.
	case err == nil:
		err = db.QueryRow("SELECT id FROM users WHERE id = ? LIMIT 1", s).Scan(&p)
		if err == sql.ErrNoRows {
			// If no user with that userID were found, assume username.
			whereClause = "users.username_safe = ?"
			p = common.SafeUsername(c.Query("u"))
		} else {
			// An user with that userID was found. Thus it's an userID.
			whereClause = "users.id = ?"
		}
	// u contains letters, so it's an username.
	default:
		whereClause = "users.username_safe = ?"
		p = common.SafeUsername(c.Query("u"))
	}
	return whereClause, p
}
