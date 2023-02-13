package nyx

import (
	"runtime"
	"strings"

	"github.com/UrbiJr/nyx/internal/client"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/gdamore/tcell/v2"
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

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("🌐 New Proxy Profile"), 0, 1, false).
		AddItem(form, 0, 8, true).
		SetBorder(true)

	return &NewProxyProfileView{Title: "New Proxy Profile", View: flex, Form: form}
}

func (nyx *Config) AddProxyProfileForm() {

	var proxyProfile user.ProxyProfile

	nyx.NewProxyProfileView.Form.AddInputField("Profile Name", "", 20, nil, func(title string) {
		proxyProfile.Name = title
	})

	nyx.NewProxyProfileView.Form.AddTextArea("Proxies", "", 0, 0, 20000, func(input string) {
		var splitChar string
		if runtime.GOOS == "windows" {
			splitChar = "\r\n"
		} else {
			splitChar = "\n"
		}
		lines := strings.Split(input, splitChar)
		for _, line := range lines {
			if client.ValidateProxyFormat(line) {
				proxyProfile.Proxies = append(proxyProfile.Proxies, line)
			}
		}

	})

	nyx.NewProxyProfileView.Form.AddButton("Save", func() {
		err := nyx.AddProxyProfile(proxyProfile)
		if err != nil {
			nyx.ShowError(nyx.NewProxyProfileView.Title, err)
		} else {
			nyx.RefreshProfilesView()
			nyx.pages.SwitchToPage(nyx.NewProxyProfileView.Title)
		}
	})

	nyx.NewProxyProfileView.Form.AddButton("Cancel", func() {
		nyx.pages.SwitchToPage(nyx.NewProxyProfileView.Title)
	})

}
