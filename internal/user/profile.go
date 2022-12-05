package user

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Profile struct {
	Title        string `json:"title"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	Postcode     string `json:"post_code"`
	State        string `json:"state"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	CardNumber   string `json:"card_number"`
	CardMonth    string `json:"card_month"`
	CardYear     string `json:"card_year"`
	CardCvv      string `json:"card_cvv"`
}
