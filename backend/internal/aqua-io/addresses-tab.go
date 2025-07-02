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

type AddressesTab struct {
	*container.TabItem
	AddressesTop    *fyne.Container
	AddressesTable  *widget.Table
	AddressesSlice  [][]any
	AddressesBottom *fyne.Container
}

// NewProfilesView returns a view for the profiles management
func (app *Config) addressesTab() *fyne.Container {
	// define profilesTab
	app.AddressesTab.AddressesTop = container.NewStack()
	app.AddressesTab.AddressesBottom = container.NewHBox()

	app.AddressesSlice = app.getAddressesSlice()

	// update content
	app.refreshAddressesTopContent()
	app.refreshAddressesBottomContent()

	// get the profiles table
	app.AddressesTable = app.getAddressesTable()

	// define the center container
	vScroll := container.NewScroll(app.AddressesTable)

	// define the profilesTab container
	ctr := container.NewWithoutLayout(app.AddressesTop, vScroll, app.AddressesBottom)

	app.AddressesTop.Move(fyne.NewPos(10, 10))
	app.AddressesTop.Resize(fyne.NewSize(1280, 40))

	vScroll.Move(fyne.NewPos(10, 60))
	vScroll.Resize(fyne.NewSize(1280, 500))

	app.AddressesBottom.Move(fyne.NewPos(300, 580))
	app.AddressesBottom.Resize(fyne.NewSize(900, 50))

	return ctr
}

func (app *Config) addAddressDialog() dialog.Dialog {

	title := widget.NewEntry()
	title.SetPlaceHolder("My Address")
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		p, _ := app.DB.GetAddressByTitle(s)
		if p != nil {
			return fmt.Errorf("An address with title \"%s\" already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	firstName := widget.NewEntry()
	firstName.SetPlaceHolder("Mario")
	firstName.Validator = utils.IsStringEmpty

	lastName := widget.NewEntry()
	lastName.SetPlaceHolder("Rossi")
	lastName.Validator = utils.IsStringEmpty

	phone := widget.NewEntry()
	phone.SetPlaceHolder("1234567890")
	phone.Validator = utils.IsStringEmpty

	addressLine1 := widget.NewEntry()
	addressLine1.SetPlaceHolder("via Aqua 1")
	addressLine1.Validator = utils.IsStringEmpty

	addressLine2 := widget.NewEntry()
	addressLine2.SetPlaceHolder("")
	addressLine2.Validator = utils.IsStringEmpty

	zipCode := widget.NewEntry()
	zipCode.SetPlaceHolder("123 XXX")
	zipCode.Validator = utils.IsStringEmpty

	province := widget.NewEntry()
	province.SetPlaceHolder("FC")
	province.Validator = utils.IsStringEmpty

	countryCode := widget.NewEntry()
	countryCode.SetPlaceHolder("IT")
	countryCode.Validator = utils.IsStringEmpty

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("First Name"),
		firstName,
		widget.NewLabel("Last Name"),
		lastName,
		widget.NewLabel("Phone"),
		phone,
		widget.NewLabel("Address"),
		addressLine1,
		widget.NewLabel("Address line 2"),
		addressLine2,
		widget.NewLabel("Zip Code"),
		zipCode,
		widget.NewLabel("Province"),
		province,
		widget.NewLabel("Country Code"),
		countryCode,
	)
	scrollContent := container.NewVScroll(vBox)

	addForm := dialog.NewCustomConfirm(
		"Add New Address",
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
				p := user.Address{
					Title:        title.Text,
					FirstName:    firstName.Text,
					LastName:     lastName.Text,
					Phone:        phone.Text,
					AddressLine1: addressLine1.Text,
					AddressLine2: addressLine2.Text,
					ZipCode:      zipCode.Text,
					Province:     province.Text,
					CountryCode:  countryCode.Text,
				}

				_, err := app.DB.InsertAddress(p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshAddressesTab()
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

func (app *Config) editAddressDialog(pf *user.Address) dialog.Dialog {
	title := widget.NewEntry()
	title.SetText(pf.Title)
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		p, _ := app.DB.GetAddressByTitle(s)
		if s != pf.Title && p != nil {
			return fmt.Errorf("An address with title \"%s\" already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
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

	phone := widget.NewEntry()
	phone.SetText(pf.Phone)
	phone.Validator = utils.IsStringEmpty

	addressLine1 := widget.NewEntry()
	addressLine1.SetText(pf.AddressLine1)
	addressLine1.Validator = utils.IsStringEmpty

	addressLine2 := widget.NewEntry()
	addressLine2.SetText(pf.AddressLine2)
	addressLine2.Validator = utils.IsStringEmpty

	zipCode := widget.NewEntry()
	zipCode.SetText(pf.ZipCode)
	zipCode.Validator = utils.IsStringEmpty

	province := widget.NewEntry()
	province.SetText(pf.Province)
	province.Validator = utils.IsStringEmpty

	countryCode := widget.NewEntry()
	countryCode.SetText(pf.CountryCode)
	countryCode.Validator = utils.IsStringEmpty

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("First Name"),
		firstName,
		widget.NewLabel("Last Name"),
		lastName,
		widget.NewLabel("Phone"),
		phone,
		widget.NewLabel("Address"),
		addressLine1,
		widget.NewLabel("Address line 2"),
		addressLine2,
		widget.NewLabel("Zip Code"),
		zipCode,
		widget.NewLabel("Province"),
		province,
		widget.NewLabel("Country Code"),
		countryCode,
	)
	scrollContent := container.NewVScroll(vBox)

	editForm := dialog.NewCustomConfirm(
		"Edit Address",
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
				p := user.Address{
					Title:        title.Text,
					FirstName:    firstName.Text,
					LastName:     lastName.Text,
					Phone:        phone.Text,
					AddressLine1: addressLine1.Text,
					AddressLine2: addressLine2.Text,
					ZipCode:      zipCode.Text,
					Province:     province.Text,
					CountryCode:  countryCode.Text,
				}

				err := app.DB.UpdateAddress(pf.ID, p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshAddressesTab()
				app.refreshProfilesTab()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	editForm.Resize(fyne.NewSize(500, 600))
	editForm.Show()

	return editForm

}

func (app *Config) getAddressesSlice() [][]any {
	var slice [][]any
	slice = append(slice, []any{"Title", "First Name", "Last Name", "Phone", "Address", "Zip Code", "Province", "Country", "Actions"})
	allAddresses, _ := app.DB.AllAddresses()

	for _, x := range allAddresses {
		var currentRow []any

		if len(x.Title) > 20 {
			currentRow = append(currentRow, x.Title[:19]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		currentRow = append(currentRow, x.FirstName)

		currentRow = append(currentRow, x.LastName)

		currentRow = append(currentRow, x.Phone)

		if len(x.AddressLine1) > 40 {
			currentRow = append(currentRow, x.AddressLine1[:39]+"...")
		} else {
			currentRow = append(currentRow, x.AddressLine1)
		}

		currentRow = append(currentRow, x.ZipCode)

		currentRow = append(currentRow, x.Province)

		currentRow = append(currentRow, x.CountryCode)

		currentRow = append(currentRow, x.ID) // for toolbar actions

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) getAddressesTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(app.AddressesSlice), len(app.AddressesSlice[0])
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

			if i.Col == (len(app.AddressesSlice[0])-1) && i.Row != 0 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
						pf, _ := app.DB.GetAddressByID(app.AddressesSlice[i.Row][len(app.AddressesSlice[0])-1].(int64))
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							_, err := app.DB.InsertAddress(*pf)
							if err != nil {
								app.Logger.Error(err)
							}
							app.refreshAddressesTab()
							app.refreshProfilesTab()
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						pf, _ := app.DB.GetAddressByID(app.AddressesSlice[i.Row][len(app.AddressesSlice[0])-1].(int64))
						if pf != nil {
							app.editAddressDialog(pf)
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
						pf, _ := app.DB.GetAddressByID(app.AddressesSlice[i.Row][len(app.AddressesSlice[0])-1].(int64))
						if pf == nil {
							return
						}

						dialog.ShowConfirm("Delete?", "", func(deleted bool) {
							if deleted {
								err := app.DB.DeleteAddress(pf.ID)
								if err != nil {
									app.Logger.Error(err)
								}
							}
							app.refreshAddressesTab()
							app.refreshProfilesTab()
						}, app.MainWindow)

					}))
				}
			} else {
				toolbar.Hide()
				lbl.Hidden = false
				// we're just putting in textual information
				lbl.SetText(
					app.AddressesSlice[i.Row][i.Col].(string))
			}
		})

	colWidths := []float32{120, 100, 100, 100, 220, 100, 100, 100, 40}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (app *Config) refreshAddressesTable() {
	app.AddressesSlice = app.getAddressesSlice()
	app.AddressesTable.Refresh()

	colWidths := []float32{120, 100, 100, 100, 220, 100, 100, 100, 40}
	for i, w := range colWidths {
		app.AddressesTable.SetColumnWidth(i, w)
	}
}

func (app *Config) refreshAddressesBottomContent() {

	all, _ := app.DB.AllAddresses()

	btnAdd := widget.NewButtonWithIcon("Add Address", theme.ContentAddIcon(), func() {
		app.addAddressDialog()
	})
	btnClear := widget.NewButtonWithIcon("Clear Addresses", theme.ContentRemoveIcon(), func() {
		dialog.ShowConfirm(
			"Delete all addresses?",
			fmt.Sprintf("Do you really want to delete %d addresses?", len(all)),
			func(deleted bool) {
				if deleted {
					for _, p := range all {
						err := app.DB.DeleteAddress(p.ID)
						if err != nil {
							app.Logger.Error(err)
						}
					}
					app.refreshAddressesTab()
					app.refreshProfilesTab()
				}
			}, app.MainWindow)
	})
	btnClear.Importance = widget.DangerImportance

	app.AddressesBottom.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		btnAdd,
		btnClear,
	}

	app.AddressesBottom.Refresh()
}

func (app *Config) refreshAddressesTopContent() {

	all, _ := app.DB.AllAddresses()

	txt := widget.NewRichTextFromMarkdown(`## ` + strconv.Itoa(len(all)) + ` Addresses Loaded`)
	app.AddressesTop.Objects = []fyne.CanvasObject{
		txt,
	}

	app.AddressesTop.Refresh()
}

func (app *Config) refreshAddressesTab() {
	app.refreshAddressesTopContent()
	app.refreshAddressesTable()
}
