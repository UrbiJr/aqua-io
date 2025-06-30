package user

// Profile contains information specific to a configuration for a trader
type Address struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"` // used to retrieve exchange settings
	Phone        string `json:"phone"`     // used to retrieve exchange settings
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	ZipCode      string `json:"zip_code"`
	Province     string `json:"province"`
	CountryCode  string `json:"country_code"`
}
