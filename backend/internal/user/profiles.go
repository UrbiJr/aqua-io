package user

// Profile contains information specific to a configuration for a trader
type Profile struct {
	ID                int64  `json:"id"`
	Title             string `json:"title"`
	BillingAddressID  int64  `json:"billing_address_id"`
	ShippingAddressID int64  `json:"shipping_address_id"`
	CardNumber        string `json:"card_number"`
	CardMonth         string `json:"card_month"`
	CardYear          string `json:"card_year"`
	CardCvv           string `json:"card_cvv"`
	TestMode          bool   `json:"test_mode"`
}