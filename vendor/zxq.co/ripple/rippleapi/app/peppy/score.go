package peppy

import (
	"database/sql"
	"strconv"
	"strings"

	"zxq.co/ripple/rippleapi/common"

	"zxq.co/x/getrank"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gopkg.in/thehowl/go-osuapi.v1"
)

// GetScores retrieve information about the top 100 scores of a specified beatmap.
func GetScores(c *gin.Context, db *sqlx.DB) {
	if c.Query("b") == "" {
		c.JSON(200, defaultResponse)
		return
	}
	var beatmapMD5 string
	err := db.Get(&beatmapMD5, "SELECT beatmap_md5 FROM beatmaps WHERE beatmap_id = ? LIMIT 1", c.Query("b"))
	switch {
	case err == sql.ErrNoRows:
		c.JSON(200, defaultResponse)
		return
	case err != nil:
		c.Error(err)
		c.JSON(200, defaultResponse)
		return
	}
	var sb = "scores.score"
	if rankable(c.Query("m")) {
		sb = "scores.pp"
	}
	var (
		extraWhere  string
		extraParams []interface{}
	)
	if c.Query("u") != "" {
		w, p := genUser(c, db)
		extraWhere = "AND " + w
		extraParams = append(extraParams, p)
	}
	mods := common.Int(c.Query("mods"))
	rows, err := db.Query(`
SELECT
	scores.id, scores.score, users.username, scores.300_count, scores.100_count,
	scores.50_count, scores.misses_count, scores.gekis_count, scores.katus_count,
	scores.max_combo, scores.full_combo, scores.mods, users.id, scores.time, scores.pp,
	scores.accuracy
FROM scores
INNER JOIN users ON users.id = scores.userid
WHERE scores.completed = '3'
  AND users.privileges & 1 > 0
  AND scores.beatmap_md5 = ?
  AND scores.play_mode = ?
  AND scores.mods & ? = ?
  `+extraWhere+`
ORDER BY `+sb+` DESC LIMIT `+strconv.Itoa(common.InString(1, c.Query("limit"), 100, 50)),
		append([]interface{}{beatmapMD5, genmodei(c.Query("m")), mods, mods}, extraParams...)...)
	if err != nil {
		c.Error(err)
		c.JSON(200, defaultResponse)
		return
	}
	var results []osuapi.GSScore
	for rows.Next() {
		var (
			s         osuapi.GSScore
			fullcombo bool
			mods      int
			date      common.UnixTimestamp
			accuracy  float64
		)
		err := rows.Scan(
			&s.ScoreID, &s.Score.Score, &s.Username, &s.Count300, &s.Count100,
			&s.Count50, &s.CountMiss, &s.CountGeki, &s.CountKatu,
			&s.MaxCombo, &fullcombo, &mods, &s.UserID, &date, &s.PP,
			&accuracy,
		)
		if err != nil {
			if err != sql.ErrNoRows {
				c.Error(err)
			}
			continue
		}
		s.FullCombo = osuapi.OsuBool(fullcombo)
		s.Mods = osuapi.Mods(mods)
		s.Date = osuapi.MySQLDate(date)
		s.Rank = strings.ToUpper(getrank.GetRank(osuapi.Mode(genmodei(c.Query("m"))), s.Mods,
			accuracy, s.Count300, s.Count100, s.Count50, s.CountMiss))
		results = append(results, s)
	}
	c.JSON(200, results)
	return
}
