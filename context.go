package main

type context struct {
	User  sessionUser
	Token string
}
type sessionUser struct {
	ID       int
	Username string
}
