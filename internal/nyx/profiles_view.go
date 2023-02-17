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
func (nyx *Config) NewProfilesView() *ProfilesView {

	return nil
}

// RefreshProfilesView refresh the list of profiles
func (nyx *Config) RefreshProfilesView() {
	/*
		for i, profile := range nyx.User.Profiles {
			// table cell containing profile name

		}
	*/
}
