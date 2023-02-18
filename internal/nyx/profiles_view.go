package nyx

import (
	"github.com/rivo/tview"
)

// ProfilesView contains information about the "Profile" view
type ProfilesView struct {
	Title string
	*tview.Table
	View *tview.Flex
}

// NewProfilesView returns a view for the profiles management
func (app *Config) NewProfilesView() *ProfilesView {

	return nil
}

// RefreshProfilesView refresh the list of profiles
func (app *Config) RefreshProfilesView() {
	/*
		for i, profile := range app.User.Profiles {
			// table cell containing profile name

		}
	*/
}
