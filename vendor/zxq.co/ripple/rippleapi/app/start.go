package app

import (
	"fmt"

	"zxq.co/ripple/rippleapi/app/internals"
	"zxq.co/ripple/rippleapi/app/peppy"
	"zxq.co/ripple/rippleapi/app/v1"
	"zxq.co/ripple/rippleapi/common"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/serenize/snaker"
	"gopkg.in/redis.v5"
)

var (
	db    *sqlx.DB
	cf    common.Conf
	doggo *statsd.Client
	red   *redis.Client
)

var commonClusterfucks = map[string]string{
	"RegisteredOn": "register_datetime",
	"UsernameAKA":  "username_aka",
}

// Start begins taking HTTP connections.
func Start(conf common.Conf, dbO *sqlx.DB) *gin.Engine {
	db = dbO
	cf = conf

	db.MapperFunc(func(s string) string {
		if x, ok := commonClusterfucks[s]; ok {
			return x
		}
		return snaker.CamelToSnake(s)
	})

	setUpLimiter()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// sentry
	if conf.SentryDSN != "" {
		ravenClient, err := raven.New(conf.SentryDSN)
		ravenClient.SetRelease(common.Version)
		if err != nil {
			fmt.Println(err)
		} else {
			r.Use(Recovery(ravenClient, false))
		}
	}

	// datadog
	var err error
	doggo, err = statsd.New("127.0.0.1:8125")
	if err != nil {
		fmt.Println(err)
	}
	doggo.Namespace = "api."
	r.Use(func(c *gin.Context) {
		doggo.Incr("requests", nil, 1)
	})

	// redis
	red = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	// token updater
	go tokenUpdater(db)

	api := r.Group("/api")
	{
		p := api.Group("/")
		{
			p.GET("/get_user", PeppyMethod(peppy.GetUser))
			p.GET("/get_match", PeppyMethod(peppy.GetMatch))
			p.GET("/get_user_recent", PeppyMethod(peppy.GetUserRecent))
			p.GET("/get_user_best", PeppyMethod(peppy.GetUserBest))
			p.GET("/get_scores", PeppyMethod(peppy.GetScores))
			p.GET("/get_beatmaps", PeppyMethod(peppy.GetBeatmap))
		}

		gv1 := api.Group("/v1")
		{
			gv1.POST("/tokens", Method(v1.TokenNewPOST))
			gv1.POST("/tokens/new", Method(v1.TokenNewPOST))
			gv1.POST("/tokens/self/delete", Method(v1.TokenSelfDeletePOST))

			// Auth-free API endpoints (public data)
			gv1.GET("/ping", Method(v1.PingGET))
			gv1.GET("/surprise_me", Method(v1.SurpriseMeGET))
			gv1.GET("/doc", Method(v1.DocGET))
			gv1.GET("/doc/content", Method(v1.DocContentGET))
			gv1.GET("/doc/rules", Method(v1.DocRulesGET))
			gv1.GET("/users", Method(v1.UsersGET))
			gv1.GET("/users/whatid", Method(v1.UserWhatsTheIDGET))
			gv1.GET("/users/full", Method(v1.UserFullGET))
			gv1.GET("/users/userpage", Method(v1.UserUserpageGET))
			gv1.GET("/users/lookup", Method(v1.UserLookupGET))
			gv1.GET("/users/scores/best", Method(v1.UserScoresBestGET))
			gv1.GET("/users/scores/recent", Method(v1.UserScoresRecentGET))
			gv1.GET("/badges", Method(v1.BadgesGET))
			gv1.GET("/beatmaps", Method(v1.BeatmapGET))
			gv1.GET("/leaderboard", Method(v1.LeaderboardGET))
			gv1.GET("/tokens", Method(v1.TokenGET))
			gv1.GET("/users/self", Method(v1.UserSelfGET))
			gv1.GET("/tokens/self", Method(v1.TokenSelfGET))
			gv1.GET("/blog/posts", Method(v1.BlogPostsGET))
			gv1.GET("/blog/posts/content", Method(v1.BlogPostsContentGET))
			gv1.GET("/scores", Method(v1.ScoresGET))
			gv1.GET("/beatmaps/rank_requests/status", Method(v1.BeatmapRankRequestsStatusGET))

			// ReadConfidential privilege required
			gv1.GET("/friends", Method(v1.FriendsGET, common.PrivilegeReadConfidential))
			gv1.GET("/friends/with", Method(v1.FriendsWithGET, common.PrivilegeReadConfidential))
			gv1.GET("/users/self/donor_info", Method(v1.UsersSelfDonorInfoGET, common.PrivilegeReadConfidential))
			gv1.GET("/users/self/favourite_mode", Method(v1.UsersSelfFavouriteModeGET, common.PrivilegeReadConfidential))
			gv1.GET("/users/self/settings", Method(v1.UsersSelfSettingsGET, common.PrivilegeReadConfidential))

			// Write privilege required
			gv1.POST("/friends/add", Method(v1.FriendsAddPOST, common.PrivilegeWrite))
			gv1.POST("/friends/del", Method(v1.FriendsDelPOST, common.PrivilegeWrite))
			gv1.POST("/users/self/settings", Method(v1.UsersSelfSettingsPOST, common.PrivilegeWrite))
			gv1.POST("/users/self/userpage", Method(v1.UserSelfUserpagePOST, common.PrivilegeWrite))
			gv1.POST("/beatmaps/rank_requests", Method(v1.BeatmapRankRequestsSubmitPOST, common.PrivilegeWrite))

			// Admin: beatmap
			gv1.POST("/beatmaps/set_status", Method(v1.BeatmapSetStatusPOST, common.PrivilegeBeatmap))
			gv1.GET("/beatmaps/ranked_frozen_full", Method(v1.BeatmapRankedFrozenFullGET, common.PrivilegeBeatmap))

			// Admin: user managing
			gv1.POST("/users/manage/set_allowed", Method(v1.UserManageSetAllowedPOST, common.PrivilegeManageUser))

			// M E T A
			// E     T    "wow thats so meta"
			// T     E                  -- the one who said "wow thats so meta"
			// A T E M
			gv1.GET("/meta/restart", Method(v1.MetaRestartGET, common.PrivilegeAPIMeta))
			gv1.GET("/meta/kill", Method(v1.MetaKillGET, common.PrivilegeAPIMeta))
			gv1.GET("/meta/up_since", Method(v1.MetaUpSinceGET, common.PrivilegeAPIMeta))
			gv1.GET("/meta/update", Method(v1.MetaUpdateGET, common.PrivilegeAPIMeta))

			// User Managing + meta
			gv1.POST("/tokens/fix_privileges", Method(v1.TokenFixPrivilegesPOST,
				common.PrivilegeManageUser, common.PrivilegeAPIMeta))

			// in the new osu-web, the old endpoints are also in /v1 it seems. So /shrug
			gv1.GET("/get_user", PeppyMethod(peppy.GetUser))
			gv1.GET("/get_match", PeppyMethod(peppy.GetMatch))
			gv1.GET("/get_user_recent", PeppyMethod(peppy.GetUserRecent))
			gv1.GET("/get_user_best", PeppyMethod(peppy.GetUserBest))
			gv1.GET("/get_scores", PeppyMethod(peppy.GetScores))
			gv1.GET("/get_beatmaps", PeppyMethod(peppy.GetBeatmap))
		}

		api.GET("/status", internals.Status)
	}

	r.NoRoute(v1.Handle404)

	return r
}
