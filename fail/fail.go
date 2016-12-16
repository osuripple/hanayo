// Package fail contains the failures and errors of Hanayo.
package fail

import (
	"errors"
)

// These are the errors that can occur commonly in hanayo.
var (
	ErrDBIsNil = errors.New("no database has been set (DB is nil)")
)
