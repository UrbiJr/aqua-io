package cactus

import (
	"errors"
	"log"

	"github.com/cactus-aio/go-cactus/internal/sites"
	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
)

// Cactus is the container of the main app, it contains the main attributes
type Cactus struct {
	*UI
	SiteList []*sites.SupportedSite
	*log.Logger
	User *user.User
}

// NewCactus returns a new instance of the app
func NewCactus() *Cactus {
	sitelist := []*sites.SupportedSite{
		{Name: sites.Nike, Category: sites.SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
		{Name: sites.Schuh, Category: sites.SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
	}

	cactus := &Cactus{
		UI:       NewUI(),
		SiteList: sitelist,
	}

	return cactus
}

// Quit exits the app gracefully
func (cactus *Cactus) Quit() {
	utils.QuitLogger()
	cactus.UI.tui.Stop()
}

// SaveProfile appends profile to user profiles list and writes the updated list to file
func (cactus *Cactus) SaveProfile(profile user.Profile) error {
	for _, p := range cactus.User.Profiles {
		if p.Title == profile.Title {
			return errors.New("a profile with this title is already existent")
		}
	}
	cactus.User.Profiles = append(cactus.User.Profiles, profile)
	user.WriteProfiles(cactus.User.Profiles)
	return nil
}

// UpdateProfile updates an existing profile and writes the updated profile list to file
func (cactus *Cactus) UpdateProfile(profile user.Profile) error {
	for _, p := range cactus.User.Profiles {
		if p.Title == profile.Title {
			p = profile
		}
	}
	user.WriteProfiles(cactus.User.Profiles)
	return nil
}

// DeleteProfile removes profile from user profiles list and writes the updated list to file
func (cactus *Cactus) DeleteProfile(profileTitle string) error {
	for i, p := range cactus.User.Profiles {
		if p.Title == profileTitle {
			// remove from slice and preserve order
			cactus.User.Profiles = append(cactus.User.Profiles[:i], cactus.User.Profiles[i+1:]...)
			user.WriteProfiles(cactus.User.Profiles)
			return nil
		}
	}
	return errors.New("cannot find a profile with this title")
}
