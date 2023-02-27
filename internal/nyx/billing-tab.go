package nyx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

type BillingTab struct {
	ProfileGroupsList   *widget.List
	Top                 *fyne.Container
	ProfilesTable       *widget.Table
	ProfilesSlice       [][]any
	Bottom              *fyne.Container
	CurrentProfileGroup int64
}

// NewProfilesView returns a view for the profiles management
func (app *Config) billingTab() *fyne.Container {
	// define billingTab
	app.BillingTab = &BillingTab{
		Top:    container.NewMax(),
		Bottom: container.NewHBox(),
	}

	// get current profiles and profile groups
	app.getProfiles()
	app.getProfileGroups()
	if len(app.User.ProfileManager.Groups) > 0 {
		app.CurrentProfileGroup = app.User.ProfileManager.Groups[0].ID
	}
	app.ProfilesSlice = app.getProfilesSlice()

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
		app.refreshProfilesTable()
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
	app.ProfilesTable = app.getProfilesTable()

	// define the center container
	vScroll := container.NewScroll(app.ProfilesTable)

	// define the billingTab container
	profilesTabContainer := container.NewWithoutLayout(leftContainer, app.Top, vScroll, app.Bottom)

	// resize and move billingTab elements
	leftContainer.Move(fyne.NewPos(10, 10))
	leftContainer.Resize(fyne.NewSize(280, 600))

	vScroll.Move(fyne.NewPos(300, 90))
	vScroll.Resize(fyne.NewSize(970, 500))

	app.Top.Move(fyne.NewPos(300, 10))
	app.Top.Resize(fyne.NewSize(900, 64))

	app.Bottom.Move(fyne.NewPos(300, 600))
	app.Bottom.Resize(fyne.NewSize(900, 50))

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
				inserted, err := app.DB.InsertProfileGroup(user.ProfileGroup{
					Name: nameEntry.Text,
				})

				if err != nil {
					app.Logger.Error(err)
				}
				app.CurrentProfileGroup = inserted.ID
				app.refreshProfileGroupsList()
				app.refreshBillingTopContent()
				app.refreshBillingBottomContent()
				app.refreshProfilesTable()
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
					app.refreshProfilesTable()
				}, app.MainWindow)
		}))

	return toolbar
}

func (app *Config) addProfileDialog() dialog.Dialog {
	title := widget.NewEntry()
	title.SetPlaceHolder("My Profile")
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
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

	vBox := container.NewVBox(
		widget.NewLabel("Title*"),
		title,
		widget.NewLabel("First name*"),
		firstName,
		widget.NewLabel("Last name"),
		lastName,
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
	)
	scrollContent := container.NewVScroll(vBox)

	addForm := dialog.NewCustomConfirm(
		"Add New Profile",
		"Create",
		"Cancel",
		scrollContent,
		func(valid bool) {

			for _, o := range vBox.Objects {
				switch o := o.(type) {
				case *widget.Entry:
					err := o.Validate()
					if err != nil {
						valid = false
						break
					}
				}
			}

			if valid {
				var countryCode, stateCode string
				countryCode, err := utils.GetCountryCode(country.Selected)
				if err != nil {
					countryCode = ""
				}
				if !state.Disabled() {
					stateCode, err = utils.GetStateCode(countryCode, state.Selected)
					if err != nil {
						stateCode = ""
					}
				}
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
				app.refreshProfilesTable()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	addForm.Resize(fyne.NewSize(500, 600))
	addForm.Show()

	return addForm
}

func (app *Config) editProfileDialog(pf *user.Profile) dialog.Dialog {
	title := widget.NewEntry()
	title.SetText(pf.Title)
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		if s != pf.Title && app.User.ProfileManager.GetProfileByTitle(s, app.CurrentProfileGroup) != nil {
			return fmt.Errorf("a profile named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else {
			return nil
		}
	}

	firstName := widget.NewEntry()
	firstName.SetText(pf.FirstName)
	firstName.Validator = utils.IsStringEmpty

	lastName := widget.NewEntry()
	lastName.SetText(pf.LastName)
	lastName.Validator = utils.IsStringEmpty

	email := widget.NewEntry()
	email.SetText(pf.Email)
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	phone := widget.NewEntry()
	phone.SetText(pf.Phone)
	phone.Validator = validation.NewRegexp(`^\d+$`, "please insert digits only")

	addressLine1 := widget.NewEntry()
	addressLine1.SetText(pf.AddressLine1)
	addressLine1.Validator = utils.IsStringEmpty

	addressLine2 := widget.NewEntry()
	addressLine2.SetText(pf.AddressLine2)

	city := widget.NewEntry()
	city.SetText(pf.City)
	city.Validator = utils.IsStringEmpty

	postcode := widget.NewEntry()
	postcode.SetText(pf.Postcode)
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
	if pf.CountryCode != "" {
		countryName, err := utils.GetCountryName(pf.CountryCode)
		if err == nil {
			country.SetSelected(countryName)
		}
	}
	if pf.State != "" {
		state.Enable()
		stateName, err := utils.GetStateName(pf.CountryCode, pf.State)
		if err == nil {
			state.SetSelected(stateName)
		}
	}

	cardNumber := widget.NewEntry()
	cardNumber.SetText(pf.CardNumber)
	cardNumber.Validator = validation.NewRegexp(`(^4[0-9]{12}(?:[0-9]{3})?$)|(^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$)|(3[47][0-9]{13})|(^3(?:0[0-5]|[68][0-9])[0-9]{11}$)|(^6(?:011|5[0-9]{2})[0-9]{12}$)|(^(?:2131|1800|35\d{3})\d{11}$)`, "not a valid credit card")

	cardMonth := widget.NewSelect(utils.CardMonths(), func(s string) {
	})
	cardMonth.SetSelected(pf.CardMonth)

	cardYear := widget.NewSelect(utils.CardYears(), func(s string) {
	})
	cardYear.SetSelected(pf.CardYear)

	cardCvv := widget.NewEntry()
	cardCvv.SetText(pf.CardCvv)
	cardCvv.Validator = func(s string) error {
		if len(s) > 5 {
			return errors.New("not a valid CVV")
		}
		return nil
	}

	vBox := container.NewVBox(
		widget.NewLabel("Title*"),
		title,
		widget.NewLabel("First name*"),
		firstName,
		widget.NewLabel("Last name"),
		lastName,
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
	)
	scrollContent := container.NewVScroll(vBox)

	editForm := dialog.NewCustomConfirm(
		"Edit Profile",
		"Update",
		"Cancel",
		scrollContent,
		func(valid bool) {

			for _, o := range vBox.Objects {
				switch o := o.(type) {
				case *widget.Entry:
					err := o.Validate()
					if err != nil {
						valid = false
						break
					}
				}
			}

			if valid {
				var countryCode, stateCode string
				countryCode, err := utils.GetCountryCode(country.Selected)
				if err != nil {
					countryCode = ""
				}
				if !state.Disabled() {
					stateCode, err = utils.GetStateCode(countryCode, state.Selected)
					if err != nil {
						stateCode = ""
					}
				}
				err = app.DB.UpdateProfile(pf.ID,
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
				app.refreshProfilesTable()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	editForm.Resize(fyne.NewSize(500, 600))
	editForm.Show()

	return editForm
}

func (app *Config) getProfilesSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"ID", "Profile", "Email", "Card", "Shipping", "Actions"})

	for _, x := range app.User.ProfileManager.FilterByGroupID(app.CurrentProfileGroup) {
		var currentRow []any

		currentRow = append(currentRow, x.ID)

		if len(x.Title) > 16 {
			currentRow = append(currentRow, x.Title[:12]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		currentRow = append(currentRow, x.Email)
		if len(x.CardNumber) > 5 {
			currentRow = append(currentRow, fmt.Sprintf("**%s", x.CardNumber[len(x.CardNumber)-4:]))
		} else {
			currentRow = append(currentRow, x.CardNumber)
		}

		currentRow = append(currentRow, x.FirstName+" "+x.LastName+"\n"+x.AddressLine1)
		currentRow = append(currentRow, widget.NewToolbar())

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) getProfilesTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			return len(app.ProfilesSlice), len(app.ProfilesSlice[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			ctr.Resize(fyne.Size{Height: 50})
			return ctr
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			// in order: if last column && not the first row
			if tci.Col == (len(app.ProfilesSlice[0])-1) && tci.Row != 0 {
				// last cell - put in a button
				w := widget.NewToolbar(
					widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[tci.Row][0].(int64), app.CurrentProfileGroup)
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							_, err := app.DB.InsertProfile(*pf)
							if err != nil {
								app.Logger.Error(err)
							}
							app.refreshProfileGroupsList()
							app.refreshBillingTopContent()
							app.refreshProfilesTable()
						}
					}),
					widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[tci.Row][0].(int64), app.CurrentProfileGroup)
						if pf != nil {
							app.editProfileDialog(pf)
						}
					}),
					widget.NewToolbarAction(theme.DeleteIcon(), func() {
						dialog.ShowConfirm("Delete?", "", func(deleted bool) {
							if deleted {
								pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[tci.Row][0].(int64), app.CurrentProfileGroup)
								err := app.DB.DeleteProfile(pf.ID)
								if err != nil {
									app.Logger.Error(err)
								}
							}
							app.refreshProfileGroupsList()
							app.refreshBillingTopContent()
							app.refreshProfilesTable()
						}, app.MainWindow)
					}))

				co.(*fyne.Container).Objects = []fyne.CanvasObject{w}
			} else if tci.Col == 0 {
				// we're just putting in textual information
				co.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(""),
				}
			} else {
				// we're just putting in textual information
				co.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(app.ProfilesSlice[tci.Row][tci.Col].(string)),
				}
			}
		})

	// hide first column (ID)
	colWidths := []float32{0, 100, 200, 200, 200, 60}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	for i := 1; i < len(app.ProfilesSlice); i++ {
		t.SetRowHeight(i, 55)
	}

	return t
}

func (app *Config) refreshProfilesTable() {
	app.ProfilesSlice = app.getProfilesSlice()
	app.ProfilesTable.Refresh()

	// hide first column (ID)
	colWidths := []float32{0, 100, 200, 200, 200, 60}
	for i, w := range colWidths {
		app.ProfilesTable.SetColumnWidth(i, w)
	}

	for i := 1; i < len(app.ProfilesSlice); i++ {
		app.ProfilesTable.SetRowHeight(i, 55)
	}
}

func (app *Config) refreshBillingBottomContent() {

	if len(app.User.ProfileManager.Groups) > 0 {
		btnAdd := widget.NewButtonWithIcon("Add Profile", theme.ContentAddIcon(), func() {
			app.addProfileDialog()
		})
		btnClear := widget.NewButtonWithIcon("Clear Profiles", theme.ContentRemoveIcon(), func() {
			dialog.ShowConfirm(
				"Delete?",
				fmt.Sprintf("Do you really want to delete %d profiles?", len(app.User.ProfileManager.FilterByGroupID(app.CurrentProfileGroup))),
				func(deleted bool) {
					if deleted {
						err := app.DB.DeleteProfilesByGroupID(app.CurrentProfileGroup)
						if err != nil {
							app.Logger.Error(err)
						}
						app.refreshProfileGroupsList()
						app.refreshBillingTopContent()
						app.refreshBillingBottomContent()
						app.refreshProfilesTable()
					}
				}, app.MainWindow)
		})
		btnClear.Importance = widget.DangerImportance

		app.Bottom.Objects = []fyne.CanvasObject{
			layout.NewSpacer(),
			btnAdd,
			btnClear,
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
