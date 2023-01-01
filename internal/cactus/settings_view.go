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
	Title  string
	Header *tview.TextView
	Footer *tview.TextView
	Form   *tview.Form
	View   *tview.Flex
}

// NewSettingsView returns a view for the profile creation
func (cactus *Cactus) NewSettingsView() *SettingsView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()
	header := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Edit settings")
	footer := tview.NewTextView().SetTextColor(tcell.ColorDarkGreen)

	flex.SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(form, 0, 8, true).
		AddItem(footer, 0, 1, false).
		SetBorder(true)

	return &SettingsView{Title: "Settings", Header: header, Footer: footer, View: flex, Form: form}
}

// EditSiteSettingsForm draws a form to edit a site settings
func (cactus *Cactus) EditSiteSettingsForm(site string) {

	cactus.SettingsView.Form.Clear(true)
	cactus.SettingsView.Header.SetText(fmt.Sprintf("Editing %s settings", site))
	var siteSettings user.SiteSettings

	cactus.SettingsView.Form.AddInputField("Captcha API", cactus.User.Settings.SitesSettings[site].Captcha, 20, nil, func(api string) {
		siteSettings.Captcha = api
	})

	cactus.SettingsView.Form.AddInputField("Delay", strconv.Itoa(cactus.User.Settings.SitesSettings[site].Delay), 20, nil, func(delay string) {
		delayInt, err := strconv.Atoi(delay)
		if err == nil {
			siteSettings.Delay = delayInt
		}
	})

	cactus.SettingsView.Form.AddButton("Save", func() {
		cactus.User.Settings.SitesSettings[site] = siteSettings
		cactus.EditSettingsForm()
	})

	cactus.SettingsView.Form.AddButton("Cancel", func() {
		cactus.EditSettingsForm()
	})

}

// EditSettingsForm draws a form to edit settings in the current view
func (cactus *Cactus) EditSettingsForm() {

	cactus.SettingsView.Form.Clear(true)

	cactus.SettingsView.Form.AddInputField("Webhook", cactus.User.Settings.Webhook, 20, nil, func(webhook string) {
		cactus.User.Settings.Webhook = webhook
	})

	cactus.SettingsView.Form.AddInputField("Delay", strconv.Itoa(cactus.User.Settings.Delay), 20, nil, func(delay string) {
		delayInt, err := strconv.Atoi(delay)
		if err == nil {
			cactus.User.Settings.Delay = delayInt
		}
	})

	cactus.SettingsView.Form.AddInputField("2captcha API key", cactus.User.Settings.APIKeys["2captcha"], 30, nil, func(apikey string) {
		cactus.User.Settings.APIKeys["2captcha"] = apikey
	})

	cactus.SettingsView.Form.AddInputField("AntiCaptcha API key", cactus.User.Settings.APIKeys["anticaptcha"], 30, nil, func(apikey string) {
		cactus.User.Settings.APIKeys["anticaptcha"] = apikey
	})

	cactus.SettingsView.Form.AddInputField("Capmonster API key", cactus.User.Settings.APIKeys["capmonster"], 30, nil, func(apikey string) {
		cactus.User.Settings.APIKeys["capmonster"] = apikey
	})

	var sitelist []string
	for _, site := range cactus.SiteList {
		sitelist = append(sitelist, site.Name.String())
	}
	cactus.SettingsView.Form.AddDropDown("Edit site settings", sitelist, 0, func(site string, index int) {
		//cactus.EditSiteSettingsForm(site)
	})

	cactus.SettingsView.Form.AddButton("Save", func() {
		user.WriteSettings(cactus.User.Settings)
		cactus.SettingsView.Footer.SetText("Settings successfully saved")
	})

	cactus.SettingsView.Form.AddButton("Go Back", func() {
		cactus.pages.SwitchToPage(cactus.MainMenuView.Title)
	})

}
