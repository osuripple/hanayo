package fail

// Failure is an error that happens at the server level, and as such information
// about it should not be given.
type Failure string

func (f Failure) Error() string {
	return string(f)
}

// Error is a mistake done by the user, and as such information about it should
// be shown to them.
type Error string

func (e Error) Error() string {
	return string(e)
}
