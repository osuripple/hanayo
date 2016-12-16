// Package fail contains the failures and errors of Hanayo.
package fail

// These are the failures that can occur commonly in hanayo.
var (
	FailDBIsNil = Failure("no database has been set (DB is nil)")
)
