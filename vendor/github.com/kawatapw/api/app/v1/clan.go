package v1

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"github.com/kawatapw/api/common"
)

type singleClan struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Icon        string `json:"icon"`
}

type multiClanData struct {
	common.ResponseBase
	Clans []singleClan `json:"clans"`
}

// clansGET retrieves all the clans on this ripple instance.
func ClansGET(md common.MethodData) common.CodeMessager {
	var (
		r    multiClanData
		rows *sql.Rows
		err  error
	)
	if md.Query("id") != "" {
		rows, err = md.DB.Query("SELECT id, name, description, tag, icon FROM clans WHERE id = ? LIMIT 1", md.Query("id"))
	} else {
		rows, err = md.DB.Query("SELECT id, name, description, tag, icon FROM clans " + common.Paginate(md.Query("p"), md.Query("l"), 50))
	}
	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nc := singleClan{}
		err = rows.Scan(&nc.ID, &nc.Name, &nc.Description, &nc.Tag, &nc.Icon)
		if err != nil {
			md.Err(err)
		}
		r.Clans = append(r.Clans, nc)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}

type clanMembersData struct {
	common.ResponseBase
	Members []userNotFullResponse `json:"members"`
}

// get total stats of clan. later.
type totalStats struct {
	common.ResponseBase
	ClanID     int      `json:"id"`
	ChosenMode modeData `json:"chosen_mode"`
	Rank       int      `json:"rank"`
}
type clanLbSingle struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tag         string   `json:"tag"`
	Icon        string   `json:"icon"`
	ChosenMode  modeData `json:"chosen_mode"`
	Rank        int      `json:"rank"`
}

type megaStats struct {
	common.ResponseBase
	Clans []clanLbSingle `json:"clans"`
}

func AllClanStatsGET(md common.MethodData) common.CodeMessager {
	var (
		r    megaStats
		rows *sql.Rows
		err  error
	)
	rows, err = md.DB.Query("SELECT id, name, description, tag, icon FROM clans")

	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nc := clanLbSingle{}
		err = rows.Scan(&nc.ID, &nc.Name, &nc.Description, &nc.Tag, &nc.Icon)
		fmt.Println(rows)
		fmt.Println(&nc.Tag)
		if err != nil {
			md.Err(err)
		}
		nc.ChosenMode.PP = 0
		r.Clans = append(r.Clans, nc)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	// anyone who ever looks into this, yes, i need to kill myself. ~Flame
	// yeah.. yeah.. i see flame ~Hazuki-san
	m, brr := strconv.ParseInt(string(md.Query("m")[19]), 10, 64)

	if brr != nil {
		fmt.Println(brr)
		m = 0
	}
	n := "std"
	if m == 1 {
		n = "taiko"
	} else if m == 2 {
		n = "ctb"
	} else if m == 3 {
		n = "mania"
	} else {
		n = "std"
	}
	fmt.Println(n)

	for i := 0; i < len(r.Clans); i++ {
		var members clanMembersData

		rid := r.Clans[i].ID

		err := md.DB.Select(&members.Members, `SELECT users.id, users.username, users.register_datetime, users.privileges,
		latest_activity, users_stats.username_aka,
		
		users_stats.country, users_stats.user_color,
		users_stats.ranked_score_std, users_stats.total_score_std, users_stats.pp_std, users_stats.playcount_std, users_stats.replays_watched_std, users_stats.total_hits_std,
		users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.pp_taiko, users_stats.playcount_taiko, users_stats.replays_watched_taiko, users_stats.total_hits_taiko,
		users_stats.ranked_score_ctb, users_stats.total_score_ctb, users_stats.pp_ctb, users_stats.playcount_ctb, users_stats.replays_watched_ctb, users_stats.total_hits_ctb,
		users_stats.ranked_score_mania, users_stats.total_score_mania, users_stats.pp_mania, users_stats.playcount_mania, users_stats.replays_watched_mania, users_stats.total_hits_mania
		
		FROM user_clans uc
		INNER JOIN users
		ON users.id = uc.user
		INNER JOIN users_stats ON users_stats.id = uc.user
		WHERE clan = ? AND privileges & 1 = 1
		`, rid)

		if err != nil {
			fmt.Println(err)
		}

		members.Code = 200

		if n == "std" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpStd
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreStd
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreStd
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountStd
			}
		} else if n == "taiko" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpTaiko
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreTaiko
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreTaiko
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountTaiko
			}
		} else if n == "ctb" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpCtb
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreCtb
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreCtb
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountCtb
			}
		} else if n == "mania" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpMania
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreMania
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreMania
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountMania
			}
		}
		r.Clans[i].ChosenMode.PP = (r.Clans[i].ChosenMode.PP / (len(members.Members) + 1))
	}

	sort.Slice(r.Clans, func(i, j int) bool {
		return r.Clans[i].ChosenMode.PP > r.Clans[j].ChosenMode.PP
	})

	for i := 0; i < len(r.Clans); i++ {
		r.Clans[i].Rank = i + 1
	}

	return r
}

func TotalClanStatsGET(md common.MethodData) common.CodeMessager {
	var (
		r    megaStats
		rows *sql.Rows
		err  error
	)
	rows, err = md.DB.Query("SELECT id, name, description, icon FROM clans")

	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nc := clanLbSingle{}
		err = rows.Scan(&nc.ID, &nc.Name, &nc.Description, &nc.Icon)
		if err != nil {
			md.Err(err)
		}
		nc.ChosenMode.PP = 0
		r.Clans = append(r.Clans, nc)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200

	id := common.Int(md.Query("id"))
	if id == 0 {
		return ErrMissingField("id")
	}
	//Uh... well... ;-;
	m, brr := strconv.ParseInt(string(md.Query("m")[11]), 10, 64)
	if brr != nil {
		fmt.Println(brr)
	}

	n := "std"
	if m == 1 {
		n = "taiko"
	} else if m == 2 {
		n = "ctb"
	} else if m == 3 {
		n = "mania"
	} else {
		n = "std"
	}
	fmt.Println(n)

	for i := 0; i < len(r.Clans); i++ {
		var members clanMembersData

		rid := r.Clans[i].ID

		err := md.DB.Select(&members.Members, `SELECT users.id, users.username, users.register_datetime, users.privileges,
		latest_activity, users_stats.username_aka,
		
		users_stats.country, users_stats.user_color,
		users_stats.ranked_score_std, users_stats.total_score_std, users_stats.pp_std, users_stats.playcount_std, users_stats.replays_watched_std, users_stats.total_hits_std,
		users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.pp_taiko, users_stats.playcount_taiko, users_stats.replays_watched_taiko, users_stats.total_hits_taiko,
		users_stats.ranked_score_ctb, users_stats.total_score_ctb, users_stats.pp_ctb, users_stats.playcount_ctb, users_stats.replays_watched_ctb, users_stats.total_hits_ctb,
		users_stats.ranked_score_mania, users_stats.total_score_mania, users_stats.pp_mania, users_stats.playcount_mania, users_stats.replays_watched_mania, users_stats.total_hits_mania
		
		FROM user_clans uc
		INNER JOIN users
		ON users.id = uc.user
		INNER JOIN users_stats ON users_stats.id = uc.user
		WHERE clan = ? AND privileges & 1 = 1
		`, rid)

		if err != nil {
			fmt.Println(err)
		}

		members.Code = 200

		if n == "std" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpStd
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreStd
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreStd
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountStd
				r.Clans[i].ChosenMode.ReplaysWatched = r.Clans[i].ChosenMode.ReplaysWatched + members.Members[u].ReplaysWatchedStd
				r.Clans[i].ChosenMode.TotalHits = r.Clans[i].ChosenMode.TotalHits + members.Members[u].TotalHitsStd
			}
		} else if n == "taiko" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpTaiko
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreTaiko
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreTaiko
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountTaiko
				r.Clans[i].ChosenMode.ReplaysWatched = r.Clans[i].ChosenMode.ReplaysWatched + members.Members[u].ReplaysWatchedTaiko
				r.Clans[i].ChosenMode.TotalHits = r.Clans[i].ChosenMode.TotalHits + members.Members[u].TotalHitsTaiko
			}
		} else if n == "ctb" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpCtb
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreCtb
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreCtb
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountCtb
				r.Clans[i].ChosenMode.ReplaysWatched = r.Clans[i].ChosenMode.ReplaysWatched + members.Members[u].ReplaysWatchedCtb
				r.Clans[i].ChosenMode.TotalHits = r.Clans[i].ChosenMode.TotalHits + members.Members[u].TotalHitsStd
			}
		} else if n == "mania" {
			for u := 0; u < len(members.Members); u++ {
				r.Clans[i].ChosenMode.PP = r.Clans[i].ChosenMode.PP + members.Members[u].PpMania
				r.Clans[i].ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore + members.Members[u].RankedScoreMania
				r.Clans[i].ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore + members.Members[u].TotalScoreMania
				r.Clans[i].ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount + members.Members[u].PlaycountMania
				r.Clans[i].ChosenMode.ReplaysWatched = r.Clans[i].ChosenMode.ReplaysWatched + members.Members[u].ReplaysWatchedMania
				r.Clans[i].ChosenMode.TotalHits = r.Clans[i].ChosenMode.TotalHits + members.Members[u].TotalHitsMania
			}
		}
		r.Clans[i].ChosenMode.PP = (r.Clans[i].ChosenMode.PP / (len(members.Members) + 1))
	}

	sort.Slice(r.Clans, func(i, j int) bool {
		return r.Clans[i].ChosenMode.PP > r.Clans[j].ChosenMode.PP
	})

	for i := 0; i < len(r.Clans); i++ {
		r.Clans[i].Rank = i + 1
	}
	b := totalStats{}
	for i := 0; i < len(r.Clans); i++ {
		if r.Clans[i].ID == id {
			b.ClanID = id
			b.ChosenMode.PP = r.Clans[i].ChosenMode.PP
			b.ChosenMode.RankedScore = r.Clans[i].ChosenMode.RankedScore
			b.ChosenMode.TotalScore = r.Clans[i].ChosenMode.TotalScore
			b.ChosenMode.PlayCount = r.Clans[i].ChosenMode.PlayCount
			b.ChosenMode.ReplaysWatched = r.Clans[i].ChosenMode.ReplaysWatched
			b.ChosenMode.TotalHits = r.Clans[i].ChosenMode.TotalHits
			b.Rank = r.Clans[i].Rank
			b.Code = 200
		}
	}

	return b
}

type isClanData struct {
	Clan  int `json:"clan"`
	User  int `json:"user"`
	Perms int `json:"perms"`
}

type isClan struct {
	common.ResponseBase
	Clan isClanData `json:"clan"`
}

func IsInClanGET(md common.MethodData) common.CodeMessager {
	ui := md.Query("uid")

	if ui == "0" {
		return ErrMissingField("uid")
	}

	var r isClan
	rows, err := md.DB.Query("SELECT user, clan, perms FROM user_clans WHERE user = ?", ui)

	if err != nil {
		md.Err(err)
		return Err500
	}

	defer rows.Close()
	for rows.Next() {
		nc := isClanData{}
		err = rows.Scan(&nc.User, &nc.Clan, &nc.Perms)
		if err != nil {
			md.Err(err)
		}
		r.Clan = nc
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}

type imFoolish struct {
	common.ResponseBase
	Invite string `json:"invite"`
}
type adminClan struct {
	Id int `json:"user"`
	Perms   int `json:"perms"`
}

func ClanInviteGET(md common.MethodData) common.CodeMessager {
	// big perms check lol ok
	n := common.Int(md.Query("id"))
	adminFoolish := adminClan{}

	var r imFoolish
	var clan int
	// get user clan, then get invite
	md.DB.QueryRow("SELECT user, clan, perms FROM user_clans WHERE user = ? LIMIT 1", n).Scan(&adminFoolish.Id, &clan, &adminFoolish.Perms)
	if adminFoolish.Perms < 8 || adminFoolish.Id != md.ID() {
		return common.SimpleResponse(500, "You are not the admin of the clan")
	}
	row := md.DB.QueryRow("SELECT invite FROM clans_invites WHERE clan = ? LIMIT 1", clan).Scan(&r.Invite)
	if row != nil {
		fmt.Println(row)
	}
	return r
}

// ClanMembersGET retrieves the people who are in a certain clan.
func ClanMembersGET(md common.MethodData) common.CodeMessager {
	i := common.Int(md.Query("id"))
	if i == 0 {
		return ErrMissingField("id")
	}
	r := common.Int(md.Query("r"))
	if r == 0 {
		var members clanMembersData

		err := md.DB.Select(&members.Members, `SELECT users.id, users.username, users.register_datetime, users.privileges,
	latest_activity, users_stats.username_aka,
	
	users_stats.country, users_stats.user_color,
	users_stats.ranked_score_std, users_stats.total_score_std, users_stats.pp_std, users_stats.playcount_std, users_stats.replays_watched_std, users_stats.total_hits_std,
	users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.pp_taiko, users_stats.playcount_taiko, users_stats.replays_watched_taiko, users_stats.total_hits_taiko
	
FROM user_clans uc
INNER JOIN users
ON users.id = uc.user
INNER JOIN users_stats ON users_stats.id = uc.user
WHERE clan = ?
ORDER BY id ASC `, i)

		if err != nil {
			md.Err(err)
			return Err500
		}

		members.Code = 200
		return members
	} else {
		var members clanMembersData

		err := md.DB.Select(&members.Members, `SELECT users.id, users.username, users.register_datetime, users.privileges,
	latest_activity, users_stats.username_aka,
	
	users_stats.country, users_stats.user_color,
	users_stats.ranked_score_std, users_stats.total_score_std, users_stats.pp_std, users_stats.playcount_std, users_stats.replays_watched_std,
	users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.pp_taiko, users_stats.playcount_taiko, users_stats.replays_watched_taiko
	
FROM user_clans uc
INNER JOIN users
ON users.id = uc.user
INNER JOIN users_stats ON users_stats.id = uc.user
WHERE clan = ? AND perms = ?
ORDER BY id ASC `, i, r)

		if err != nil {
			md.Err(err)
			return Err500
		}

		members.Code = 200
		return members
	}
}

// Zunhapan likes this.