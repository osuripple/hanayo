package funcmap

// IsTFAEnabled checks whether 2FA is enabled.
func IsTFAEnabled(u int) (int, error) {
	return TFAService.Enabled(u)
}
