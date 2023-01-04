package nyx

import (
	"github.com/UrbiJr/go-cactus/internal/user"
	"github.com/UrbiJr/go-cactus/internal/utils"
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
func (nyx *Nyx) NewNewProfileView() *NewProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("New Profile"), 0, 1, false).
		AddItem(form, 0, 8, true).
		SetBorder(true)

	return &NewProfileView{Title: "New Profile", View: flex, Form: form}
}

func (nyx *Nyx) AddProfileForm() {

	var profile user.Profile

	nyx.NewProfileView.Form.AddInputField("Profile Title", "", 20, nil, func(title string) {
		profile.Title = title
	})

	nyx.NewProfileView.Form.AddInputField("First Name", "", 20, nil, func(firstName string) {
		profile.FirstName = firstName
	})

	nyx.NewProfileView.Form.AddInputField("Last Name", "", 20, nil, func(lastName string) {
		profile.LastName = lastName
	})

	nyx.NewProfileView.Form.AddInputField("Address Line 1", "", 30, nil, func(addressLine1 string) {
		profile.AddressLine1 = addressLine1
	})

	nyx.NewProfileView.Form.AddInputField("Address Line 2", "", 30, nil, func(addressLine2 string) {
		profile.AddressLine2 = addressLine2
	})

	nyx.NewProfileView.Form.AddInputField("City", "", 20, nil, func(city string) {
		profile.City = city
	})

	nyx.NewProfileView.Form.AddInputField("Postcode", "", 20, nil, func(postcode string) {
		profile.Postcode = postcode
	})

	nyx.NewProfileView.Form.AddInputField("Email", "", 30, nil, func(email string) {
		profile.Email = email
	})

	nyx.NewProfileView.Form.AddInputField("Phone", "", 30, nil, func(phone string) {
		profile.Phone = phone
	})

	nyx.NewProfileView.Form.AddDropDown("Country", utils.CountryNames, 0, func(country string, index int) {
		countryCode, err := utils.GetCountryCode(country)
		if err == nil {
			profile.CountryCode = countryCode
			nyx.RefreshStateSelectionView(nyx.NewProfileView.Title, countryCode, &profile)
			nyx.pages.SwitchToPage(nyx.StateSelectionView.Title)
		}
	})

	nyx.NewProfileView.Form.AddInputField("Card number", "", 20, nil, func(cardNumber string) {
		profile.CardNumber = cardNumber
	})

	nyx.NewProfileView.Form.AddDropDown("Card month", utils.CardMonths(), 0, func(month string, index int) {
		profile.CardMonth = month
	})

	nyx.NewProfileView.Form.AddDropDown("Card year", utils.CardYears(), 0, func(year string, index int) {
		profile.CardYear = year
	})

	nyx.NewProfileView.Form.AddInputField("CVV", "", 5, nil, func(cvv string) {
		profile.CardCvv = cvv
	})

	nyx.NewProfileView.Form.AddButton("Save", func() {
		err := nyx.AddProfile(profile)
		if err != nil {
			nyx.ShowError(nyx.NewProfileView.Title, err)
		} else {
			nyx.RefreshProfileView()
			nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
		}
	})

	nyx.NewProfileView.Form.AddButton("Cancel", func() {
		nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
	})

}
