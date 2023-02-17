package nyx

import (
	"github.com/rivo/tview"
)

// NewProfileView contains information about the "New Profile" view
type NewProfileView struct {
	Title string
	Form  *tview.Form
	View  *tview.Flex
}

// NewNewProfileView returns a view for the profile creation
func (nyx *Config) NewNewProfileView() *NewProfileView {

	return nil
}

func (nyx *Config) AddProfileForm() {

	//var profile user.Profile

}
