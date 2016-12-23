package funcmap

import (
	"git.zxq.co/ripple/hanayo"
)

// These are a list of services that will be set by the HTTP server.
var (
	UserService          hanayo.UserService
	TFAService           hanayo.TFAService
	SystemSettingService hanayo.SystemSettingService
)
