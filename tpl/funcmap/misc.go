package funcmap

import (
	"git.zxq.co/ripple/hanayo"
)

// SystemSettings retrieves system settings from the database.
func SystemSettings(names ...string) (map[string]hanayo.SystemSetting, error) {
	return SystemSettingService.Setting(names...)
}

// Version returns the version of hanayo.
func Version() string {
	return hanayo.Version
}
