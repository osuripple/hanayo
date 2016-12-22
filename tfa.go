package hanayo

// TFAService is a service that handles information about Ripple's Two Factor
// Authentication system.
type TFAService interface {
	Enabled(u int) (int, error)

	TelegramVerify(u int, ip string, token string) error // when verification failed, return ErrTFATelegramVerificationFailed
	TelegramCreate(u int, ip string, token string) error

	TOTPInfo(u int) (*TOTPInfo, error)
}

// TOTPInfo contains the information about a TOTP token stored in the database.
type TOTPInfo struct {
	UserID   int
	Secret   string
	Recovery string
	Enabled  bool
}
