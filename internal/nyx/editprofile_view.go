package nyx

import (
	"fmt"

	"github.com/UrbiJr/go-cactus/internal/user"
	"github.com/UrbiJr/go-cactus/internal/utils"
	"github.com/gdamore/tcell/v2"
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
func (nyx *Nyx) NewEditProfileView() *EditProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()
	header := tview.NewTextView().SetTextColor(tcell.ColorGreen)

	flex.SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(form, 0, 8, true).
		SetBorder(true)

	return &EditProfileView{Title: "Edit Profile", Header: header, View: flex, Form: form}
}

// EditProfileForm draws a form to edit a profile in the current view
func (nyx *Nyx) EditProfileForm(profile user.Profile) {

	nyx.EditProfileView.Form.Clear(true)
	nyx.EditProfileView.Header.SetText("👤 " + fmt.Sprintf("Editing profile \"%s\"", profile.Title))

	nyx.EditProfileView.Form.AddInputField("First Name", profile.FirstName, 20, nil, func(firstName string) {
		profile.FirstName = firstName
	})

	nyx.EditProfileView.Form.AddInputField("Last Name", profile.LastName, 20, nil, func(lastName string) {
		profile.LastName = lastName
	})

	nyx.EditProfileView.Form.AddInputField("Address Line 1", profile.AddressLine1, 30, nil, func(addressLine1 string) {
		profile.AddressLine1 = addressLine1
	})

	nyx.EditProfileView.Form.AddInputField("Address Line 2", profile.AddressLine2, 30, nil, func(addressLine2 string) {
		profile.AddressLine2 = addressLine2
	})

	nyx.EditProfileView.Form.AddInputField("City", profile.City, 20, nil, func(city string) {
		profile.City = city
	})

	nyx.EditProfileView.Form.AddInputField("Postcode", profile.Postcode, 20, nil, func(postcode string) {
		profile.Postcode = postcode
	})

	nyx.EditProfileView.Form.AddInputField("Email", profile.Email, 30, nil, func(email string) {
		profile.Email = email
	})

	nyx.EditProfileView.Form.AddInputField("Phone", profile.Phone, 30, nil, func(phone string) {
		profile.Phone = phone
	})

	idx := 0
	country, err := utils.GetCountryName(profile.CountryCode)
	if err == nil {
		idx = utils.StringSliceIndex(country, utils.CountryNames)
	}
	nyx.EditProfileView.Form.AddDropDown("Country", utils.CountryNames, idx, func(country string, index int) {
		countryCode, err := utils.GetCountryCode(country)
		if err == nil {
			profile.CountryCode = countryCode
			nyx.RefreshStateSelectionView(nyx.EditProfileView.Title, countryCode, &profile)
			nyx.pages.SwitchToPage(nyx.StateSelectionView.Title)
		}
	})

	nyx.EditProfileView.Form.AddInputField("Card number", profile.CardNumber, 20, nil, func(cardNumber string) {
		profile.CardNumber = cardNumber
	})

	nyx.EditProfileView.Form.AddDropDown("Card month", utils.CardMonths(), utils.StringSliceIndex(profile.CardMonth, utils.CardMonths()), func(month string, index int) {
		profile.CardMonth = month
	})

	nyx.EditProfileView.Form.AddDropDown("Card year", utils.CardYears(), utils.StringSliceIndex(profile.CardYear, utils.CardYears()), func(year string, index int) {
		profile.CardYear = year
	})

	nyx.EditProfileView.Form.AddInputField("CVV", profile.CardCvv, 5, nil, func(cvv string) {
		profile.CardCvv = cvv
	})

	nyx.EditProfileView.Form.AddButton("Save", func() {
		err := nyx.UpdateProfile(profile)
		if err != nil {
			nyx.ShowError(nyx.EditProfileView.Title, err)
		} else {
			nyx.RefreshProfileView()
			nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
		}
	})

	nyx.EditProfileView.Form.AddButton("Cancel", func() {
		nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
	})

}

// RenameProfileForm draws a form to rename a profile in the current view
func (nyx *Nyx) RenameProfileForm(profile user.Profile) {

	oldTitle := profile.Title
	nyx.EditProfileView.Form.Clear(true)
	nyx.EditProfileView.Header.SetText("👤 Rename profile")

	nyx.EditProfileView.Form.AddInputField("New Title", profile.Title, 20, nil, func(title string) {
		profile.Title = title
	})

	nyx.EditProfileView.Form.AddButton("Save", func() {
		err := nyx.UpdateProfileTitle(oldTitle, profile)
		if err != nil {
			nyx.ShowError(nyx.EditProfileView.Title, err)
		} else {
			nyx.RefreshProfileView()
			nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
		}
	})

	nyx.EditProfileView.Form.AddButton("Cancel", func() {
		nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
	})

}
