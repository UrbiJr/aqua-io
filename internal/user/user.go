package user

// User is a type associated with users of this software.
type User struct {
	Email           string
	Discord         any
	Username        string
	LicenseKey      string
	ExpiresAt       int64
	PersistentLogin bool
	Settings        *Settings
	*ProfileManager
}

// NewUser creates an instance of the User class
func NewUser(email, uid, username, licenseKey string, discord any, expiresAt int64) *User {
	return &User{
		Email:      email,
		Discord:    discord,
		Username:   username,
		LicenseKey: licenseKey,
		ExpiresAt:  expiresAt,
	}
}
