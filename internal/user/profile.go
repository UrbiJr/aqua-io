package user

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Profile struct {
	Title        string `json:"title"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
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

// ReadProfiles reads profiles.json and returns read data as []Profile
func ReadProfiles() ([]Profile, error) {

	var profiles []Profile
	jsonFile, err := os.Open("profiles.json")

	if err != nil {
		// file does not exist
		profiles = []Profile{}
		file, _ := json.MarshalIndent(profiles, "", " ")
		_ = ioutil.WriteFile("profiles.json", file, 0644)
		return profiles, nil
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &profiles)
	if err != nil {
		// error reading file
		return []Profile{}, err
	}

	return profiles, nil

}

// WriteProfiles writes profiles to profiles.json
func WriteProfiles(profiles []Profile) {
	file, _ := json.MarshalIndent(profiles, "", " ")
	_ = ioutil.WriteFile("profiles.json", file, 0644)
}
