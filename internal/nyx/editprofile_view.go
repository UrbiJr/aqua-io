package nyx

import (
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/rivo/tview"
)

// EditProfileView contains information about the "Edit Profile" view
type EditProfileView struct {
	Title  string
	Header *tview.TextView
	Form   *tview.Form
	View   *tview.Flex
}

// NewEditProfileView returns a view for the profile creation
func (nyx *Config) NewEditProfileView() *EditProfileView {

	return nil
}

// EditProfileForm draws a form to edit a profile in the current view
func (nyx *Config) EditProfileForm(profile user.Profile) {

}

// RenameProfileForm draws a form to rename a profile in the current view
func (nyx *Config) RenameProfileForm(profile user.Profile) {

	//oldTitle := profile.Title

}
