package v1

import (
	"fmt"

	"zxq.co/ripple/ocl"
	"zxq.co/ripple/rippleapi/common"
)

type leaderboardUser struct {
	userData
	ChosenMode    modeData `json:"chosen_mode"`
	PlayStyle     int      `json:"play_style"`
	FavouriteMode int      `json:"favourite_mode"`
}

type leaderboardResponse struct {
	common.ResponseBase
	Users []leaderboardUser `json:"users"`
}

const lbUserQuery = `
SELECT
	users.id, users.username, users.register_datetime, users.privileges, users.latest_activity,

	users_stats.username_aka, users_stats.country,
	users_stats.play_style, users_stats.favourite_mode,

	users_stats.ranked_score_%[1]s, users_stats.total_score_%[1]s, users_stats.playcount_%[1]s,
	users_stats.replays_watched_%[1]s, users_stats.total_hits_%[1]s,
	users_stats.avg_accuracy_%[1]s, users_stats.pp_%[1]s, leaderboard_%[1]s.position as %[1]s_position
FROM leaderboard_%[1]s
INNER JOIN users ON users.id = leaderboard_%[1]s.user
INNER JOIN users_stats ON users_stats.id = leaderboard_%[1]s.user
%[2]s`

// LeaderboardGET gets the leaderboard.
func LeaderboardGET(md common.MethodData) common.CodeMessager {
	m := getMode(md.Query("mode"))
	w := &common.WhereClause{
		Clause: "WHERE " + md.User.OnlyUserPublic(md.HasQuery("see_everything")),
	}
	w.Where("users_stats.country = ?", md.Query("country"))
	// Admins may not want to see banned users on the leaderboard.
	// This is the default setting. In case they do, they have to activate see_everything.
	query := fmt.Sprintf(lbUserQuery, m, w.Clause+
		` ORDER BY leaderboard_`+m+`.position `+common.Paginate(md.Query("p"), md.Query("l"), 500))
	rows, err := md.DB.Query(query, w.Params...)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var resp leaderboardResponse
	for rows.Next() {
		var u leaderboardUser
		err := rows.Scan(
			&u.ID, &u.Username, &u.RegisteredOn, &u.Privileges, &u.LatestActivity,

			&u.UsernameAKA, &u.Country, &u.PlayStyle, &u.FavouriteMode,

			&u.ChosenMode.RankedScore, &u.ChosenMode.TotalScore, &u.ChosenMode.PlayCount,
			&u.ChosenMode.ReplaysWatched, &u.ChosenMode.TotalHits,
			&u.ChosenMode.Accuracy, &u.ChosenMode.PP, &u.ChosenMode.GlobalLeaderboardRank,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		u.ChosenMode.Level = ocl.GetLevelPrecise(int64(u.ChosenMode.TotalScore))
		resp.Users = append(resp.Users, u)
	}
	resp.Code = 200
	return resp
}
