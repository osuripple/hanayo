package hanayo

import (
	"git.zxq.co/ripple/rippleapi/common"
)

// An User is simply an user on Ripple.
type User struct {
	ID              int
	Username        string
	UsernameSafe    string
	PasswordMD5     string
	PasswordVersion uint8
	Email           string
	Privileges      uint64
	Registered      common.UnixTimestamp
	LatestActivity  common.UnixTimestamp
	Flags           uint64
}

// UserService represents a service able to return user information.
type UserService interface {
	User(id int) (*User, error)
	UserByEmail(email string) (*User, error)
	RegisterUser(u User) error
	GetCountry(id int) (*string, error)
	SetCountry(id int, country string) error
}
