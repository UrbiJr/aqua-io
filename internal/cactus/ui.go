package cactus

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI contains information about the user interface
type UI struct {
	tui             *tview.Application
	pages           *tview.Pages
	MainMenuView    *tview.Flex
	SitelistView    *tview.Flex
	ProfilesView    *ProfileView
	SelectStateView *tview.Flex
	NewProfileView  *NewProfileView
}

// ProfileView contains information about the "Profile" view
type ProfileView struct {
	ProfilesList *tview.List
	View         *tview.Flex
}

// NewProfileView contains information about the "New Profile" view
type NewProfileView struct {
	Form *tview.Form
	View *tview.Flex
}

// MenuEntry contains information for each entry of a menu.
// name is the entry name to display
// description is a secondary text
// selected is the callback function which gets called when the entry is selected
type MenuEntry struct {
	name        string
	label       rune
	description string
	selected    func()
}

// Menu contains information about a Cactus-AIO ui menu
type Menu struct {
	name    string
	entries []MenuEntry
}

// indexToRune is a utility function to cast an int index to rune type in order to use it as menu shortcut
func indexToRune(index int) (rune, error) {
	s2 := strconv.Itoa(index)
	var b rune
	if len(s2) > 1 {
		return 0, errors.New("index cannot exceed 1 digit")
	} else {
		b = rune(s2[0])
	}
	return b, nil
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

	cactus.NewProfileView.Form.AddInputField("Address Line 1", "", 20, nil, func(addressLine1 string) {
		profile.AddressLine1 = addressLine1
	})

	cactus.NewProfileView.Form.AddInputField("Address Line 2", "", 20, nil, func(addressLine2 string) {
		profile.AddressLine2 = addressLine2
	})

	cactus.NewProfileView.Form.AddInputField("City", "", 20, nil, func(city string) {
		profile.City = city
	})

	cactus.NewProfileView.Form.AddInputField("Postcode", "", 20, nil, func(postcode string) {
		profile.Postcode = postcode
	})

	cactus.NewProfileView.Form.AddInputField("Email", "", 20, nil, func(email string) {
		profile.Email = email
	})

	cactus.NewProfileView.Form.AddInputField("Phone", "", 20, nil, func(phone string) {
		profile.Phone = phone
	})

	cactus.NewProfileView.Form.AddDropDown("Country", utils.CountryNames, 0, func(country string, index int) {
		countryCode, err := utils.GetCountryCode(country)
		if err == nil {
			profile.CountryCode = countryCode
			cactus.pages.AddPage("Select state", cactus.NewSelectStateView(countryCode, &profile), true, false)
			cactus.pages.SwitchToPage("Select state")
		}
	})

	cactus.NewProfileView.Form.AddButton("Save", func() {
		utils.Info(profile.State)
		cactus.User.Profiles = append(cactus.User.Profiles, profile)
		cactus.RefreshProfileView()
		cactus.pages.SwitchToPage("Profiles")
	})

	cactus.NewProfileView.Form.AddButton("Cancel", func() {
		cactus.pages.SwitchToPage("Profiles")
	})

}

// NewSelectStateView returns a view for the state selection by country
func (cactus *Cactus) NewSelectStateView(countryCode string, profile *user.Profile) *tview.Flex {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	states, err := utils.GetStates(countryCode)
	if err != nil {
		cactus.pages.SwitchToPage("New Profile")
	}

	form.AddDropDown("States", states, 0, func(state string, index int) {
		profile.State = state
	})

	form.AddButton("Save", func() {
		cactus.pages.SwitchToPage("New Profile")
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Select state"), 0, 1, false).
		AddItem(form, 0, 4, true).
		SetBorder(true)

	return flex
}

// NewNewProfileView returns a view for the profile creation
func (cactus *Cactus) NewNewProfileView() *NewProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("New Profile"), 0, 1, false).
		AddItem(form, 0, 4, true).
		SetBorder(true)

	return &NewProfileView{View: flex, Form: form}
}

// RefreshProfileView refresh the list of profiles
func (cactus *Cactus) RefreshProfileView() {
	cactus.ProfilesView.ProfilesList.Clear()
	for _, entry := range cactus.User.Profiles {
		// set 0 as shortcut for no binding
		cactus.ProfilesView.ProfilesList.AddItem(entry.Title, fmt.Sprintf("edit %s profile", entry.Title), 0, nil)
	}
	cactus.ProfilesView.ProfilesList.AddItem("Go Back", "return to the previous page", 'b', cactus.UI.OnGoBackSelected)
}

// NewSitelistView returns a view for the profiles management
func (cactus *Cactus) NewProfilesView() *ProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	profilesList := tview.NewList()
	for _, entry := range cactus.User.Profiles {
		// set 0 as shortcut for no binding
		profilesList.AddItem(entry.Title, fmt.Sprintf("edit %s profile", entry.Title), 0, nil)
	}
	profilesList.AddItem("Go Back", "return to the previous page", 'b', cactus.UI.OnGoBackSelected)

	helpInfo1 := tview.NewTextView().SetText("Press (a) to add a new Profile").SetTextAlign(tview.AlignLeft)
	helpInfo2 := tview.NewTextView().SetText("Press (d) to delete, (r) to rename, <ENTER> to edit a profile").SetTextAlign(tview.AlignRight)

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Profiles"), 0, 1, false).
		AddItem(profilesList, 0, 4, true).
		AddItem(
			tview.NewGrid().SetRows(3, 3).
				AddItem(helpInfo1, 1, 0, 1, 1, 0, 0, false).
				AddItem(helpInfo2, 1, 1, 1, 1, 0, 0, false), 0, 4, false).
		SetBorder(true)

	return &ProfileView{View: flex, ProfilesList: profilesList}
}

// NewSitelistView returns a view for the sitelist
func (cactus *Cactus) NewSitelistView() *tview.Flex {

	list := tview.NewList()
	i := 1
	for _, entry := range cactus.SiteList {
		shortcut, _ := indexToRune(i)
		list.AddItem(entry.Name.String(), "", shortcut, nil)
		i += 1
	}

	list.AddItem("Go Back", "return to the previous page", 'b', cactus.UI.OnGoBackSelected)

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Sitelist"), 0, 1, false).
		AddItem(list, 0, 4, true).
		SetBorder(true)

	return flex
}

// NewMainMenuView initializes the menu for the main view
func (ui *UI) NewMainMenuView(welcomeMessage string, entries []MenuEntry) *tview.Flex {

	list := tview.NewList()
	for _, entry := range entries {
		list.AddItem(entry.name, entry.description, entry.label, entry.selected)
	}

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText(welcomeMessage), 0, 1, false).
		AddItem(list, 0, 4, true).
		SetBorder(true)

	return flex
}

// OnGoBackSelected should be called when a user choose to go to the previous view
func (ui *UI) OnGoBackSelected() {

	// switch selection based on current page
	switch currentPageTitle, _ := ui.pages.GetFrontPage(); currentPageTitle {
	case "Sitelist", "Profiles":
		// if current page is Sitelist, go back to main menu
		ui.pages.SwitchToPage("Main Menu")
	}

}

// OnProfilesSelected should be called when a user choose Profiles entry on main menu
func (ui *UI) OnProfilesSelected() {

	// switch current view to Sitelist
	ui.pages.SwitchToPage("Profiles")
}

// OnSitelistSelected should be called when a user choose Sitelist entry on main menu
func (ui *UI) OnSitelistSelected() {

	// switch current view to Sitelist
	ui.pages.SwitchToPage("Sitelist")
}

// NewUI initializes a tview and assigns it to the returned UI object
func NewUI() *UI {
	ui := new(UI)
	ui.tui = tview.NewApplication()

	return ui
}

func (ui *UI) SetListeners() {
	ui.tui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// switch selection based on current page
		switch currentPageTitle, _ := ui.pages.GetFrontPage(); currentPageTitle {
		case "Profiles":
			switch pressedKey := event.Rune(); pressedKey {
			case 97: // user presses 'a' key
				// switch current view to Sitelist
				ui.pages.SwitchToPage("New Profile")
			case 100: // user presses 'd' key
			case 114: // user presses 'r' key
			}
		}
		return event
	})

}
