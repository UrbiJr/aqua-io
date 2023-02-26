package nyx

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

type BillingTab struct {
	ProfileGroupsList   *widget.List
	Top                 *fyne.Container
	ProfilesTable       *widget.Table
	Bottom              *fyne.Container
	CurrentProfileGroup int64
}

// NewProfilesView returns a view for the profiles management
func (app *Config) billingTab() *fyne.Container {
	// define billingTab
	app.BillingTab = &BillingTab{
		Top:    container.NewMax(),
		Bottom: container.NewMax(),
	}

	// get current profiles and profile groups
	app.getProfiles()
	app.getProfileGroups()
	if len(app.User.ProfileManager.Groups) > 0 {
		app.CurrentProfileGroup = app.User.ProfileManager.Groups[0].ID
	}

	// define a list to display profile groups
	list := widget.NewList(
		func() int {
			return len(app.User.ProfileManager.Groups)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(idx widget.ListItemID, item fyne.CanvasObject) {
			if label, ok := item.(*widget.Label); ok {
				label.SetText(
					fmt.Sprintf("%s\t (%d)", app.User.ProfileManager.Groups[idx].Name, len(app.User.ProfileManager.FilterByGroupID(app.User.ProfileManager.Groups[idx].ID))))
			}
		},
	)
	list.OnSelected = func(idx widget.ListItemID) {
		app.CurrentProfileGroup = app.User.ProfileManager.Groups[idx].ID
		app.refreshBillingTopContent()
		app.refreshBillingBottomContent()
	}
	app.ProfileGroupsList = list

	// define a button to create a new group
	addProfileGroupButton := widget.NewButtonWithIcon("New Group", theme.ContentAddIcon(), func() {
		app.addProfileGroupDialog()
	})

	// define the left container
	// use a border layout instead of VBox because VBox and HBox collapse their content to minsize with an equal width
	leftContainer := container.NewBorder(
		container.NewVBox(canvas.NewText("PROFILE GROUPS", nil), addProfileGroupButton), // top
		nil,  // bottom
		nil,  // left
		nil,  // right
		list) // scroll content (center)

	// update content based on currently selected group
	app.refreshBillingTopContent()
	app.refreshBillingBottomContent()

	// get the profiles table

	// define the center container
	centerContainer := container.NewVBox()

	// define the billingTab container
	profilesTabContainer := container.NewWithoutLayout(leftContainer, app.Top, centerContainer, app.Bottom)

	// resize and move billingTab elements
	leftContainer.Move(fyne.NewPos(10, 10))
	leftContainer.Resize(fyne.NewSize(280, 600))

	app.Top.Move(fyne.NewPos(300, 10))
	app.Top.Resize(fyne.NewSize(900, 64))

	app.Bottom.Move(fyne.NewPos(300, 600))
	app.Bottom.Resize(fyne.NewSize(900, 64))

	return profilesTabContainer
}

func (app *Config) addProfileGroupDialog() dialog.Dialog {
	nameEntry := widget.NewEntry()

	nameEntry.Validator = utils.IsStringEmpty

	// create a dialog
	addForm := dialog.NewForm(
		"Create Group",
		"Create",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Group name", Widget: nameEntry},
		},
		func(valid bool) {
			if valid {
				_, err := app.DB.InsertProfileGroup(user.ProfileGroup{
					Name: nameEntry.Text,
				})

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfileGroupsList()
			}
		},
		app.MainWindow)

	// size and show the dialog
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}

func (app *Config) getProfiles() {
	profiles, err := app.DB.AllProfiles()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.ProfileManager.Profiles = profiles
}

func (app *Config) getProfileGroups() {
	groups, err := app.DB.AllProfileGroups()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.ProfileManager.Groups = groups
}

func (app *Config) refreshProfileGroupsList() {
	app.getProfiles()
	app.getProfileGroups()
	app.ProfileGroupsList.Refresh()
}

func (app *Config) getProfileGroupToolBar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			dialog.ShowConfirm(
				"Delete?",
				fmt.Sprintf("Deleting the group will ALSO delete its %d profiles.", len(app.User.ProfileManager.FilterByGroupID(app.CurrentProfileGroup))),
				func(deleted bool) {
					if deleted {
						err := app.DB.DeleteProfileGroup(app.CurrentProfileGroup)
						if err != nil {
							app.Logger.Error(err)
						}
					}
					app.refreshProfileGroupsList()
					app.refreshBillingTopContent()
					app.refreshBillingBottomContent()
				}, app.MainWindow)
		}))

	return toolbar
}

func (app *Config) addProfileDialog() dialog.Dialog {
	title := widget.NewEntry()
	title.SetPlaceHolder("My Profile")
	title.Validator = func(s string) error {
		if app.User.ProfileManager.GetProfileByTitle(s, app.CurrentProfileGroup) != nil {
			return fmt.Errorf("a profile named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else {
			return nil
		}
	}

	firstName := widget.NewEntry()
	firstName.SetPlaceHolder("John")
	firstName.Validator = utils.IsStringEmpty

	lastName := widget.NewEntry()
	lastName.SetPlaceHolder("Smith")
	lastName.Validator = utils.IsStringEmpty

	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	phone := widget.NewEntry()
	phone.SetPlaceHolder("3389384736")
	phone.Validator = validation.NewRegexp(`^\d+$`, "please insert digits only")

	addressLine1 := widget.NewEntry()
	addressLine1.SetPlaceHolder("104 Chapter Road")
	addressLine1.Validator = utils.IsStringEmpty

	addressLine2 := widget.NewEntry()
	addressLine2.SetPlaceHolder("")

	city := widget.NewEntry()
	city.SetPlaceHolder("London")
	city.Validator = utils.IsStringEmpty

	postcode := widget.NewEntry()
	postcode.SetPlaceHolder("NW2 5NG")
	postcode.Validator = utils.IsStringEmpty

	stateLabel := widget.NewLabel("State")
	state := widget.NewSelect([]string{}, func(s string) {

	})
	state.Disable()

	country := widget.NewSelect(utils.CountryNames, func(s string) {
		countryCode, err := utils.GetCountryCode(s)
		if err == nil {
			states, err := utils.GetStates(countryCode)
			if err == nil {
				state.Options = states
				state.Enable()
				stateLabel.SetText("State*")
			} else {
				state.Options = []string{}
				state.ClearSelected()
				state.Disable()
				stateLabel.SetText("State")
			}
		}
	})

	cardNumber := widget.NewEntry()
	cardNumber.SetPlaceHolder("")
	cardNumber.Validator = validation.NewRegexp(`(^4[0-9]{12}(?:[0-9]{3})?$)|(^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$)|(3[47][0-9]{13})|(^3(?:0[0-5]|[68][0-9])[0-9]{11}$)|(^6(?:011|5[0-9]{2})[0-9]{12}$)|(^(?:2131|1800|35\d{3})\d{11}$)`, "not a valid credit card")

	cardMonth := widget.NewSelect(utils.CardMonths(), func(s string) {
	})

	cardYear := widget.NewSelect(utils.CardYears(), func(s string) {
	})

	cardCvv := widget.NewEntry()
	cardCvv.SetPlaceHolder("")
	cardCvv.Validator = func(s string) error {
		if len(s) > 5 {
			return errors.New("not a valid CVV")
		}
		return nil
	}

	scrollContent := container.NewVScroll(container.NewVBox(
		widget.NewLabel("Title*"),
		title,
		widget.NewLabel("Email*"),
		email,
		widget.NewLabel("Phone*"),
		phone,
		widget.NewLabel("Address Line 1*"),
		addressLine1,
		widget.NewLabel("Address Line 2"),
		addressLine2,
		widget.NewLabel("City*"),
		city,
		widget.NewLabel("Postcode*"),
		postcode,
		stateLabel,
		state,
		widget.NewLabel("Country*"),
		country,
		widget.NewLabel("Card Number*"),
		cardNumber,
		widget.NewLabel("Card Month*"),
		cardMonth,
		widget.NewLabel("Card Year*"),
		cardYear,
		widget.NewLabel("Card CVV*"),
		cardCvv,
	))

	addForm := dialog.NewCustomConfirm(
		"Add New Profile",
		"Create",
		"Cancel",
		scrollContent,
		func(valid bool) {
			var countryCode, stateCode string
			countryCode, err := utils.GetCountryCode(country.Selected)
			if err != nil {
				dialog.ShowError(errors.New("please choose a country"), app.MainWindow)
				valid = false
			}
			if !state.Disabled() {
				stateCode, err = utils.GetStateCode(countryCode, state.Selected)
				if err != nil {
					dialog.ShowError(errors.New("please choose a state"), app.MainWindow)
					valid = false
				}
			}

			if valid {
				_, err = app.DB.InsertProfile(
					user.Profile{
						GroupID:      app.CurrentProfileGroup,
						Title:        title.Text,
						Email:        email.Text,
						FirstName:    firstName.Text,
						LastName:     lastName.Text,
						AddressLine1: addressLine1.Text,
						AddressLine2: addressLine2.Text,
						City:         city.Text,
						Postcode:     postcode.Text,
						State:        stateCode,
						CountryCode:  countryCode,
						Phone:        phone.Text,
						CardNumber:   cardNumber.Text,
						CardMonth:    cardMonth.Selected,
						CardYear:     cardYear.Selected,
						CardCvv:      cardCvv.Text,
					})

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfileGroupsList()
				app.refreshBillingTopContent()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	addForm.Resize(fyne.NewSize(500, 600))
	addForm.Show()

	return addForm
}

func (app *Config) refreshBillingBottomContent() {

	if len(app.User.ProfileManager.Groups) > 0 {
		btnAdd := widget.NewButtonWithIcon("Add Profile", theme.ContentAddIcon(), func() {
			app.addProfileDialog()
		})

		app.Bottom.Objects = []fyne.CanvasObject{
			btnAdd,
		}
	} else {
		app.Bottom.Objects = []fyne.CanvasObject{}
	}

	app.Bottom.Refresh()
}

func (app *Config) refreshBillingTopContent() {

	if len(app.User.ProfileManager.Groups) > 0 {
		if app.User.ProfileManager.GetGroupByID(app.CurrentProfileGroup) == nil {
			app.CurrentProfileGroup = app.User.ProfileManager.Groups[0].ID
		}
		txt := widget.NewRichTextFromMarkdown(
			`# ` + app.User.ProfileManager.GetGroupByID(app.CurrentProfileGroup).Name + `

## ` + strconv.Itoa(len(app.User.ProfileManager.FilterByGroupID(app.CurrentProfileGroup))) + ` Profiles Loaded`)
		toolbar := app.getProfileGroupToolBar()
		app.Top.Objects = []fyne.CanvasObject{
			txt,
			toolbar,
		}
	} else {
		app.Top.Objects = []fyne.CanvasObject{}
	}

	app.Top.Refresh()
}
