package aqua_io

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/UrbiJr/aqua-io/backend/internal/user"
	"github.com/UrbiJr/aqua-io/backend/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProfilesTab struct {
	*container.TabItem
	ProfileTop           *fyne.Container
	ProfilesTable *widget.Table
	ProfilesSlice [][]any
	ProfileBottom        *fyne.Container
}

// NewProfilesView returns a view for the profiles management
func (app *Config) profilesTab() *fyne.Container {
	// define profilesTab
	app.ProfilesTab.ProfileTop = container.NewStack()
	app.ProfilesTab.ProfileBottom = container.NewHBox()

	app.ProfilesSlice = app.getProfilesSlice()

	// update content
	app.refreshProfilesTopContent()
	app.refreshProfilesBottomContent()

	// get the profiles table
	app.ProfilesTable = app.getProfilesTable()

	// define the center container
	vScroll := container.NewScroll(app.ProfilesTable)

	// define the profilesTab container
	profilesTabContainer := container.NewWithoutLayout(app.ProfileTop, vScroll, app.ProfileBottom)

	app.ProfileTop.Move(fyne.NewPos(10, 10))
	app.ProfileTop.Resize(fyne.NewSize(1280, 40))

	vScroll.Move(fyne.NewPos(10, 60))
	vScroll.Resize(fyne.NewSize(1280, 500))

	app.ProfileBottom.Move(fyne.NewPos(300, 580))
	app.ProfileBottom.Resize(fyne.NewSize(900, 50))

	return profilesTabContainer
}

func (app *Config) getAllProfilesTitles() ([]string, error) {
	profiles, err := app.DB.AllProfiles()
	if err != nil {
		return nil, err
	}

	var titles []string
	for _, p := range profiles {
		titles = append(titles, p.Title)
	}

	return titles, nil
}

func (app *Config) addProfileDialog() dialog.Dialog {

	title := widget.NewEntry()
	title.SetPlaceHolder("My Profile")
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		p, _ := app.DB.GetProfileByTitle(s)
		if p != nil {
			return fmt.Errorf("a profile named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	var addressesIDs []string

	addresses, err := app.DB.AllAddresses()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for _, a := range addresses {
			addressesIDs = append(addressesIDs, a.Title)
		}
	}

	billingAddressSelect := widget.NewSelect(addressesIDs, nil)
	billingAddressSelect.ClearSelected()

	shippingAddressSelect := widget.NewSelect(addressesIDs, nil)
	shippingAddressSelect.ClearSelected()

	cardNumber := widget.NewEntry()
	cardNumber.SetPlaceHolder("XXXX XXXX XXXX XX34")
	cardNumber.Validator = utils.IsStringEmpty

	cardMonth := widget.NewEntry()
	cardMonth.SetPlaceHolder("12")
	cardMonth.Validator = utils.IsStringEmpty

	cardYear := widget.NewEntry()
	cardYear.SetPlaceHolder("29")
	cardYear.Validator = utils.IsStringEmpty

	cardCvv := widget.NewEntry()
	cardCvv.SetPlaceHolder("123")
	cardCvv.Validator = utils.IsStringEmpty

	testMode := widget.NewCheck("", func(b bool) {})

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Billing Address"),
		billingAddressSelect,
		widget.NewLabel("Shipping Address"),
		shippingAddressSelect,
		widget.NewLabel("Card Number"),
		cardNumber,
		widget.NewLabel("Card Month"),
		cardMonth,
		widget.NewLabel("Card Year"),
		cardYear,
		widget.NewLabel("CVV"),
		cardCvv,
		widget.NewLabel("Test"),
		testMode,
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

			billingAddress, _ := app.DB.GetAddressByTitle(billingAddressSelect.Selected)
			shippingAddress, _ := app.DB.GetAddressByTitle(shippingAddressSelect.Selected)

			if valid {
				p := user.Profile{
					Title:             title.Text,
					BillingAddressID:  billingAddress.ID,
					ShippingAddressID: shippingAddress.ID,
					CardNumber:        cardNumber.Text,
					CardMonth:         cardMonth.Text,
					CardYear:          cardYear.Text,
					CardCvv:           cardCvv.Text,
					TestMode:          testMode.Checked,
				}

				_, err := app.DB.InsertProfile(p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfilesTab()
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
		p, _ := app.DB.GetProfileByTitle(s)
		if s != pf.Title && p != nil {
			return fmt.Errorf("a profile named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	var addressesIDs []string

	addresses, err := app.DB.AllAddresses()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for _, a := range addresses {
			addressesIDs = append(addressesIDs, strconv.Itoa(int(a.ID)))
		}
	}

	billingAddressSelect := widget.NewSelect(addressesIDs, nil)
	billingAddressSelect.SetSelected(string(pf.BillingAddressID))

	shippingAddressSelect := widget.NewSelect(addressesIDs, nil)
	shippingAddressSelect.SetSelected(string(pf.ShippingAddressID))

	cardNumber := widget.NewEntry()
	cardNumber.SetText(pf.CardNumber)
	cardNumber.Validator = utils.IsStringEmpty

	cardMonth := widget.NewEntry()
	cardMonth.SetText(pf.CardMonth)
	cardMonth.Validator = utils.IsStringEmpty

	cardYear := widget.NewEntry()
	cardYear.SetText(pf.CardYear)
	cardYear.Validator = utils.IsStringEmpty

	cardCvv := widget.NewEntry()
	cardCvv.SetText(pf.CardCvv)
	cardCvv.Validator = utils.IsStringEmpty

	testMode := widget.NewCheck("", func(b bool) {})
	testMode.SetChecked(pf.TestMode)

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Billing Address"),
		billingAddressSelect,
		widget.NewLabel("Shipping Address"),
		shippingAddressSelect,
		widget.NewLabel("Card Number"),
		cardNumber,
		widget.NewLabel("Card Month"),
		cardMonth,
		widget.NewLabel("Card Year"),
		cardYear,
		widget.NewLabel("CVV"),
		cardCvv,
		widget.NewLabel("Test"),
		testMode,
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

			billingAddressID, _ := strconv.Atoi(billingAddressSelect.Selected)
			shippingAddressID, _ := strconv.Atoi(shippingAddressSelect.Selected)

			if valid {
				p := user.Profile{
					Title:             title.Text,
					BillingAddressID:  int64(billingAddressID),
					ShippingAddressID: int64(shippingAddressID),
					CardNumber:        cardNumber.Text,
					CardMonth:         cardMonth.Text,
					CardYear:          cardYear.Text,
					CardCvv:           cardCvv.Text,
					TestMode:          testMode.Checked,
				}

				err := app.DB.UpdateProfile(pf.ID, p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfilesTab()
				app.refreshTasksTab()
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
	slice = append(slice, []any{"Title", "Billing Address", "Shipping Address", "Card Number", "Test", "Actions"})
	allProfiles, _ := app.DB.AllProfiles()

	for _, x := range allProfiles {
		var currentRow []any

		if len(x.Title) > 30 {
			currentRow = append(currentRow, x.Title[:29]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		billingAddress, err := app.DB.GetAddressByID(x.BillingAddressID)
		if err != nil {
			continue
		}
		if len(billingAddress.Title) > 30 {
			currentRow = append(currentRow, billingAddress.Title[:29]+"...")
		} else {
			currentRow = append(currentRow, billingAddress.Title)
		}

		shippingAddress, err := app.DB.GetAddressByID(x.ShippingAddressID)
		if err != nil {
			continue
		}
		if len(shippingAddress.Title) > 30 {
			currentRow = append(currentRow, shippingAddress.Title[:29]+"...")
		} else {
			currentRow = append(currentRow, shippingAddress.Title)
		}

		currentRow = append(currentRow, x.CardNumber[:12]+"XXXX")

		currentRow = append(currentRow, x.TestMode)

		currentRow = append(currentRow, x.ID) // for toolbar actions

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
			lbl := widget.NewLabel("")
			toolbar := widget.NewToolbar()
			toolbar.Hide()
			return container.NewStack(lbl, toolbar)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			toolbar := container.Objects[1].(*widget.Toolbar)

			if i.Row != 0 && i.Col == 5 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
						pf, _ := app.DB.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							_, err := app.DB.InsertProfile(*pf)
							if err != nil {
								app.Logger.Error(err)
							}
							app.refreshProfilesTab()
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						pf, _ := app.DB.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf != nil {
							app.editProfileDialog(pf)
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
						pf, _ := app.DB.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf == nil {
							return
						}

						dialog.ShowConfirm("Delete?", "", func(deleted bool) {
							if deleted {
								err := app.DB.DeleteProfile(pf.ID)
								if err != nil {
									app.Logger.Error(err)
								}
							}
							app.refreshProfilesTab()
						}, app.MainWindow)

					}))
				}
			} else if i.Col == 4 && i.Row != 0 {
				toolbar.Hide()
				lbl.Hidden = false
				if app.ProfilesSlice[i.Row][i.Col].(bool) == true {
					lbl.SetText("Yes")
				} else {
					lbl.SetText("No")
				}
			} else {
				toolbar.Hide()
				lbl.Hidden = false
				// we're just putting in textual information
				lbl.SetText(
					app.ProfilesSlice[i.Row][i.Col].(string))
			}
		})

	colWidths := []float32{220, 270, 200, 200, 100, 40}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (app *Config) refreshProfilesTable() {
	app.ProfilesSlice = app.getProfilesSlice()
	app.ProfilesTable.Refresh()

	colWidths := []float32{220, 270, 200, 200, 100, 40}
	for i, w := range colWidths {
		app.ProfilesTable.SetColumnWidth(i, w)
	}
}

func (app *Config) refreshProfilesBottomContent() {

	allProfiles, _ := app.DB.AllProfiles()

	btnAdd := widget.NewButtonWithIcon("Add Profile", theme.ContentAddIcon(), func() {
		app.addProfileDialog()
	})
	btnClear := widget.NewButtonWithIcon("Clear Profiles", theme.ContentRemoveIcon(), func() {
		dialog.ShowConfirm(
			"Delete all profiles?",
			fmt.Sprintf("Do you really want to delete %d profiles?", len(allProfiles)),
			func(deleted bool) {
				if deleted {
					for _, p := range allProfiles {
						err := app.DB.DeleteProfile(p.ID)
						if err != nil {
							app.Logger.Error(err)
						}
					}
					app.refreshProfilesTab()
				}
			}, app.MainWindow)
	})
	btnClear.Importance = widget.DangerImportance

	app.ProfileBottom.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		btnAdd,
		btnClear,
	}

	app.ProfileBottom.Refresh()
}

func (app *Config) refreshProfilesTopContent() {

	allProfiles, _ := app.DB.AllProfiles()

	txt := widget.NewRichTextFromMarkdown(`## ` + strconv.Itoa(len(allProfiles)) + ` Profiles Loaded`)
	app.ProfileTop.Objects = []fyne.CanvasObject{
		txt,
	}

	app.ProfileTop.Refresh()
}

func (app *Config) refreshProfilesTab() {
	app.refreshProfilesTopContent()
	app.refreshProfilesTable()
}
