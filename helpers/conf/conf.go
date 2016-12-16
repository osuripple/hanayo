package conf

import (
	"git.zxq.co/x/rs"
	"github.com/thehowl/conf"
)

// Conf is the config file for hanayo.
var Conf = struct {
	ListenTo string `description:"ip:port from which to take requests."`
	Unix     bool   `description:"Whether ListenTo is an unix socket."`

	DSN string `description:"MySQL server DSN"`

	CookieSecret string

	RedisEnable         bool `description:"Whether to use redis for sessions"`
	RedisMaxConnections int
	RedisNetwork        string
	RedisAddress        string
	RedisPassword       string

	AvatarURL     string
	BaseURL       string
	DiscordServer string

	API           string
	BanchoAPI     string
	APISecret     string
	BaseAPIPublic string

	IPAPI string

	Offline          bool   `description:"If this is true, files will be served from the local server instead of the CDN."`
	MainRippleFolder string `description:"Folder where all the non-go projects are contained, such as old-frontend, lets, ci-system."`
	AvatarsFolder    string `description:"location folder of avatars"`

	MailgunDomain        string
	MailgunPrivateAPIKey string
	MailgunPublicAPIKey  string
	MailgunFrom          string

	RecaptchaSite    string
	RecaptchaPrivate string

	DiscordOAuthID     string
	DiscordOAuthSecret string
	DonorBotURL        string
	DonorBotSecret     string

	SentryDSN string

	AnalyticsID string
}{
	ListenTo:         ":45221",
	CookieSecret:     rs.String(46),
	AvatarURL:        "https://a.ripple.moe",
	BaseURL:          "https://ripple.moe",
	BanchoAPI:        "https://c.ripple.moe",
	API:              "http://localhost:40001/api/v1/",
	APISecret:        "Potato",
	IPAPI:            "https://ip.zxq.co",
	DiscordServer:    "#",
	MainRippleFolder: "/home/ripple/ripple",
	MailgunFrom:      `"Ripple" <noreply@ripple.moe>`,
}

// Load loads the configuration file into Conf.
func Load(name string) error {
	if name == "" {
		name = "hanayo.conf"
	}
	return conf.Load(&Conf, name)
}

// Export exports the configuration file from Conf.
func Export(name string) error {
	if name == "" {
		name = "hanayo.conf"
	}
	return conf.Export(Conf, name)
}
