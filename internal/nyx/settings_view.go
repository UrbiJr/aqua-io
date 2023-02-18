package nyx

import (
	"github.com/rivo/tview"
)

// SettingsView contains information about the "Settings" view
type SettingsView struct {
	Title string
	View  *tview.Flex
}

// NewSettingsView returns a view for the profile creation
func (app *Config) NewSettingsView() *SettingsView {
	return nil
}

// EditSiteSettingsForm draws a form to edit a site settings
func (app *Config) EditSiteSettingsForm(site string) {

}

// EditSettingsForm draws a form to edit settings in the current view
func (app *Config) EditSettingsForm() {

}
