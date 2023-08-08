package aqua_io

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"
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
	Top           *fyne.Container
	ProfilesTable *widget.Table
	ProfilesSlice [][]any
	Bottom        *fyne.Container
}

// NewProfilesView returns a view for the profiles management
func (app *Config) profilesTab() *fyne.Container {
	// define profilesTab
	app.ProfilesTab.Top = container.NewMax()
	app.ProfilesTab.Bottom = container.NewHBox()

	app.ProfilesSlice = app.getProfilesSlice()

	// update content
	app.refreshProfilesTopContent()
	app.refreshProfilesBottomContent()

	// get the profiles table
	app.ProfilesTable = app.getProfilesTable()

	// define the center container
	vScroll := container.NewScroll(app.ProfilesTable)

	// define the profilesTab container
	profilesTabContainer := container.NewWithoutLayout(app.Top, vScroll, app.Bottom)

	app.Top.Move(fyne.NewPos(10, 10))
	app.Top.Resize(fyne.NewSize(1280, 40))

	vScroll.Move(fyne.NewPos(10, 60))
	vScroll.Resize(fyne.NewSize(1280, 500))

	app.Bottom.Move(fyne.NewPos(300, 580))
	app.Bottom.Resize(fyne.NewSize(900, 50))

	return profilesTabContainer
}

func (app *Config) addProfileDialog() dialog.Dialog {

	title := widget.NewEntry()
	title.SetPlaceHolder("My Profile")
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		if app.User.ProfileManager.GetProfileByTitle(s) != nil {
			return fmt.Errorf("a profile named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	exchange := widget.NewSelect([]string{
		constants.ByBit,
		constants.OKX,
		constants.Binance,
		constants.Phemex,
		constants.Bitget,
	}, nil)
	exchange.ClearSelected()

	accountName := widget.NewEntry()
	accountName.SetPlaceHolder("")
	accountName.Validator = utils.IsStringEmpty

	publicAPI := widget.NewEntry()
	publicAPI.SetPlaceHolder("XXXX")
	publicAPI.Validator = utils.IsStringEmpty

	secretAPI := widget.NewPasswordEntry()
	secretAPI.SetPlaceHolder("bybit api secret")
	secretAPI.Validator = utils.IsStringEmpty

	passphrase := widget.NewPasswordEntry()
	passphrase.SetPlaceHolder("passphrase")
	passphrase.Validator = utils.IsStringEmpty

	stopIfFallUnder := widget.NewEntry()
	stopIfFallUnder.SetPlaceHolder("0")
	stopIfFallUnder.Validator = utils.IsFloat

	testMode := widget.NewCheck("", func(b bool) {})

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Exchange"),
		exchange,
		widget.NewLabel("Account Name"),
		accountName,
		widget.NewLabel("Public API Key"),
		publicAPI,
		widget.NewLabel("Secret API Key"),
		secretAPI,
		widget.NewLabel("Passphrase"),
		passphrase,
		widget.NewLabel("Stop If Fall Under"),
		stopIfFallUnder,
		widget.NewLabel("Test Mode"),
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

			if valid {
				p := user.Profile{
					Title:       title.Text,
					Exchange:    constants.Exchange(exchange.Selected),
					AccountName: accountName.Text,
					PublicAPI:   publicAPI.Text,
					SecretAPI:   secretAPI.Text,
					Passphrase:  passphrase.Text,
					TestMode:    testMode.Checked,
				}
				p.StopIfFallUnder, _ = strconv.ParseFloat(stopIfFallUnder.Text, 64)

				inserted, err := app.DB.InsertProfile(p)

				if err != nil {
					app.Logger.Error(err)
				} else {
					app.User.ProfileManager.Profiles = append(app.User.ProfileManager.Profiles, *inserted)
				}
				app.refreshProfilesTab()
				app.RefreshProfileSelector()
				app.RefreshLeaderboardWithoutFetch()
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
		if s != pf.Title && app.User.ProfileManager.GetProfileByTitle(s) != nil {
			return fmt.Errorf("a profile named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	exchange := widget.NewSelect([]string{
		constants.ByBit,
		constants.OKX,
		constants.Binance,
		constants.Phemex,
		constants.Bitget,
	}, nil)
	exchange.SetSelected(string(pf.Exchange))

	accountName := widget.NewEntry()
	accountName.SetText(pf.AccountName)
	accountName.Validator = utils.IsStringEmpty

	publicAPI := widget.NewEntry()
	publicAPI.SetText(pf.PublicAPI)
	publicAPI.Validator = utils.IsStringEmpty

	secretAPI := widget.NewPasswordEntry()
	secretAPI.SetText(pf.SecretAPI)
	secretAPI.Validator = utils.IsStringEmpty

	passphrase := widget.NewPasswordEntry()
	passphrase.SetText(pf.Passphrase)
	passphrase.Validator = utils.IsStringEmpty

	stopIfFallUnder := widget.NewEntry()
	stopIfFallUnder.SetText(fmt.Sprintf("%f", pf.StopIfFallUnder))
	stopIfFallUnder.Validator = utils.IsFloat

	testMode := widget.NewCheck("", func(b bool) {})
	testMode.SetChecked(pf.TestMode)

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Exchange"),
		exchange,
		widget.NewLabel("Account Name"),
		accountName,
		widget.NewLabel("Public API Key"),
		publicAPI,
		widget.NewLabel("Secret API Key"),
		secretAPI,
		widget.NewLabel("Passphrase"),
		passphrase,
		widget.NewLabel("Stop If Fall Under"),
		stopIfFallUnder,
		widget.NewLabel("Test Mode"),
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

			if valid {
				p := user.Profile{
					Title:       title.Text,
					Exchange:    constants.Exchange(exchange.Selected),
					AccountName: accountName.Text,
					PublicAPI:   publicAPI.Text,
					SecretAPI:   secretAPI.Text,
					Passphrase:  passphrase.Text,
					TestMode:    testMode.Checked,
				}
				p.StopIfFallUnder, _ = strconv.ParseFloat(stopIfFallUnder.Text, 64)

				err := app.DB.UpdateProfile(pf.ID, p)

				if err != nil {
					app.Logger.Error(err)
				} else {
					app.User.ProfileManager.UpdateProfile(pf.ID, p)
				}
				app.refreshProfilesTab()
				app.RefreshProfileSelector()
				app.RefreshLeaderboardWithoutFetch()
				app.refreshCopiedTradersTab()
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
	slice = append(slice, []any{"Profile Title", "Exchange", "Public API", "Stop If Fall Under", "Test", "Actions"})

	for _, x := range app.User.ProfileManager.Profiles {
		var currentRow []any

		if len(x.Title) > 30 {
			currentRow = append(currentRow, x.Title[:29]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		currentRow = append(currentRow, x.Exchange)

		if len(x.PublicAPI) > 5 {
			currentRow = append(currentRow, fmt.Sprintf("%s****", x.PublicAPI[0:4]))
		} else {
			currentRow = append(currentRow, x.PublicAPI)
		}

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.StopIfFallUnder))

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
			return container.NewMax(lbl, toolbar)
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
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							inserted, err := app.DB.InsertProfile(*pf)
							if err != nil {
								app.Logger.Error(err)
							} else {
								app.User.ProfileManager.Profiles = append(app.User.ProfileManager.Profiles, *inserted)
							}
							app.refreshProfilesTab()
							app.RefreshProfileSelector()
							app.RefreshLeaderboardWithoutFetch()
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf != nil {
							app.editProfileDialog(pf)
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf == nil {
							return
						}

						dialog.ShowConfirm("Delete?", "", func(deleted bool) {
							if deleted {
								err := app.DB.DeleteProfile(pf.ID)
								if err != nil {
									app.Logger.Error(err)
								} else {
									app.User.ProfileManager.DeleteProfile(pf.ID)
								}
							}
							app.refreshProfilesTab()
							app.RefreshProfileSelector()
							app.RefreshLeaderboardWithoutFetch()
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

	btnAdd := widget.NewButtonWithIcon("Add Profile", theme.ContentAddIcon(), func() {
		app.addProfileDialog()
	})
	btnClear := widget.NewButtonWithIcon("Clear Profiles", theme.ContentRemoveIcon(), func() {
		dialog.ShowConfirm(
			"Delete all profiles?",
			fmt.Sprintf("Do you really want to delete %d profiles?", len(app.User.ProfileManager.Profiles)),
			func(deleted bool) {
				if deleted {
					for _, p := range app.User.ProfileManager.Profiles {
						err := app.DB.DeleteProfile(p.ID)
						if err != nil {
							app.Logger.Error(err)
						} else {
							app.User.ProfileManager.DeleteProfile(p.ID)
						}
					}
					app.refreshProfilesTab()
					app.RefreshProfileSelector()
					app.RefreshLeaderboardWithoutFetch()
				}
			}, app.MainWindow)
	})
	btnClear.Importance = widget.DangerImportance

	app.Bottom.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		btnAdd,
		btnClear,
	}

	app.Bottom.Refresh()
}

func (app *Config) refreshProfilesTopContent() {

	txt := widget.NewRichTextFromMarkdown(`## ` + strconv.Itoa(len(app.User.ProfileManager.Profiles)) + ` Profiles Loaded`)
	app.Top.Objects = []fyne.CanvasObject{
		txt,
	}

	app.Top.Refresh()
}

func (app *Config) refreshProfilesTab() {
	app.refreshProfilesTopContent()
	app.refreshProfilesTable()
	app.RefreshProfileSelector()
}
