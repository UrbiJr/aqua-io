package copy_io

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"
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
	leverage.Validator = utils.IsInteger

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
				p := user.Profile{
					Title:               title.Text,
					BybitApiKey:         bybitApiKey.Text,
					BybitApiSecret:      bybitApiSecret.Text,
					BlacklistCoins:      strings.Split(blackListCoins.Text, ","),
					BlockAddsAboveEntry: blockAddsAboveEntry.Checked,
					TestMode:            testMode.Checked,
				}
				p.MaxBybitBinancePriceDifferentPercent, _ = strconv.ParseFloat(maxBybitBinancePriceDifferentPercent.Text, 64)
				p.Leverage, _ = strconv.ParseInt(leverage.Text, 10, 64)
				p.MaxAddMultiplier, _ = strconv.ParseFloat(maxAddMultiplier.Text, 64)
				p.OpenDelay, _ = strconv.ParseFloat(openDelay.Text, 64)
				p.OneCoinMaxPercent, _ = strconv.ParseFloat(oneCoinMaxPercent.Text, 64)
				p.AddPreventionPercent, _ = strconv.ParseFloat(addPreventionPercent.Text, 64)
				p.AutoTP, _ = strconv.ParseFloat(autoTP.Text, 64)
				p.AutoSL, _ = strconv.ParseFloat(autoSL.Text, 64)
				p.MaxOpenPositions, _ = strconv.ParseInt(maxOpenPositions.Text, 10, 64)

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
	leverage.SetText(fmt.Sprintf("%d", pf.Leverage))
	leverage.Validator = utils.IsInteger

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
				p := user.Profile{
					Title:               title.Text,
					TraderID:            pf.TraderID,
					BybitApiKey:         bybitApiKey.Text,
					BybitApiSecret:      bybitApiSecret.Text,
					BlacklistCoins:      strings.Split(blackListCoins.Text, ","),
					BlockAddsAboveEntry: blockAddsAboveEntry.Checked,
					TestMode:            testMode.Checked,
				}
				p.MaxBybitBinancePriceDifferentPercent, _ = strconv.ParseFloat(maxBybitBinancePriceDifferentPercent.Text, 64)
				p.Leverage, _ = strconv.ParseInt(leverage.Text, 10, 64)
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
				} else {
					app.User.ProfileManager.UpdateProfile(pf.ID, p)
				}
				app.refreshProfilesTab()
				app.RefreshProfileSelector()
				app.RefreshLeaderboardWithoutFetch()
				app.refreshCopiedTradersTab(true)
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
	slice = append(slice, []any{"Profile Title", "Trader", "ByBit API Key", "Auto TP/SL", "Test", "Actions"})

	for _, x := range app.User.ProfileManager.Profiles {
		var currentRow []any

		if len(x.Title) > 16 {
			currentRow = append(currentRow, x.Title[:12]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		if x.TraderID != "" {
			currentRow = append(currentRow, x.TraderID)
		} else {
			currentRow = append(currentRow, "Unset")
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

			if i.Row != 0 && i.Col == 5 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
						pf := app.User.ProfileManager.GetProfileByID(app.ProfilesSlice[i.Row][5].(int64))
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							pf.TraderID = ""
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
						if pf.TraderID != "" {
							dialog.ShowError(
								fmt.Errorf("You must first stop copying trader %s.\nYou can do that from Copied Traders tab.", pf.TraderID),
								app.MainWindow)
						} else {
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
						}

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

	colWidths := []float32{120, 270, 200, 200, 200, 40}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (app *Config) refreshProfilesTable() {
	app.ProfilesSlice = app.getProfilesSlice()
	app.ProfilesTable.Refresh()

	colWidths := []float32{120, 270, 200, 200, 200, 40}
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
			fmt.Sprintf("Do you really want to delete %d profiles?\nNote that profiles with a trader will NOT be deleted: you have to stop copying the trader first.", len(app.User.ProfileManager.GetProfilesWithTrader())),
			func(deleted bool) {
				if deleted {
					for _, p := range app.User.ProfileManager.Profiles {
						if p.TraderID == "" {
							err := app.DB.DeleteProfile(p.ID)
							if err != nil {
								app.Logger.Error(err)
							} else {
								app.User.ProfileManager.DeleteProfile(p.ID)
							}
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
