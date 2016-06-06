package main

// user represents the user a logged in visitor could be.
// I can't really find any better explanation.
type user struct {
	ID       int
	Username string
	Password string
	Allowed  int
	APIToken string
}
