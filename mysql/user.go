package mysql

import (
	"database/sql"
	"time"

	"git.zxq.co/ripple/hanayo"
	// Go away golint.
	_ "github.com/go-sql-driver/mysql"
)

const userFields = `
id, username, username_safe, password_md5, password_version, email,
privileges, register_datetime as registered, latest_activity, flags
`

// User retrieves an user knowing their ID.
func (s *ServiceProvider) User(id int) (*hanayo.User, error) {
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
func (s *ServiceProvider) UserByEmail(email string) (*hanayo.User, error) {
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
func (s *ServiceProvider) RegisterUser(u hanayo.User) (int, error) {
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
		return 0, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = s.DB.Exec(
		`INSERT INTO users_stats(id, username, user_color)
		VALUES (?, ?, 'black');`, lid, u.Username,
	)

	return int(lid), err
}

// GetCountry retrieves an user's country.
func (s *ServiceProvider) GetCountry(id int) (*string, error) {
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
func (s *ServiceProvider) SetCountry(id int, country string) error {
	_, err := s.DB.Exec("UPDATE users_stats SET country = ? WHERE id = ?",
		country, id)
	return err
}
