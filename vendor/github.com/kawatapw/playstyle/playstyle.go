// Package playstyle provides an enum for Ripple's playstyles.
package playstyle

import "strings"

// PlayStyle is a bitwise enum containing the instruments a Ripple user likes
// to play with.
type PlayStyle int

// various playstyles on ripple.
const (
	Mouse int = 1 << iota
	Tablet
	Keyboard
	Touchscreen
	Relax
	Timewarp
	AimAssist
	Unmod
	EZ
	TrueHD
)

// Styles are string representations of the various playstyles someone can have.
var Styles = [...]string{
	"Mouse",
	"Tablet",
	"Keyboard",
	"Touchscreen",
	"Relax",
	"Timewarp",
	"Aim Assist",
	"Unmod",
	"EZ",
	"TrueHD",
}

// String is the string representation of a playstyle.
func (p PlayStyle) String() string {
	var parts []string

	i := int(p)
	for k, v := range Styles {
		if i&(1<<uint(k)) > 0 {
			parts = append(parts, v)
		}
	}

	return strings.Join(parts, ", ")
}
