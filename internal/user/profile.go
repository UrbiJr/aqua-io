package user

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Profile struct {
	ProfileName  string `csv:"PROFILE NAME"`
	Email        string `csv:"EMAIL"`
	FirstName    string `csv:"FIRST NAME"`
	LastName     string `csv:"LAST NAME"`
	FullName     string
	AddressLine1 string `csv:"ADDRESS LINE 1"`
	AddressLine2 string `csv:"ADDRESS LINE 2"`
	City         string `csv:"CITY"`
	Postcode     string `csv:"POST CODE"`
	State        string `csv:"STATE"`
	CountryCode  string `csv:"COUNTRY CODE"`
	Phone        string `csv:"PHONE"`
	CardNumber   string `csv:"CARD NUMBER"`
	CardMonth    string `csv:"CARD MONTH"`
	CardYear     string `csv:"CARD YEAR"`
	CardCvv      string `csv:"CARD CVV"`
}
