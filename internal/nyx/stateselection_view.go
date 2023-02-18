package nyx

import (
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/rivo/tview"
)

// StateSelectionView contains information about the "State selection" view
type StateSelectionView struct {
	Title string
	Form  *tview.Form
	View  *tview.Flex
}

// RefreshStateSelectionView fills the state selection form with the states of the specified country
func (app *Config) RefreshStateSelectionView(currentPage string, countryCode string, profile *user.Profile) {

}

// NewStateSelectionView returns a view for the state selection by country
func (app *Config) NewStateSelectionView() *StateSelectionView {

	return nil
}
