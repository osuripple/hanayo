package hanayo

// SystemSettingService is a service providing infomation about the system
// settings.
type SystemSettingService interface {
	Setting(settings ...string) (map[string]SystemSetting, error)
	AllSettings() (map[string]SystemSetting, error)
}

// SystemSetting represents a system setting on Ripple.
type SystemSetting struct {
	Name   string
	Int    int
	String string
}
