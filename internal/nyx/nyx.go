package nyx

import (
	"errors"
	"net/http"

	"fyne.io/fyne/v2"
	"github.com/UrbiJr/nyx/internal/client"
	"github.com/UrbiJr/nyx/internal/sites"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

// Config is the container of the main app, it contains the main attributes
type Config struct {
	App                     fyne.App
	MainWindow              fyne.Window
	CheckoutsChartContainer *fyne.Container
	SiteList                []*sites.SupportedSite
	Logger                  *utils.AppLogger
	User                    *user.User
	HTTPClient              *http.Client
	TLSClient               *client.Client
}

// NewNyx returns a new instance of the app
func NewNyx() *Config {
	sitelist := []*sites.SupportedSite{
		{Name: sites.Kickz, Category: sites.SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
	}

	nyx := &Config{
		SiteList: sitelist,
	}

	return nyx
}

// Quit exits the app gracefully
func (app *Config) Quit() {
	app.Logger.QuitLogger()
}

// AddProxyProfile appends profile to user profiles list and writes the updated list to file
func (app *Config) AddProxyProfile(profile user.ProxyProfile) error {
	profile.Id = utils.RandString(12, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	app.User.ProxyProfiles = append(app.User.ProxyProfiles, profile)
	user.WriteProxies(app.User.ProxyProfiles)
	return nil
}

// AddProfile appends profile to user profiles list and writes the updated list to file
func (app *Config) AddProfile(profile user.Profile) error {
	for _, p := range app.User.Profiles {
		if p.Title == profile.Title {
			return errors.New("a profile with this title is already existent")
		}
	}
	app.User.Profiles = append(app.User.Profiles, profile)
	user.WriteProfiles(app.User.Profiles)
	return nil
}

// UpdateProfile updates an existing profile and writes the updated profile list to file
func (app *Config) UpdateProfile(profile user.Profile) error {
	idx := -1
	for i, p := range app.User.Profiles {
		if p.Title == profile.Title {
			idx = i
		}
	}
	if idx == -1 {
		return errors.New("match not found")
	}
	app.User.Profiles[idx] = profile
	user.WriteProfiles(app.User.Profiles)
	return nil
}

// UpdateProfileTitle updates an existing profile title and writes the updated profile list to file
func (app *Config) UpdateProfileTitle(oldtitle string, profile user.Profile) error {
	idx := -1
	for i, p := range app.User.Profiles {
		if p.Title == oldtitle {
			idx = i
		}
	}
	if idx == -1 {
		return errors.New("match not found")
	}
	app.User.Profiles[idx].Title = profile.Title
	user.WriteProfiles(app.User.Profiles)
	return nil
}

// DeleteProfile removes profile from user profiles list and writes the updated list to file
func (app *Config) DeleteProfile(profileTitle string) error {
	for i, p := range app.User.Profiles {
		if p.Title == profileTitle {
			// remove from slice and preserve order
			app.User.Profiles = append(app.User.Profiles[:i], app.User.Profiles[i+1:]...)
			user.WriteProfiles(app.User.Profiles)
			return nil
		}
	}
	return errors.New("cannot find a profile with this title")
}
