package user

import "github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"

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

type ProfileManager struct {
	Profiles []Profile
}

func (pfm *ProfileManager) GetProfileByTitle(title string) *Profile {
	for _, p := range pfm.Profiles {
		if p.Title == title {
			return &p
		}
	}

	return nil
}

func (pfm *ProfileManager) UpdateProfile(ID int64, updated Profile) *Profile {
	for i, p := range pfm.Profiles {
		if p.ID == ID {
			pfm.Profiles[i].Title = updated.Title
			pfm.Profiles[i].PublicAPI = updated.PublicAPI
			pfm.Profiles[i].SecretAPI = updated.SecretAPI
			pfm.Profiles[i].TestMode = updated.TestMode
			return &pfm.Profiles[i]
		}
	}

	return nil
}

func (pfm *ProfileManager) DeleteProfile(ID int64) {
	for i, p := range pfm.Profiles {
		if p.ID == ID {
			pfm.Profiles = append(pfm.Profiles[:i], pfm.Profiles[i+1:]...)
			break
		}
	}
}

func (pfm *ProfileManager) GetProfileByID(ID int64) *Profile {
	for _, p := range pfm.Profiles {
		if p.ID == ID {
			return &p
		}
	}
	return nil
}

func (pfm *ProfileManager) GetAllTitles() []string {
	var titles []string
	for _, p := range pfm.Profiles {
		titles = append(titles, p.Title)
	}

	return titles
}
