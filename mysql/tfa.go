package mysql

import (
	"database/sql"
	"time"

	"git.zxq.co/ripple/hanayo"
	"git.zxq.co/ripple/hanayo/fail"
	"git.zxq.co/ripple/rippleapi/common"
)

// Enabled returns whether two factor authentication is enabled for a certain
// user.
func (s *ServiceProvider) Enabled(id int) (int, error) {
	var u int
	err := s.DB.Get(
		&u,
		"SELECT IFNULL((SELECT 1 FROM 2fa_telegram WHERE userid = ?), 0) | "+
			"IFNULL((SELECT 2 FROM 2fa_totp WHERE userid = ? AND enabled = 1), 0) "+
			"as x", id, id,
	)
	return u, err
}

// TelegramVerify checks the token the user provided is valid.
func (s *ServiceProvider) TelegramVerify(u int, ip string, token string) error {
	var t common.UnixTimestamp
	err := s.DB.Get(&t,
		"SELECT expire FROM 2fa WHERE userid = ? AND ip = ? AND token = ?",
		u, ip, token,
	)

	switch err {
	case nil:
		// move on
	case sql.ErrNoRows:
		return fail.ErrTFATelegramVerificationFailed
	default:
		return err
	}

	if time.Now().After(time.Time(t)) {
		return fail.ErrTFATelegramVerificationFailed
	}

	return nil
}

// TelegramCreate inserts a new token in the 2fa table containing the Telegram
// token.
func (s *ServiceProvider) TelegramCreate(u int, ip string, token string) error {
	_, err := s.DB.Exec(
		"INSERT INTO 2fa(userid, token, ip, expire, sent) VALUES (?, ?, ?, ?, 0);",
		u, token, ip, time.Now().Add(time.Hour).Unix(),
	)
	return err
}

// TOTPInfo retrieves information about TOTP setup.
func (s *ServiceProvider) TOTPInfo(u int) (*hanayo.TOTPInfo, error) {
	var ti hanayo.TOTPInfo
	err := s.DB.Get(&ti, "SELECT userid, secret, recovery, enabled FROM 2fa_totp WHERE userid = ?", u)
	switch err {
	case nil:
		return &ti, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
