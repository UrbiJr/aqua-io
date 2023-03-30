package nyx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

type ProfilesTab struct {
	*container.TabItem
	ProfileGroupsList   *widget.List
	Top                 *fyne.Container
	ProfilesTable       *widget.Table
	ProfilesSlice       [][]any
	Bottom              *fyne.Container
	CurrentProfileGroup int64
}

// NewProfilesView returns a view for the profiles management
func (app *Config) profilesTab() *fyne.Container {
	// define profilesTab
	app.ProfilesTab.Top = container.NewMax()
	app.ProfilesTab.Bottom = container.NewHBox()

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
		app.refreshProfilesTopContent()
		app.refreshProfilesBottomContent()
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
	app.refreshProfilesTopContent()
	app.refreshProfilesBottomContent()

	// get the profiles table
	app.ProfilesTable = app.getProfilesTable()

	// define the center container
	vScroll := container.NewScroll(app.ProfilesTable)

	// define the profilesTab container
	profilesTabContainer := container.NewWithoutLayout(leftContainer, app.Top, vScroll, app.Bottom)

	// resize and move profilesTab elements
	leftContainer.Move(fyne.NewPos(10, 10))
	leftContainer.Resize(fyne.NewSize(280, 780))

	vScroll.Move(fyne.NewPos(300, 90))
	vScroll.Resize(fyne.NewSize(1120, 660))

	app.Top.Move(fyne.NewPos(300, 10))
	app.Top.Resize(fyne.NewSize(1120, 64))

	app.Bottom.Move(fyne.NewPos(300, 760))
	app.Bottom.Resize(fyne.NewSize(1120, 50))

	return profilesTabContainer
}

func (app *Config) addProfileGroupDialog() dialog.Dialog {
	nameEntry := widget.NewEntry()

	nameEntry.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		if app.User.ProfileManager.GetGroupByName(s) != nil {
			return fmt.Errorf("a group named %s already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a name")
		} else {
			return nil
		}
	}

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
				app.refreshProfilesTopContent()
				app.refreshProfilesBottomContent()
				app.refreshProfilesTable()
				app.refreshProfileSelector()
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
					app.refreshProfilesTopContent()
					app.refreshProfilesBottomContent()
					app.refreshProfilesTable()
					app.refreshProfileSelector()
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

	bybitApiKey := widget.NewEntry()
	bybitApiKey.SetPlaceHolder("XXXX")
	bybitApiKey.Validator = utils.IsStringEmpty

	bybitApiSecret := widget.NewPasswordEntry()
	bybitApiSecret.SetPlaceHolder("bybit api secret")
	bybitApiSecret.Validator = utils.IsStringEmpty

	maxBybitBinancePriceDifferentPercent := widget.NewEntry()
	maxBybitBinancePriceDifferentPercent.SetPlaceHolder("0.5")
	maxBybitBinancePriceDifferentPercent.Validator = utils.IsFloat

	leverage := widget.NewEntry()
	leverage.SetPlaceHolder("15")
	leverage.Validator = utils.IsFloat

	initialOpenPercent := widget.NewEntry()
	initialOpenPercent.SetPlaceHolder("2.5")
	initialOpenPercent.Validator = utils.IsFloat

	maxAddMultiplier := widget.NewEntry()
	maxAddMultiplier.SetPlaceHolder("5")
	maxAddMultiplier.Validator = utils.IsFloat

	openDelay := widget.NewEntry()
	openDelay.SetPlaceHolder("5")
	openDelay.Validator = utils.IsFloat

	oneCoinMaxPercent := widget.NewEntry()
	oneCoinMaxPercent.SetPlaceHolder("25")
	oneCoinMaxPercent.Validator = utils.IsFloat

	blackListCoins := widget.NewMultiLineEntry()
	blackListCoins.SetPlaceHolder("coin1,coin2,coin3")

	addPreventionPercent := widget.NewEntry()
	addPreventionPercent.SetPlaceHolder("0.5")
	addPreventionPercent.Validator = utils.IsFloat

	blockAddsAboveEntry := widget.NewCheck("", func(b bool) {})

	maxOpenPositions := widget.NewEntry()
	maxOpenPositions.SetPlaceHolder("10")
	maxOpenPositions.Validator = utils.IsInteger

	autoTP := widget.NewEntry()
	autoTP.SetPlaceHolder("7")
	autoTP.Validator = utils.IsFloat

	autoSL := widget.NewEntry()
	autoSL.SetPlaceHolder("0")
	autoSL.Validator = utils.IsFloat

	testMode := widget.NewCheck("", func(b bool) {})

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Bybit Api Key"),
		bybitApiKey,
		widget.NewLabel("Bybit API Secret"),
		bybitApiSecret,
		widget.NewLabel("Binance/ByBit Price Difference"),
		maxBybitBinancePriceDifferentPercent,
		widget.NewLabel("Leverage"),
		leverage,
		widget.NewLabel("Initial Open Percent"),
		initialOpenPercent,
		widget.NewLabel("Max Add Multiplier"),
		maxAddMultiplier,
		widget.NewLabel("Open Delay"),
		openDelay,
		widget.NewLabel("One Coin Max Percent*"),
		oneCoinMaxPercent,
		widget.NewLabel("Blacklist Coins"),
		blackListCoins,
		widget.NewLabel("Add Prevention Percent"),
		addPreventionPercent,
		widget.NewLabel("Blocks Adds Above Entry"),
		blockAddsAboveEntry,
		widget.NewLabel("Max Open Positions"),
		maxOpenPositions,
		widget.NewLabel("Auto TP"),
		autoTP,
		widget.NewLabel("Auto SL"),
		autoSL,
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
				app.Logger.Debug("Form submitted")
				p := user.Profile{
					GroupID:             app.CurrentProfileGroup,
					Title:               title.Text,
					BybitApiKey:         bybitApiKey.Text,
					BybitApiSecret:      bybitApiSecret.Text,
					BlacklistCoins:      strings.Split(blackListCoins.Text, ","),
					BlockAddsAboveEntry: blockAddsAboveEntry.Checked,
					TestMode:            testMode.Checked,
				}
				p.MaxBybitBinancePriceDifferentPercent, _ = strconv.ParseFloat(maxBybitBinancePriceDifferentPercent.Text, 64)
				p.Leverage, _ = strconv.ParseFloat(leverage.Text, 64)
				p.MaxAddMultiplier, _ = strconv.ParseFloat(maxAddMultiplier.Text, 64)
				p.OpenDelay, _ = strconv.ParseFloat(openDelay.Text, 64)
				p.OneCoinMaxPercent, _ = strconv.ParseFloat(oneCoinMaxPercent.Text, 64)
				p.AddPreventionPercent, _ = strconv.ParseFloat(addPreventionPercent.Text, 64)
				p.AutoTP, _ = strconv.ParseFloat(autoTP.Text, 64)
				p.AutoSL, _ = strconv.ParseFloat(autoSL.Text, 64)
				p.MaxOpenPositions, _ = strconv.ParseInt(maxOpenPositions.Text, 10, 64)

				_, err := app.DB.InsertProfile(p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfileGroupsList()
				app.refreshProfilesTopContent()
				app.refreshProfilesTable()
				app.refreshProfileSelector()
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

	bybitApiKey := widget.NewEntry()
	bybitApiKey.SetText(pf.BybitApiKey)
	bybitApiKey.Validator = utils.IsStringEmpty

	bybitApiSecret := widget.NewPasswordEntry()
	bybitApiSecret.SetText(pf.BybitApiSecret)
	bybitApiSecret.Validator = utils.IsStringEmpty

	maxBybitBinancePriceDifferentPercent := widget.NewEntry()
	maxBybitBinancePriceDifferentPercent.SetText(fmt.Sprintf("%.2f", pf.MaxBybitBinancePriceDifferentPercent))
	maxBybitBinancePriceDifferentPercent.Validator = utils.IsFloat

	leverage := widget.NewEntry()
	leverage.SetText(fmt.Sprintf("%.2f", pf.Leverage))
	leverage.Validator = utils.IsFloat

	initialOpenPercent := widget.NewEntry()
	initialOpenPercent.SetText(fmt.Sprintf("%.2f", pf.InitialOpenPercent))
	initialOpenPercent.Validator = utils.IsFloat

	maxAddMultiplier := widget.NewEntry()
	maxAddMultiplier.SetText(fmt.Sprintf("%.2f", pf.MaxAddMultiplier))
	maxAddMultiplier.Validator = utils.IsFloat

	openDelay := widget.NewEntry()
	openDelay.SetText(fmt.Sprintf("%.2f", pf.OpenDelay))
	openDelay.Validator = utils.IsFloat

	oneCoinMaxPercent := widget.NewEntry()
	oneCoinMaxPercent.SetText(fmt.Sprintf("%.2f", pf.OneCoinMaxPercent))
	oneCoinMaxPercent.Validator = utils.IsFloat

	blackListCoins := widget.NewMultiLineEntry()
	blackListCoins.SetPlaceHolder(strings.Join(pf.BlacklistCoins, ","))

	addPreventionPercent := widget.NewEntry()
	addPreventionPercent.SetText(fmt.Sprintf("%.2f", pf.AddPreventionPercent))
	addPreventionPercent.Validator = utils.IsFloat

	blockAddsAboveEntry := widget.NewCheck("", func(b bool) {})
	blockAddsAboveEntry.SetChecked(pf.BlockAddsAboveEntry)

	maxOpenPositions := widget.NewEntry()
	maxOpenPositions.SetText(fmt.Sprintf("%d", pf.MaxOpenPositions))
	maxOpenPositions.Validator = utils.IsInteger

	autoTP := widget.NewEntry()
	autoTP.SetText(fmt.Sprintf("%.2f", pf.AutoTP))
	autoTP.Validator = utils.IsFloat

	autoSL := widget.NewEntry()
	autoSL.SetText(fmt.Sprintf("%.2f", pf.AutoSL))
	autoSL.Validator = utils.IsFloat

	testMode := widget.NewCheck("", func(b bool) {})
	testMode.SetChecked(pf.TestMode)

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Bybit Api Key"),
		bybitApiKey,
		widget.NewLabel("Bybit API Secret"),
		bybitApiSecret,
		widget.NewLabel("Binance/ByBit Price Difference"),
		maxBybitBinancePriceDifferentPercent,
		widget.NewLabel("Leverage"),
		leverage,
		widget.NewLabel("Initial Open Percent"),
		initialOpenPercent,
		widget.NewLabel("Max Add Multiplier"),
		maxAddMultiplier,
		widget.NewLabel("Open Delay"),
		openDelay,
		widget.NewLabel("One Coin Max Percent*"),
		oneCoinMaxPercent,
		widget.NewLabel("Blacklist Coins"),
		blackListCoins,
		widget.NewLabel("Add Prevention Percent"),
		addPreventionPercent,
		widget.NewLabel("Blocks Adds Above Entry"),
		blockAddsAboveEntry,
		widget.NewLabel("Max Open Positions"),
		maxOpenPositions,
		widget.NewLabel("Auto TP"),
		autoTP,
		widget.NewLabel("Auto SL"),
		autoSL,
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
				app.Logger.Debug("Form submitted")
				p := user.Profile{
					GroupID:             app.CurrentProfileGroup,
					Title:               title.Text,
					BybitApiKey:         bybitApiKey.Text,
					BybitApiSecret:      bybitApiSecret.Text,
					BlacklistCoins:      strings.Split(blackListCoins.Text, ","),
					BlockAddsAboveEntry: blockAddsAboveEntry.Checked,
					TestMode:            testMode.Checked,
				}
				p.MaxBybitBinancePriceDifferentPercent, _ = strconv.ParseFloat(maxBybitBinancePriceDifferentPercent.Text, 64)
				p.Leverage, _ = strconv.ParseFloat(leverage.Text, 64)
				p.MaxAddMultiplier, _ = strconv.ParseFloat(maxAddMultiplier.Text, 64)
				p.OpenDelay, _ = strconv.ParseFloat(openDelay.Text, 64)
				p.OneCoinMaxPercent, _ = strconv.ParseFloat(oneCoinMaxPercent.Text, 64)
				p.AddPreventionPercent, _ = strconv.ParseFloat(addPreventionPercent.Text, 64)
				p.AutoTP, _ = strconv.ParseFloat(autoTP.Text, 64)
				p.AutoSL, _ = strconv.ParseFloat(autoSL.Text, 64)
				p.MaxOpenPositions, _ = strconv.ParseInt(maxOpenPositions.Text, 10, 64)

				err := app.DB.UpdateProfile(pf.ID, p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfileGroupsList()
				app.refreshProfilesTopContent()
				app.refreshProfilesTable()
				app.refreshProfileSelector()
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

	slice = append(slice, []any{"Profile", "ByBit API Key", "Auto TP/SL", "Test", "Actions"})

	for _, x := range app.User.ProfileManager.FilterByGroupID(app.CurrentProfileGroup) {
		var currentRow []any

		if len(x.Title) > 16 {
			currentRow = append(currentRow, x.Title[:12]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		if len(x.BybitApiKey) > 5 {
			currentRow = append(currentRow, fmt.Sprintf("%s****", x.BybitApiKey[0:4]))
		} else {
			currentRow = append(currentRow, x.BybitApiKey)
		}

		currentRow = append(currentRow, fmt.Sprintf("%.2f / %.2f", x.AutoTP, x.AutoSL))

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

			if i.Row != 0 && i.Col == 4 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][4].(int64), app.CurrentProfileGroup)
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							_, err := app.DB.InsertProfile(*pf)
							if err != nil {
								app.Logger.Error(err)
							}
							app.refreshProfileGroupsList()
							app.refreshProfilesTopContent()
							app.refreshProfilesTable()
							app.refreshProfileSelector()
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][4].(int64), app.CurrentProfileGroup)
						if pf != nil {
							app.editProfileDialog(pf)
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
						dialog.ShowConfirm("Delete?", "", func(deleted bool) {
							if deleted {
								pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][4].(int64), app.CurrentProfileGroup)
								err := app.DB.DeleteProfile(pf.ID)
								if err != nil {
									app.Logger.Error(err)
								}
							}
							app.refreshProfileGroupsList()
							app.refreshProfilesTopContent()
							app.refreshProfilesTable()
							app.refreshProfileSelector()
						}, app.MainWindow)
					}))
				}

			} else if i.Col == 3 && i.Row != 0 {
				toolbar.Hide()
				lbl.Hidden = false
				if app.ProfilesSlice[i.Row][i.Col].(bool) == true {
					lbl.SetText("Test")
				} else {
					lbl.SetText("Real")
				}
			} else {
				toolbar.Hide()
				lbl.Hidden = false
				// we're just putting in textual information
				lbl.SetText(
					app.ProfilesSlice[i.Row][i.Col].(string))
			}
		})

	colWidths := []float32{100, 200, 200, 200, 60}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (app *Config) refreshProfilesTable() {
	app.ProfilesSlice = app.getProfilesSlice()
	app.ProfilesTable.Refresh()

	colWidths := []float32{100, 200, 200, 200, 60}
	for i, w := range colWidths {
		app.ProfilesTable.SetColumnWidth(i, w)
	}

}

func (app *Config) refreshProfilesBottomContent() {

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
						app.refreshProfilesTopContent()
						app.refreshProfilesBottomContent()
						app.refreshProfilesTable()
						app.refreshProfileSelector()
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

func (app *Config) refreshProfilesTopContent() {

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
