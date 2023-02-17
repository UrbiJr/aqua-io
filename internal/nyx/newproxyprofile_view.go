package nyx

import (
	"github.com/rivo/tview"
)

// NewProxyProfileView contains information about the "New Profile" view
type NewProxyProfileView struct {
	Title string
	Form  *tview.Form
	View  *tview.Flex
}

// NewNewProxyProfileView returns a view for the profile creation
func (nyx *Config) NewNewProxyProfileView() *NewProxyProfileView {

	return nil
}

func (nyx *Config) AddProxyProfileForm() {

}
