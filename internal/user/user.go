package user

// User is a type associated with users of this software.
type User struct {
	Email          string
	Username       string
	UID            string
	Settings       *Settings
	ProfileManager *ProfileManager
	ProxyManager   *ProxyManager
}

// NewUser creates an instance of the User class
func NewUser(email string, uid string, username string) *User {
	return &User{
		Email:    email,
		Username: username,
		UID:      uid,
	}
}
