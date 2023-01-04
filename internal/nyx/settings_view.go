package nyx

import (
	"fmt"
	"strconv"

	"github.com/UrbiJr/go-cactus/internal/user"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SettingsView contains information about the "Settings" view
type SettingsView struct {
	Title string
	View  *tview.Flex
}

// NewSettingsView returns a view for the profile creation
func (nyx *Nyx) NewSettingsView() *SettingsView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		SetBorder(true)

	return &SettingsView{Title: "Settings", View: flex}
}

// EditSiteSettingsForm draws a form to edit a site settings
func (nyx *Nyx) EditSiteSettingsForm(site string) {

	header := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText(fmt.Sprintf("⚙️ %s settings", site))
	form := tview.NewForm()

	var siteSettings user.SiteSettings

	form.AddInputField("Captcha API", nyx.User.Settings.SitesSettings[site].Captcha, 20, nil, func(api string) {
		siteSettings.Captcha = api
	})

	form.AddInputField("Delay", strconv.Itoa(nyx.User.Settings.SitesSettings[site].Delay), 20, nil, func(delay string) {
		delayInt, err := strconv.Atoi(delay)
		if err == nil {
			siteSettings.Delay = delayInt
		}
	})

	form.AddButton("Save", func() {
		nyx.User.Settings.SitesSettings[site] = siteSettings
		nyx.EditSettingsForm()
	})

	form.AddButton("Cancel", func() {
		nyx.EditSettingsForm()
	})

	nyx.SettingsView.View.Clear()
	nyx.SettingsView.View.AddItem(header, 0, 1, false)
	nyx.SettingsView.View.AddItem(form, 0, 4, true)

}

// EditSettingsForm draws a form to edit settings in the current view
func (nyx *Nyx) EditSettingsForm() {

	header := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("⚙️ Settings")
	form := tview.NewForm()
	footer := tview.NewTextView().SetTextColor(tcell.ColorDarkGreen)

	form.AddInputField("Webhook", nyx.User.Settings.Webhook, 20, nil, func(webhook string) {
		nyx.User.Settings.Webhook = webhook
	})

	form.AddInputField("Delay", strconv.Itoa(nyx.User.Settings.Delay), 20, nil, func(delay string) {
		delayInt, err := strconv.Atoi(delay)
		if err == nil {
			nyx.User.Settings.Delay = delayInt
		}
	})

	form.AddInputField("2captcha API key", nyx.User.Settings.APIKeys["2captcha"], 30, nil, func(apikey string) {
		nyx.User.Settings.APIKeys["2captcha"] = apikey
	})

	form.AddInputField("AntiCaptcha API key", nyx.User.Settings.APIKeys["anticaptcha"], 30, nil, func(apikey string) {
		nyx.User.Settings.APIKeys["anticaptcha"] = apikey
	})

	form.AddInputField("Capmonster API key", nyx.User.Settings.APIKeys["capmonster"], 30, nil, func(apikey string) {
		nyx.User.Settings.APIKeys["capmonster"] = apikey
	})

	var sitelist []string
	for _, site := range nyx.SiteList {
		sitelist = append(sitelist, site.Name.String())
	}
	form.AddDropDown("Edit site settings", sitelist, 0, func(site string, index int) {
		nyx.EditSiteSettingsForm(site)
	})

	form.AddButton("Save", func() {
		user.WriteSettings(nyx.User.Settings)
		footer.SetText("💾 Successfully saved!")
	})

	form.AddButton("Go Back", func() {
		nyx.pages.SwitchToPage(nyx.MainMenuView.Title)
	})

	nyx.SettingsView.View.Clear()
	nyx.SettingsView.View.AddItem(header, 0, 1, false)
	nyx.SettingsView.View.AddItem(form, 0, 4, true)
	nyx.SettingsView.View.AddItem(footer, 0, 1, false)
}
