package main

import "git.zxq.co/ripple/rippleapi/common"

type context struct {
	User  sessionUser
	Token string
}
type sessionUser struct {
	ID         int
	Username   string
	Privileges common.UserPrivileges
}
