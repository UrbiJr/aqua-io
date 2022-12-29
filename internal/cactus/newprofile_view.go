package cactus

import (
	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewProfileView contains information about the "New Profile" view
type NewProfileView struct {
	Title string
	Form  *tview.Form
	View  *tview.Flex
}

// NewNewProfileView returns a view for the profile creation
func (cactus *Cactus) NewNewProfileView() *NewProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("New Profile"), 0, 1, false).
		AddItem(form, 0, 4, true).
		SetBorder(true)

	return &NewProfileView{Title: "New Profile", View: flex, Form: form}
}

func (cactus *Cactus) AddProfileForm() {

	var profile user.Profile

	cactus.NewProfileView.Form.AddInputField("Profile Title", "", 20, nil, func(title string) {
		profile.Title = title
	})

	cactus.NewProfileView.Form.AddInputField("First Name", "", 20, nil, func(firstName string) {
		profile.FirstName = firstName
	})

	cactus.NewProfileView.Form.AddInputField("Last Name", "", 20, nil, func(lastName string) {
		profile.LastName = lastName
	})

	cactus.NewProfileView.Form.AddInputField("Address Line 1", "", 30, nil, func(addressLine1 string) {
		profile.AddressLine1 = addressLine1
	})

	cactus.NewProfileView.Form.AddInputField("Address Line 2", "", 30, nil, func(addressLine2 string) {
		profile.AddressLine2 = addressLine2
	})

	cactus.NewProfileView.Form.AddInputField("City", "", 20, nil, func(city string) {
		profile.City = city
	})

	cactus.NewProfileView.Form.AddInputField("Postcode", "", 20, nil, func(postcode string) {
		profile.Postcode = postcode
	})

	cactus.NewProfileView.Form.AddInputField("Email", "", 30, nil, func(email string) {
		profile.Email = email
	})

	cactus.NewProfileView.Form.AddInputField("Phone", "", 30, nil, func(phone string) {
		profile.Phone = phone
	})

	cactus.NewProfileView.Form.AddDropDown("Country", utils.CountryNames, 0, func(country string, index int) {
		countryCode, err := utils.GetCountryCode(country)
		if err == nil {
			profile.CountryCode = countryCode
			cactus.RefreshStateSelectionView(countryCode, &profile)
			cactus.pages.SwitchToPage(cactus.StateSelectionView.Title)
		}
	})

	cactus.NewProfileView.Form.AddInputField("Card number", "", 20, nil, func(cardNumber string) {
		profile.CardNumber = cardNumber
	})

	cactus.NewProfileView.Form.AddDropDown("Card month", utils.CardMonths(), 0, func(month string, index int) {
		profile.CardMonth = month
	})

	cactus.NewProfileView.Form.AddDropDown("Card year", utils.CardYears(), 0, func(year string, index int) {
		profile.CardYear = year
	})

	cactus.NewProfileView.Form.AddInputField("CVV", "", 5, nil, func(cvv string) {
		profile.CardCvv = cvv
	})

	cactus.NewProfileView.Form.AddButton("Save", func() {
		err := cactus.SaveProfile(profile)
		if err != nil {
			cactus.ShowError(cactus.NewProfileView.Title, err)
		} else {
			cactus.RefreshProfileView()
			cactus.pages.SwitchToPage(cactus.ProfilesView.Title)
		}
	})

	cactus.NewProfileView.Form.AddButton("Cancel", func() {
		cactus.pages.SwitchToPage(cactus.ProfilesView.Title)
	})

}
