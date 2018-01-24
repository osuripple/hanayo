package v1

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/thehowl/go-osuapi.v1"
	"zxq.co/ripple/rippleapi/common"
	"zxq.co/x/getrank"
)

// Score is a score done on Ripple.
type Score struct {
	ID         int                  `json:"id"`
	BeatmapMD5 string               `json:"beatmap_md5"`
	Score      int64                `json:"score"`
	MaxCombo   int                  `json:"max_combo"`
	FullCombo  bool                 `json:"full_combo"`
	Mods       int                  `json:"mods"`
	Count300   int                  `json:"count_300"`
	Count100   int                  `json:"count_100"`
	Count50    int                  `json:"count_50"`
	CountGeki  int                  `json:"count_geki"`
	CountKatu  int                  `json:"count_katu"`
	CountMiss  int                  `json:"count_miss"`
	Time       common.UnixTimestamp `json:"time"`
	PlayMode   int                  `json:"play_mode"`
	Accuracy   float64              `json:"accuracy"`
	PP         float32              `json:"pp"`
	Rank       string               `json:"rank"`
	Completed  int                  `json:"completed"`
}

// beatmapScore is to differentiate from userScore, as beatmapScore contains
// also an user, while userScore contains the beatmap.
type beatmapScore struct {
	Score
	User userData `json:"user"`
}

type scoresResponse struct {
	common.ResponseBase
	Scores []beatmapScore `json:"scores"`
}

// ScoresGET retrieves the top scores for a certain beatmap.
func ScoresGET(md common.MethodData) common.CodeMessager {
	var (
		beatmapMD5 string
		r          scoresResponse
	)
	switch {
	case md.Query("md5") != "":
		beatmapMD5 = md.Query("md5")
	case md.Query("b") != "":
		err := md.DB.Get(&beatmapMD5, "SELECT beatmap_md5 FROM beatmaps WHERE beatmap_id = ? LIMIT 1", md.Query("b"))
		switch {
		case err == sql.ErrNoRows:
			r.Code = 200
			return r
		case err != nil:
			md.Err(err)
			return Err500
		}
	default:
		return ErrMissingField("md5|b")
	}

	sort := common.Sort(md, common.SortConfiguration{
		Default: "scores.pp DESC, scores.score DESC",
		Table:   "scores",
		Allowed: []string{"pp", "score", "accuracy", "id"},
	})

	rows, err := md.DB.Query(`
SELECT
	scores.id, scores.beatmap_md5, scores.score,
	scores.max_combo, scores.full_combo, scores.mods,
	scores.300_count, scores.100_count, scores.50_count,
	scores.gekis_count, scores.katus_count, scores.misses_count,
	scores.time, scores.play_mode, scores.accuracy, scores.pp,
	scores.completed,

	users.id, users.username, users.register_datetime, users.privileges,
	users.latest_activity, users_stats.username_aka, users_stats.country
FROM scores
INNER JOIN users ON users.id = scores.userid
INNER JOIN users_stats ON users_stats.id = scores.userid
WHERE scores.beatmap_md5 = ? AND scores.completed = '3' AND `+md.User.OnlyUserPublic(false)+
		` `+genModeClause(md)+`
`+sort+common.Paginate(md.Query("p"), md.Query("l"), 100), beatmapMD5)
	if err != nil {
		md.Err(err)
		return Err500
	}
	for rows.Next() {
		var (
			s beatmapScore
			u userData
		)
		err := rows.Scan(
			&s.ID, &s.BeatmapMD5, &s.Score.Score,
			&s.MaxCombo, &s.FullCombo, &s.Mods,
			&s.Count300, &s.Count100, &s.Count50,
			&s.CountGeki, &s.CountKatu, &s.CountMiss,
			&s.Time, &s.PlayMode, &s.Accuracy, &s.PP,
			&s.Completed,

			&u.ID, &u.Username, &u.RegisteredOn, &u.Privileges,
			&u.LatestActivity, &u.UsernameAKA, &u.Country,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		s.User = u
		s.Rank = strings.ToUpper(getrank.GetRank(
			osuapi.Mode(s.PlayMode),
			osuapi.Mods(s.Mods),
			s.Accuracy,
			s.Count300,
			s.Count100,
			s.Count50,
			s.CountMiss,
		))
		r.Scores = append(r.Scores, s)
	}
	r.Code = 200
	return r
}

func getMode(m string) string {
	switch m {
	case "1":
		return "taiko"
	case "2":
		return "ctb"
	case "3":
		return "mania"
	default:
		return "std"
	}
}

func genModeClause(md common.MethodData) string {
	var modeClause string
	if md.Query("mode") != "" {
		m, err := strconv.Atoi(md.Query("mode"))
		if err == nil && m >= 0 && m <= 3 {
			modeClause = fmt.Sprintf("AND scores.play_mode = '%d'", m)
		}
	}
	return modeClause
}
