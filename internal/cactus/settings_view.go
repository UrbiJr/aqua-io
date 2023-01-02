package cactus

import (
	"fmt"
	"strconv"

	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SettingsView contains information about the "Settings" view
type SettingsView struct {
	Title string
	View  *tview.Flex
}

// NewSettingsView returns a view for the profile creation
func (cactus *Cactus) NewSettingsView() *SettingsView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		SetBorder(true)

	return &SettingsView{Title: "Settings", View: flex}
}

// EditSiteSettingsForm draws a form to edit a site settings
func (cactus *Cactus) EditSiteSettingsForm(site string) {

	header := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText(fmt.Sprintf("⚙️ %s settings", site))
	form := tview.NewForm()

	var siteSettings user.SiteSettings

	form.AddInputField("Captcha API", cactus.User.Settings.SitesSettings[site].Captcha, 20, nil, func(api string) {
		siteSettings.Captcha = api
	})

	form.AddInputField("Delay", strconv.Itoa(cactus.User.Settings.SitesSettings[site].Delay), 20, nil, func(delay string) {
		delayInt, err := strconv.Atoi(delay)
		if err == nil {
			siteSettings.Delay = delayInt
		}
	})

	form.AddButton("Save", func() {
		cactus.User.Settings.SitesSettings[site] = siteSettings
		cactus.EditSettingsForm()
	})

	form.AddButton("Cancel", func() {
		cactus.EditSettingsForm()
	})

	cactus.SettingsView.View.Clear()
	cactus.SettingsView.View.AddItem(header, 0, 1, false)
	cactus.SettingsView.View.AddItem(form, 0, 4, true)

}

// EditSettingsForm draws a form to edit settings in the current view
func (cactus *Cactus) EditSettingsForm() {

	header := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("⚙️ Settings")
	form := tview.NewForm()
	footer := tview.NewTextView().SetTextColor(tcell.ColorDarkGreen)

	form.AddInputField("Webhook", cactus.User.Settings.Webhook, 20, nil, func(webhook string) {
		cactus.User.Settings.Webhook = webhook
	})

	form.AddInputField("Delay", strconv.Itoa(cactus.User.Settings.Delay), 20, nil, func(delay string) {
		delayInt, err := strconv.Atoi(delay)
		if err == nil {
			cactus.User.Settings.Delay = delayInt
		}
	})

	form.AddInputField("2captcha API key", cactus.User.Settings.APIKeys["2captcha"], 30, nil, func(apikey string) {
		cactus.User.Settings.APIKeys["2captcha"] = apikey
	})

	form.AddInputField("AntiCaptcha API key", cactus.User.Settings.APIKeys["anticaptcha"], 30, nil, func(apikey string) {
		cactus.User.Settings.APIKeys["anticaptcha"] = apikey
	})

	form.AddInputField("Capmonster API key", cactus.User.Settings.APIKeys["capmonster"], 30, nil, func(apikey string) {
		cactus.User.Settings.APIKeys["capmonster"] = apikey
	})

	var sitelist []string
	for _, site := range cactus.SiteList {
		sitelist = append(sitelist, site.Name.String())
	}
	form.AddDropDown("Edit site settings", sitelist, 0, func(site string, index int) {
		cactus.EditSiteSettingsForm(site)
	})

	form.AddButton("Save", func() {
		user.WriteSettings(cactus.User.Settings)
		footer.SetText("💾 Successfully saved!")
	})

	form.AddButton("Go Back", func() {
		cactus.pages.SwitchToPage(cactus.MainMenuView.Title)
	})

	cactus.SettingsView.View.Clear()
	cactus.SettingsView.View.AddItem(header, 0, 1, false)
	cactus.SettingsView.View.AddItem(form, 0, 4, true)
	cactus.SettingsView.View.AddItem(footer, 0, 1, false)
}
