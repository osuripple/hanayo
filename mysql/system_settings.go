package mysql

import "git.zxq.co/ripple/hanayo"

// Setting retrieves settings from the database.
func (s *ServiceProvider) Setting(settings ...string) (map[string]hanayo.SystemSetting, error) {
	if len(settings) == 0 {
		return nil, nil
	}
	var (
		in   = ""
		args []interface{}
	)
	for _, s := range settings {
		in += "?, "
		args = append(args, s)
	}
	in = in[:len(in)-2]

	ss := make([]hanayo.SystemSetting, 0, len(settings))
	err := s.DB.Select(
		&ss,
		"SELECT name, value_int as `int`, value_string as `string` "+
			"WHERE name IN ("+in+")",
		args...,
	)
	if err != nil {
		return nil, err
	}

	return settingsToMap(ss), nil
}

// AllSettings retrieves all the system settings from the database.
func (s *ServiceProvider) AllSettings() (map[string]hanayo.SystemSetting, error) {
	ss := make([]hanayo.SystemSetting, 0, 10)
	err := s.DB.Select(
		&ss, "SELECT name, value_int as `int`, value_string as `string`",
	)
	if err != nil {
		return nil, err
	}

	return settingsToMap(ss), nil
}

func settingsToMap(ss []hanayo.SystemSetting) map[string]hanayo.SystemSetting {
	m := make(map[string]hanayo.SystemSetting, len(ss))
	for _, s := range ss {
		m[s.Name] = s
	}
	return m
}
