package main

// about using johnniedoe/contrib/gzip:
// johnniedoe's fork fixes a critical issue for which .String resulted in
// an ERR_DECODING_FAILED. This is an actual pull request on the contrib
// repo, but apparently, gin is dead.

import (
	"encoding/gob"
	"fmt"

	"git.zxq.co/ripple/schiavolib"
	"git.zxq.co/x/rs"
	"github.com/fatih/structs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/johnniedoe/contrib/gzip"
	"github.com/thehowl/conf"
	"github.com/thehowl/qsql"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// version is the version of hanayo
const version = "0.5.1b"

var (
	config struct {
		ListenTo string `description:"ip:port from which to take requests."`
		Unix     bool   `description:"Whether ListenTo is an unix socket."`

		DSN string `description:"MySQL server DSN"`

		CookieSecret string

		RedisEnable         bool
		RedisMaxConnections int
		RedisNetwork        string
		RedisAddress        string
		RedisPassword       string

		AvatarURL     string
		BaseURL       string
		DiscordServer string

		API       string
		BanchoAPI string
		APISecret string

		IP_API string

		Offline          bool   `json:"If this is true, files will be served from the local server instead of the CDN."`
		MainRippleFolder string `json:"Folder where all the non-go projects are contained, such as old-frontend, lets, ci-system."`

		MailgunDomain        string
		MailgunPrivateAPIKey string
		MailgunPublicAPIKey  string
		MailgunFrom          string

		RecaptchaSite    string
		RecaptchaPrivate string
	}
	configMap map[string]interface{}
	db        *sqlx.DB
	qb        *qsql.DB
	mg        mailgun.Mailgun
)

func main() {
	fmt.Println("hanayo v" + version)

	err := conf.Load(&config, "hanayo.conf")
	switch err {
	case nil:
		// carry on
	case conf.ErrNoFile:
		conf.Export(config, "hanayo.conf")
		fmt.Println("The configuration file was not found. We created one for you.")
		return
	default:
		panic(err)
	}

	var configDefaults = map[*string]string{
		&config.ListenTo:         ":45221",
		&config.CookieSecret:     rs.String(46),
		&config.AvatarURL:        "https://a.ripple.moe",
		&config.BaseURL:          "https://ripple.moe",
		&config.BanchoAPI:        "https://c.ripple.moe",
		&config.API:              "http://localhost:40001/api/v1/",
		&config.APISecret:        "Potato",
		&config.IP_API:           "https://ip.zxq.co",
		&config.DiscordServer:    "#",
		&config.MainRippleFolder: "/home/ripple/ripple",
		&config.MailgunFrom:      `"Ripple" <noreply@ripple.moe>`,
	}
	for key, value := range configDefaults {
		if *key == "" {
			*key = value
		}
	}

	configMap = structs.Map(config)

	// initialise db
	db, err = sqlx.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}
	qb = qsql.New(db.DB)
	if err != nil {
		panic(err)
	}

	// initialise mailgun
	mg = mailgun.NewMailgun(
		config.MailgunDomain,
		config.MailgunPrivateAPIKey,
		config.MailgunPublicAPIKey,
	)

	if gin.Mode() == gin.DebugMode {
		fmt.Println("Development environment detected. Starting fsnotify on template folder...")
		err := reloader()
		if err != nil {
			fmt.Println(err)
		}
	}

	schiavo.Prefix = "hanayo"
	schiavo.Bunker.Send(fmt.Sprintf("STARTUATO, mode: %s", gin.Mode()))

	// even if it's not release, we say that it's release
	// so that gin doesn't spam
	gin.SetMode(gin.ReleaseMode)

	gobRegisters := []interface{}{
		[]message{},
		errorMessage{},
		infoMessage{},
		neutralMessage{},
		warningMessage{},
		successMessage{},
	}
	for _, el := range gobRegisters {
		gob.Register(el)
	}

	fmt.Println("Importing templates...")
	loadTemplates("")

	fmt.Println("Setting up rate limiter...")
	setUpLimiter()

	fmt.Println("Exporting configuration...")

	conf.Export(config, "hanayo.conf")

	httpLoop()
}

func httpLoop() {
	for {
		e := generateEngine()
		fmt.Println("Starting webserver...")
		if !startuato(e) {
			break
		}
	}
}

func generateEngine() *gin.Engine {
	fmt.Println("Starting session system...")
	var store sessions.Store
	if config.RedisMaxConnections != 0 {
		var err error
		store, err = sessions.NewRedisStore(
			config.RedisMaxConnections,
			config.RedisNetwork,
			config.RedisAddress,
			config.RedisPassword,
			[]byte(config.CookieSecret),
		)
		if err != nil {
			fmt.Println(err)
			store = sessions.NewCookieStore([]byte(config.CookieSecret))
		}
	} else {
		store = sessions.NewCookieStore([]byte(config.CookieSecret))
	}

	r := gin.Default()

	r.Use(
		gzip.Gzip(gzip.DefaultCompression),
		sessions.Sessions("session", store),
		sessionInitializer(),
		rateLimiter(false),
		twoFALock,
	)

	r.Static("/static", "static")
	r.StaticFile("/favicon.ico", "static/favicon.ico")

	r.POST("/login", loginSubmit)
	r.GET("/logout", logout)

	r.GET("/register", register)
	r.POST("/register", registerSubmit)
	r.GET("/register/verify", verifyAccount)
	r.GET("/register/welcome", welcome)

	r.GET("/u/:user", userProfile)

	r.POST("/pwreset", passwordReset)
	r.GET("/pwreset/continue", passwordResetContinue)
	r.POST("/pwreset/continue", passwordResetContinueSubmit)

	r.GET("/2fa_gateway", tfaGateway)
	r.GET("/2fa_gateway/clear", clear2fa)
	r.GET("/2fa_gateway/verify", verify2fa)

	r.GET("/irc/generate", ircGenToken)

	r.POST("/settings/userpage/parse", parseBBCode)

	loadSimplePages(r)

	r.NoRoute(notFound)

	return r
}

const alwaysRespondText = `Ooops! Looks like something went really wrong while trying to process your request.
Perhaps report this to a Ripple developer?
Retrying doing again what you were trying to do might work, too.`
