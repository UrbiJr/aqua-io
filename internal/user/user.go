package user

// User is a type associated with users of this software.
type User struct {
	ID                 int64   `json:"id"`
	Email              string  `json:"-"`
	DiscordID          string  `json:"-"`
	Username           string  `json:"-"`
	ProfilePictureURL  string  `json:"-"`
	ProfilePicturePath string  `json:"profile_picture_path"`
	LicenseKey         string  `json:"license_key"`
	ExpiresAt          float64 `json:"-"`
	PersistentLogin    bool    `json:"persistent_login"`
	Theme              string  `json:"theme"`
	*ProfileManager
}

// NewUser creates an instance of the User class
func NewUser(email, discordID, username, profilePicture, licenseKey string, expiresAt float64, persistentLogin bool) *User {
	return &User{
		Email:              email,
		DiscordID:          discordID,
		Username:           username,
		ProfilePicturePath: "",
		ProfilePictureURL:  profilePicture,
		LicenseKey:         licenseKey,
		ExpiresAt:          expiresAt,
		PersistentLogin:    persistentLogin,
		ProfileManager:     &ProfileManager{},
	}
}
