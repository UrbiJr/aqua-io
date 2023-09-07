package user

import "github.com/UrbiJr/aqua-io/backend/internal/utils/constants"

// Profile contains information specific to a configuration for a trader
type Profile struct {
	ID              int64              `json:"id"`
	Title           string             `json:"title"`
	Exchange        constants.Exchange `json:"exchange"`
	AccountName     string             `json:"account_name"`
	PublicAPI       string             `json:"public_api"`
	SecretAPI       string             `json:"secret_api"`
	Passphrase      string             `json:"passphrase"`
	StopIfFallUnder float64            `json:"stop_if_fall_under"`
	TestMode        bool               `json:"test_mode"`
}
