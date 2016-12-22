package mysql

import (
	"database/sql"
	"time"

	"git.zxq.co/ripple/hanayo"
	"git.zxq.co/ripple/hanayo/fail"
	"git.zxq.co/ripple/rippleapi/common"
	"github.com/jmoiron/sqlx"
)

// TFAService implements the TFAService as specified by hanayo.TFAService.
type TFAService struct {
	DB *sqlx.DB
}

func (s *TFAService) init() error {
	if s.DB == nil {
		if DB == nil {
			return fail.FailDBIsNil
		}
		s.DB = DB
	}
	return nil
}

// Enabled returns whether two factor authentication is enabled for a certain
// user.
func (s *TFAService) Enabled(id int) (int, error) {
	if err := s.init(); err != nil {
		return 0, err
	}

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
func (s *TFAService) TelegramVerify(u int, ip string, token string) error {
	if err := s.init(); err != nil {
		return err
	}

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
func (s *TFAService) TelegramCreate(u int, ip string, token string) error {
	if err := s.init(); err != nil {
		return err
	}

	_, err := s.DB.Exec(
		"INSERT INTO 2fa(userid, token, ip, expire, sent) VALUES (?, ?, ?, ?, 0);",
		u, token, ip, time.Now().Add(time.Hour).Unix(),
	)
	return err
}

// TOTPInfo retrieves information about TOTP setup.
func (s *TFAService) TOTPInfo(u int) (*hanayo.TOTPInfo, error) {
	if err := s.init(); err != nil {
		return nil, err
	}

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
