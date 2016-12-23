package funcmap

import (
	"strconv"
	"time"
)

var hanayoStarted = time.Now().UnixNano()

// UnixNano returns the UNIX timestamp of when hanayo was started in nanoseconds.
func UnixNano() string {
	return strconv.FormatInt(hanayoStarted, 10)
}

// CurrYear returns an int containing the current year.
func CurrYear() int {
	return time.Now().Year()
}
