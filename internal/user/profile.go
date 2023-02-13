package user

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Profile struct {
	Title        string
	Email        string
	FirstName    string
	LastName     string
	FullName     string
	AddressLine1 string
	AddressLine2 string
	City         string
	Postcode     string
	State        string
	CountryCode  string
	Phone        string
	CardNumber   string
	CardMonth    string
	CardYear     string
	CardCvv      string
}

// ReadProfiles reads profiles from DB and returns read data as []Profile
func ReadProfiles() ([]Profile, error) {

	var profiles []Profile

	return profiles, nil

}

// WriteProfiles writes profiles to the DB
func WriteProfiles(profiles []Profile) {

}
