package nyx

import (
	"errors"
	"log"

	"github.com/UrbiJr/go-cactus/internal/sites"
	"github.com/UrbiJr/go-cactus/internal/user"
	"github.com/UrbiJr/go-cactus/internal/utils"
)

// Nyx is the container of the main app, it contains the main attributes
type Nyx struct {
	*UI
	SiteList []*sites.SupportedSite
	*log.Logger
	User *user.User
}

// NewNyx returns a new instance of the app
func NewNyx() *Nyx {
	sitelist := []*sites.SupportedSite{
		{Name: sites.Kickz, Category: sites.SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
	}

	nyx := &Nyx{
		UI:       NewUI(),
		SiteList: sitelist,
	}

	return nyx
}

// Quit exits the app gracefully
func (nyx *Nyx) Quit() {
	utils.QuitLogger()
	nyx.UI.tui.Stop()
}

// AddProfile appends profile to user profiles list and writes the updated list to file
func (nyx *Nyx) AddProfile(profile user.Profile) error {
	for _, p := range nyx.User.Profiles {
		if p.Title == profile.Title {
			return errors.New("a profile with this title is already existent")
		}
	}
	nyx.User.Profiles = append(nyx.User.Profiles, profile)
	user.WriteProfiles(nyx.User.Profiles)
	return nil
}

// UpdateProfile updates an existing profile and writes the updated profile list to file
func (nyx *Nyx) UpdateProfile(profile user.Profile) error {
	idx := -1
	for i, p := range nyx.User.Profiles {
		if p.Title == profile.Title {
			idx = i
		}
	}
	if idx == -1 {
		return errors.New("match not found")
	}
	nyx.User.Profiles[idx] = profile
	user.WriteProfiles(nyx.User.Profiles)
	return nil
}

// UpdateProfileTitle updates an existing profile title and writes the updated profile list to file
func (nyx *Nyx) UpdateProfileTitle(oldtitle string, profile user.Profile) error {
	idx := -1
	for i, p := range nyx.User.Profiles {
		if p.Title == oldtitle {
			idx = i
		}
	}
	if idx == -1 {
		return errors.New("match not found")
	}
	nyx.User.Profiles[idx].Title = profile.Title
	user.WriteProfiles(nyx.User.Profiles)
	return nil
}

// DeleteProfile removes profile from user profiles list and writes the updated list to file
func (nyx *Nyx) DeleteProfile(profileTitle string) error {
	for i, p := range nyx.User.Profiles {
		if p.Title == profileTitle {
			// remove from slice and preserve order
			nyx.User.Profiles = append(nyx.User.Profiles[:i], nyx.User.Profiles[i+1:]...)
			user.WriteProfiles(nyx.User.Profiles)
			return nil
		}
	}
	return errors.New("cannot find a profile with this title")
}
