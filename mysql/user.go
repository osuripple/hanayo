// Package mysql implements hanayo's domain types with MySQL and sqlx.
package mysql

import (
	"database/sql"
	"time"

	"git.zxq.co/ripple/hanayo"
	"git.zxq.co/ripple/hanayo/fail"
	// Go away golint.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB is the default database that will be used unless it's passed in the
// service.
var DB *sqlx.DB

// UserService implements the UserService as specified by hanayo.UserService.
type UserService struct {
	DB *sqlx.DB
}

func (s *UserService) init() error {
	if s.DB == nil {
		if DB == nil {
			return fail.ErrDBIsNil
		}
		s.DB = DB
	}
	return nil
}

const userFields = `
id, username, username_safe, password_md5, password_version, email,
privileges, register_datetime as registered, latest_activity, flags
`

// User retrieves an user knowing their ID.
func (s *UserService) User(id int) (*hanayo.User, error) {
	if err := s.init(); err != nil {
		return nil, err
	}

	var u hanayo.User
	err := s.DB.Get(&u, "SELECT "+userFields+" FROM users WHERE id = ?", id)
	switch err {
	case nil:
		return &u, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// UserByEmail retrieves an user knowing their email address.
func (s *UserService) UserByEmail(email string) (*hanayo.User, error) {
	if err := s.init(); err != nil {
		return nil, err
	}

	var u hanayo.User
	err := s.DB.Get(&u, "SELECT "+userFields+
		" FROM users WHERE email = ?", email)
	switch err {
	case nil:
		return &u, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// RegisterUser creates a new user in the tables users and users_stats.
func (s *UserService) RegisterUser(u hanayo.User) error {
	if err := s.init(); err != nil {
		return err
	}

	n := time.Now().Unix()
	res, err := s.DB.Exec(
		`INSERT INTO users(
			username, username_safe, password_md5, email,
			register_datetime, privileges, latest_activity, password_version
		) VALUES (
			?, ?, ?, ?,
			?, ?, ?, 2
		)`,
		u.Username, u.UsernameSafe, u.PasswordMD5, u.Email,
		n, u.Privileges, n,
	)
	if err != nil {
		return err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		`INSERT INTO users_stats(id, username, user_color)
		VALUES (?, ?, 'black');`, lid, u.Username,
	)

	return err
}

// GetCountry retrieves an user's country.
func (s *UserService) GetCountry(id int) (*string, error) {
	if err := s.init(); err != nil {
		return nil, err
	}

	var c string
	err := s.DB.Get(&c, "SELECT country FROM users_stats WHERE id = ?", id)
	switch err {
	case nil:
		return &c, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// SetCountry sets an user's country.
func (s *UserService) SetCountry(id int, country string) error {
	if err := s.init(); err != nil {
		return err
	}

	_, err := s.DB.Exec("UPDATE users_stats SET country = ? WHERE id = ?",
		country, id)
	return err
}
