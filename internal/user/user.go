package user

// User is a type associated with users of this software.
type User struct {
	LicenseKey string
	Email      string
	Username   string
	UID        string
	Profiles   []Profile
}

// NewUser creates an instance of the User class
func NewUser(licenseKey string, email string, uid string, username string) *User {
	return &User{
		LicenseKey: licenseKey,
		Email:      email,
		Username:   username,
		UID:        uid,
		Profiles:   []Profile{},
	}
}
