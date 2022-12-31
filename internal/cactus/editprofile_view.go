package cactus

import (
	"fmt"

	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
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
func (cactus *Cactus) NewEditProfileView() *EditProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()
	header := tview.NewTextView().SetTextColor(tcell.ColorGreen)

	flex.SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(form, 0, 8, true).
		SetBorder(true)

	return &EditProfileView{Title: "Edit Profile", Header: header, View: flex, Form: form}
}

func (cactus *Cactus) EditProfileForm(profile user.Profile) {

	cactus.EditProfileView.Header.SetText(fmt.Sprintf("Editing profile \"%s\"", profile.Title))
	cactus.EditProfileView.Form.Clear(true)

	cactus.EditProfileView.Form.AddInputField("First Name", profile.FirstName, 20, nil, func(firstName string) {
		profile.FirstName = firstName
	})

	cactus.EditProfileView.Form.AddInputField("Last Name", profile.LastName, 20, nil, func(lastName string) {
		profile.LastName = lastName
	})

	cactus.EditProfileView.Form.AddInputField("Address Line 1", profile.AddressLine1, 30, nil, func(addressLine1 string) {
		profile.AddressLine1 = addressLine1
	})

	cactus.EditProfileView.Form.AddInputField("Address Line 2", profile.AddressLine2, 30, nil, func(addressLine2 string) {
		profile.AddressLine2 = addressLine2
	})

	cactus.EditProfileView.Form.AddInputField("City", profile.City, 20, nil, func(city string) {
		profile.City = city
	})

	cactus.EditProfileView.Form.AddInputField("Postcode", profile.Postcode, 20, nil, func(postcode string) {
		profile.Postcode = postcode
	})

	cactus.EditProfileView.Form.AddInputField("Email", profile.Email, 30, nil, func(email string) {
		profile.Email = email
	})

	cactus.EditProfileView.Form.AddInputField("Phone", profile.Phone, 30, nil, func(phone string) {
		profile.Phone = phone
	})

	idx := 0
	country, err := utils.GetCountryName(profile.CountryCode)
	if err == nil {
		idx = utils.StringSliceIndex(country, utils.CountryNames)
	}
	cactus.EditProfileView.Form.AddDropDown("Country", utils.CountryNames, idx, func(country string, index int) {
		countryCode, err := utils.GetCountryCode(country)
		if err == nil {
			profile.CountryCode = countryCode
			cactus.RefreshStateSelectionView(cactus.EditProfileView.Title, countryCode, &profile)
			cactus.pages.SwitchToPage(cactus.StateSelectionView.Title)
		}
	})

	cactus.EditProfileView.Form.AddInputField("Card number", profile.CardNumber, 20, nil, func(cardNumber string) {
		profile.CardNumber = cardNumber
	})

	cactus.EditProfileView.Form.AddDropDown("Card month", utils.CardMonths(), utils.StringSliceIndex(profile.CardMonth, utils.CardMonths()), func(month string, index int) {
		profile.CardMonth = month
	})

	cactus.EditProfileView.Form.AddDropDown("Card year", utils.CardYears(), utils.StringSliceIndex(profile.CardYear, utils.CardYears()), func(year string, index int) {
		profile.CardYear = year
	})

	cactus.EditProfileView.Form.AddInputField("CVV", profile.CardCvv, 5, nil, func(cvv string) {
		profile.CardCvv = cvv
	})

	cactus.EditProfileView.Form.AddButton("Save", func() {
		err := cactus.UpdateProfile(profile)
		if err != nil {
			cactus.ShowError(cactus.EditProfileView.Title, err)
		} else {
			cactus.RefreshProfileView()
			cactus.pages.SwitchToPage(cactus.ProfilesView.Title)
		}
	})

	cactus.EditProfileView.Form.AddButton("Cancel", func() {
		cactus.pages.SwitchToPage(cactus.ProfilesView.Title)
	})

}
