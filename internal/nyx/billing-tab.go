package nyx

import (
	"errors"
	"fmt"
	"regexp"
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

	// get current profile groups
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
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if label, ok := item.(*widget.Label); ok {
				label.SetText(
					fmt.Sprintf("%s\t (%d)", app.User.ProfileManager.Groups[id].Name, len(app.User.ProfileManager.FilterByGroupName(app.User.ProfileManager.Groups[id].Name))))
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		app.CurrentProfileGroup = app.User.ProfileManager.Groups[id].ID
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

func (app *Config) getProfileGroups() {
	groups, err := app.DB.AllProfileGroups()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.ProfileManager.Groups = groups
}

func (app *Config) refreshProfileGroupsList() {
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
			} else {
				state.Options = []string{}
				state.ClearSelected()
				state.Disable()
			}
		}
	})

	cardNumber := widget.NewEntry()
	cardNumber.SetPlaceHolder("")
	cardNumber.Validator = func(s string) error {
		pattern := `(^4[0-9]{12}(?:[0-9]{3})?$)|(^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$)|(3[47][0-9]{13})|(^3(?:0[0-5]|[68][0-9])[0-9]{11}$)|(^6(?:011|5[0-9]{2})[0-9]{12}$)|(^(?:2131|1800|35\d{3})\d{11}$)`
		if len(s) > 0 {
			_, e := regexp.MatchString(pattern, s)
			if e != nil {
				return errors.New("not a valid credit card")
			}
		}
		return nil
	}

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

	// create a dialog
	addForm := dialog.NewForm(
		"Add New Profile",
		"Create",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Title", Widget: title},
			{Text: "Email", Widget: email},
			{Text: "Phone", Widget: phone},
			{Text: "Address Line 1", Widget: addressLine1},
			{Text: "Address Line 2", Widget: addressLine2},
			{Text: "City", Widget: city},
			{Text: "Postcode", Widget: postcode},
			{Text: "State", Widget: state},
			{Text: "Country", Widget: country},
			{Text: "Card Number", Widget: cardNumber},
			{Text: "Card Month", Widget: cardMonth},
			{Text: "Card Year", Widget: cardYear},
			{Text: "Card CVV", Widget: cardCvv},
		},
		func(valid bool) {
			if valid {

			}
		},
		app.MainWindow)

	// size and show the dialog
	addForm.Resize(fyne.NewSize(800, 700))
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
