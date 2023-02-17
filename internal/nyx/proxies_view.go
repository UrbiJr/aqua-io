package nyx

import (
	"github.com/rivo/tview"
)

// ProxiesView contains information about the "Profile" view
type ProxiesView struct {
	Title string
	View  *tview.Flex
}

// NewTaskCreationView returns a view for the profiles management
func (nyx *Config) NewProxiesView() *ProxiesView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	return &ProxiesView{
		Title: "Proxies",
		View:  flex,
	}
}

// RefreshProxiesView refreshes this view with updated proxy profiles and the current selection, if any
func (nyx *Config) RefreshProxiesView(proxyProfileIdx int) {

	/*
		for i, proxyProfile := range nyx.User.ProxyProfiles {

		}
	*/

}
