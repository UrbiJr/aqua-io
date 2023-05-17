package copy_io

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"net/url"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	x_widget "fyne.io/x/fyne/widget"
	"github.com/UrbiJr/aqua-io/internal/resources"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"
)

const (
	defaultPeriodType     = "WEEKLY"
	defaultStatisticsType = "ROI"
)

type LeaderboardTab struct {
	Traders []user.Trader
	*container.TabItem
	ProfileSelector      *widget.Select
	SelectedProfile      *user.Profile
	TraderPositionsSlice [][]any
	CardsContainer       *fyne.Container
}

func (app *Config) leaderboardTab() *fyne.Container {

	grid := container.NewAdaptiveGrid(4)
	// fetch the traders from binance
	app.LeaderboardTab.Traders, _ = app.fetchTraders(defaultStatisticsType, defaultPeriodType)
	// get the leaderboard
	cards := app.makeTradersCards()
	for _, card := range cards {
		grid.Add(card)
	}
	// make it scrollable
	vScroll := container.NewVScroll(grid)

	// get the profile selector
	app.LeaderboardTab.ProfileSelector = widget.NewSelect(app.User.ProfileManager.GetAllTitles(), func(s string) {
		app.LeaderboardTab.SelectedProfile = app.User.ProfileManager.GetProfileByTitle(s)
		app.RefreshLeaderboardWithoutFetch()
	})

	// get the search filter
	searchEntry := widget.NewSelectEntry([]string{})
	searchEntry.SetPlaceHolder("Search by nickname...")
	searchResults := fyne.NewMenu("Search results")
	onSearchSubmitted := func() {
		go func() {
			searchResults.Items = []*fyne.MenuItem{}
			filtered, err := app.searchByNickname(searchEntry.Text)
			if err != nil {
				searchResults.Items = append(searchResults.Items, fyne.NewMenuItem("API error", func() {}))
			} else {
				for _, trader := range filtered {
					if len(searchResults.Items) > 5 {
						break
					}

					n := trader.NickName
					photoUrl := trader.UserPhotoUrl
					uid := trader.EncryptedUid
					followers := trader.FollowerCount
					twitterUrl := trader.TwitterUrl
					pnl := trader.Pnl
					roi := trader.Roi
					searchResults.Items = append(searchResults.Items, fyne.NewMenuItem(n, func() {
						app.Logger.Debug("Selected trader " + n)
						go func() {
							app.traderDialog(user.Trader{
								NickName:      n,
								UserPhotoUrl:  photoUrl,
								EncryptedUid:  uid,
								FollowerCount: followers,
								TwitterUrl:    twitterUrl,
								Pnl:           pnl,
								Roi:           roi,
							}, "")
						}()

					}))
				}
			}
			entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(searchEntry)
			popUp := widget.NewPopUpMenu(searchResults, app.App.Driver().CanvasForObject(searchEntry))
			popUp.ShowAtPosition(entryPos.Add(fyne.NewPos(0, searchEntry.Size().Height)))
			popUp.Resize(fyne.NewSize(searchEntry.Size().Width, 50))
		}()
	}
	searchEntry.OnSubmitted = func(s string) {
		onSearchSubmitted()
	}
	btn := widget.NewButton("", onSearchSubmitted)
	btn.Importance = widget.LowImportance
	btn.SetIcon(theme.SearchIcon())
	searchEntry.ActionItem = btn

	// get the other filters
	sortByStatistics := widget.NewSelect([]string{"ROI", "PNL"}, nil)
	sortByStatistics.SetSelected("ROI")
	filterByPeriod := widget.NewSelect([]string{"DAILY", "WEEKLY", "MONTHLY", "TOTAL"}, nil)
	filterByPeriod.SetSelected("WEEKLY")
	sortByStatistics.OnChanged = func(s string) {
		go func() {
			app.RefreshLeaderboard(sortByStatistics.Selected, filterByPeriod.Selected)
		}()

	}
	filterByPeriod.OnChanged = func(s string) {
		go func() {
			app.RefreshLeaderboard(sortByStatistics.Selected, filterByPeriod.Selected)
		}()
	}

	// display and resize content
	leftTopContainer := container.NewVBox(container.NewHBox(widget.NewLabel("Time"), filterByPeriod, widget.NewLabel("Sort by"), sortByStatistics), searchEntry)
	rightTopContainer := container.NewVBox(widget.NewLabel("Select Profile"), app.LeaderboardTab.ProfileSelector)
	topContainer := container.NewAdaptiveGrid(2, leftTopContainer, rightTopContainer)
	mainContainer := container.NewWithoutLayout(vScroll, topContainer)

	topContainer.Move(fyne.NewPos(10, 10))
	topContainer.Resize(fyne.NewSize(1260, 50))
	vScroll.Move(fyne.NewPos(0, 90))
	vScroll.Resize(fyne.NewSize(1310, 540))

	app.LeaderboardTab.CardsContainer = grid

	return mainContainer
}

func (app *Config) getTraderPositionsSlice(t user.Trader) [][]any {
	var slice [][]any

	slice = append(slice, []any{"Symbol", "Size", "Entry Price", "Mark Price", "PNL"})
	positions, err := app.fetchTraderPositions(t.EncryptedUid)
	if err != nil {
		return slice
	}

	for _, x := range positions {
		var currentRow []any

		currentRow = append(currentRow, fmt.Sprintf("%s Perpetual", x.Symbol))

		currentRow = append(currentRow, x.Amount)

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.EntryPrice))

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.MarkPrice))

		currentRow = append(currentRow, x.Pnl)

		slice = append(slice, currentRow)
	}

	return slice
}

// traderDialog represents the "Trader Overview" which allows user to copy the trader and view its positions
func (app *Config) traderDialog(t user.Trader, traderID string) dialog.Dialog {

	if t.EncryptedUid == "" {
		if traderID != "" {
			// fetch trader
			traderInfo, err := app.fetchTraderByUid(traderID)
			if err != nil {
				return nil
			}
			t = *traderInfo
		} else {
			return nil
		}
	}

	// get the trader card (nickname + image + copy button)
	traderCard, gif := app.getTraderCard(t, true, false)

	// get the positions table
	var slice [][]any
	slice = append(slice, []any{"Symbol", "Size", "Entry Price", "Mark Price", "PNL"})
	app.TraderPositionsSlice = slice
	positionsTable := widget.NewTable(
		func() (int, int) {
			return len(app.TraderPositionsSlice), len(app.TraderPositionsSlice[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			canvasText := canvas.NewText("", nil)
			canvasText.Hide()
			return container.NewMax(lbl, canvasText)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			canvasText := container.Objects[1].(*canvas.Text)

			if i.Col == 4 && i.Row != 0 {
				lbl.Hide()
				canvasText.Hidden = false
				pnl := app.TraderPositionsSlice[i.Row][i.Col].(float64)
				if pnl > 0 {
					canvasText.Color = color.RGBA{R: 14, G: 203, B: 129, A: 255}
				} else {
					canvasText.Color = color.RGBA{R: 246, G: 70, B: 93, A: 255}
				}
				canvasText.Text = fmt.Sprintf("%.2f", pnl)
			} else if i.Col == 0 && i.Row != 0 {
				lbl.Hide()
				canvasText.Hidden = false
				amount := app.TraderPositionsSlice[i.Row][1].(float64)
				if amount > 0 {
					canvasText.Color = color.RGBA{R: 14, G: 203, B: 129, A: 255}
					canvasText.Text = fmt.Sprintf("%s (Long)", app.TraderPositionsSlice[i.Row][i.Col].(string))
				} else {
					canvasText.Color = color.RGBA{R: 246, G: 70, B: 93, A: 255}
					canvasText.Text = fmt.Sprintf("%s (Short)", app.TraderPositionsSlice[i.Row][i.Col].(string))
				}
			} else if i.Col == 1 && i.Row != 0 {
				canvasText.Hide()
				lbl.Hidden = false
				lbl.SetText(
					fmt.Sprintf("%.3f", math.Abs(app.TraderPositionsSlice[i.Row][i.Col].(float64))))
			} else {
				canvasText.Hide()
				lbl.Hidden = false
				lbl.SetText(
					app.TraderPositionsSlice[i.Row][i.Col].(string))
			}
		})
	colWidths := []float32{200, 100, 100, 100, 200}
	for i, w := range colWidths {
		positionsTable.SetColumnWidth(i, w)
	}
	positionsCard := widget.NewCard("Positions", "", positionsTable)
	go func() {
		app.TraderPositionsSlice = app.getTraderPositionsSlice(t)
		positionsTable.Refresh()
	}()

	var grid *fyne.Container
	if gif != nil {
		grid = container.NewGridWithRows(3, container.NewMax(gif), traderCard, positionsCard)
		gif.Start()
	} else {
		grid = container.NewGridWithRows(2, traderCard, positionsCard)
	}

	scrollContent := container.NewVScroll(grid)

	// merge everything
	traderDialog := dialog.NewCustom(
		"Trader Overview",
		"Close",
		scrollContent,
		app.MainWindow)

	traderDialog.Resize(fyne.NewSize(900, 600))
	traderDialog.Show()

	return traderDialog
}

// makeTradersCards returns a list of cards to be displayed in the leaderboard
func (app *Config) makeTradersCards() []*widget.Card {
	var cards []*widget.Card

	for _, trader := range app.LeaderboardTab.Traders {
		card, _ := app.getTraderCard(trader, false, true)
		cards = append(cards, card)
	}

	return cards
}

func (app *Config) noProfileSelectedDialog() dialog.Dialog {
	d := dialog.NewInformation(
		"No Profile Selected",
		"Please select a profile first, then try again.",
		app.MainWindow)
	d.Show()

	return d
}

func (app *Config) copyTraderDialog(t user.Trader) dialog.Dialog {
	var d dialog.Dialog

	if app.LeaderboardTab.SelectedProfile.TraderID != "" {
		d = dialog.NewError(
			errors.New("Selected profile already has a trader.\nPlease select a different profile, then try again"),
			app.MainWindow,
		)
	} else {
		d = dialog.NewCustomConfirm(
			"Copy?",
			"Confirm",
			"Cancel",
			widget.NewLabelWithStyle(fmt.Sprintf("Copy %s positions?", t.NickName), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			func(b bool) {
				if b {
					app.copyTrader(t, app.LeaderboardTab.SelectedProfile)
				}
			},
			app.MainWindow)
	}

	d.Show()
	return d
}

func (app *Config) stopCopyingTraderDialog(t user.Trader) dialog.Dialog {
	d := dialog.NewCustomConfirm(
		"Stop Copying?",
		"Confirm",
		"Cancel",
		widget.NewLabelWithStyle(fmt.Sprintf("Stop Copying %s positions?", t.NickName), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		func(b bool) {
			if b {
				app.stopCopyingTrader(t, "")
			}
		},
		app.MainWindow)

	d.Show()

	return d
}

func (app *Config) makeOrderDialog(p *user.Profile, symbol, orderType string, price float64) dialog.Dialog {

	amount := widget.NewEntry()
	amount.SetPlaceHolder("0")
	amount.Validator = utils.IsFloat

	d := dialog.NewForm(
		fmt.Sprintf("Order %s", symbol),
		"Confirm",
		"Skip",
		[]*widget.FormItem{
			widget.NewFormItem("Amount (USDT): ", amount),
		},
		func(b bool) {
			if b {
				a, err := strconv.ParseFloat(amount.Text, 64)
				if err != nil {
					app.Logger.Error(err)
					return
				}
				_, err = app.createOrder(p, symbol, orderType, a, price)
				if err != nil {
					app.App.SendNotification(&fyne.Notification{
						Title:   "❌ Error Creating Order",
						Content: err.Error(),
					})
					app.Logger.Error(fmt.Sprintf("failed to create %s order: %s", symbol, err.Error()))
				} else {
					app.App.SendNotification(&fyne.Notification{
						Title:   "🤑 Success!",
						Content: fmt.Sprintf("Successfully created %s order", symbol),
					})
					app.Logger.Debug(fmt.Sprintf("successfully created %s order", symbol))
				}
			}
		},
		app.MainWindow)

	d.Show()

	return d
}

// getTraderCard returns a card widget containing information about the trader and relative actions,
// it eventually also returns a pointer to a GIF widget if the trader image is a GIF (as this cannot be displayed in the card)
func (app *Config) getTraderCard(trader user.Trader, showImage bool, showPopUpButton bool) (*widget.Card, *x_widget.AnimatedGif) {
	var twitterLink, binanceLink fyne.CanvasObject
	var canvasImage *canvas.Image
	var btn *widget.Button

	if app.LeaderboardTab.SelectedProfile != nil && app.LeaderboardTab.SelectedProfile.TraderID == trader.EncryptedUid {
		btn = widget.NewButton("Stop Copying", func() {
			app.stopCopyingTraderDialog(trader)
		})
	} else {
		btn = widget.NewButton("Copy", func() {
			if app.LeaderboardTab.SelectedProfile == nil {
				app.noProfileSelectedDialog()
				return
			}
			app.copyTraderDialog(trader)
		})
	}

	btn.Importance = widget.HighImportance

	if trader.TwitterUrl == nil {
		twitterLink = widget.NewLabel("")
	} else {
		twitterUrl, _ := url.Parse(trader.TwitterUrl.(string))
		twitterLink = widget.NewHyperlink("Twitter", twitterUrl)
	}

	if trader.EncryptedUid == "" {
		binanceLink = widget.NewLabel("")
	} else {
		binanceUrl, _ := url.Parse("https://www.binance.com/en/futures-activity/leaderboard/user/um?encryptedUid=" + trader.EncryptedUid)
		binanceLink = widget.NewHyperlink("Binance", binanceUrl)
	}

	var card *widget.Card
	if showPopUpButton {
		card = widget.NewCard(
			trader.NickName,
			fmt.Sprintf("%d Followers", trader.FollowerCount),
			container.NewGridWithColumns(2,
				widget.NewButtonWithIcon("View Positions", theme.VisibilityIcon(), func() {
					go func() {
						app.traderDialog(trader, "")
					}()
				}), widget.NewLabel(""),
				widget.NewLabel(fmt.Sprintf("ROI: %.2f%%", trader.Roi*100)), widget.NewLabel(fmt.Sprintf("PNL (USD): %.2f", trader.Pnl)),
				container.NewHBox(binanceLink, twitterLink), btn))
	} else {
		card = widget.NewCard(
			trader.NickName,
			fmt.Sprintf("%d Followers", trader.FollowerCount),
			container.NewGridWithColumns(2,
				widget.NewLabel(fmt.Sprintf("ROI: %.2f%%", trader.Roi*100)), widget.NewLabel(fmt.Sprintf("PNL (USD): %.2f", trader.Pnl)),
				container.NewHBox(binanceLink, twitterLink), btn))
	}

	if showImage {
		ext := filepath.Ext(trader.UserPhotoUrl)
		if !utils.DoesFileExist(fmt.Sprintf("downloads/%s%s", trader.EncryptedUid, ext)) {
			// image not stored locally, download it
			err := app.downloadFile(trader.UserPhotoUrl, trader.EncryptedUid, ext)
			if err != nil {
				// return bundled error image
				canvasImage = canvas.NewImageFromResource(resources.ResourceNoImageAvailablePng)
				canvasImage.SetMinSize(fyne.NewSize(25, 25))
				canvasImage.FillMode = canvas.ImageFillContain
				card.SetImage(canvasImage)
				return card, nil
			}
		}
		switch ext {
		case ".jpg", ".png":
			canvasImage = canvas.NewImageFromFile(fmt.Sprintf("downloads/%s%s", trader.EncryptedUid, ext))
			canvasImage.SetMinSize(fyne.NewSize(25, 25))
			canvasImage.FillMode = canvas.ImageFillContain
			card.SetImage(canvasImage)
		case ".gif":
			gif, err := x_widget.NewAnimatedGif(storage.NewFileURI(fmt.Sprintf("downloads/%s%s", trader.EncryptedUid, ext)))
			if err != nil {
				app.Logger.Error(err)
			} else {
				return card, gif
			}
		}

	}

	return card, nil
}

func (app *Config) RefreshLeaderboard(statisticsType, periodType string) {
	app.LeaderboardTab.Traders, _ = app.fetchTraders(statisticsType, periodType)
	cards := app.makeTradersCards()
	app.LeaderboardTab.CardsContainer.RemoveAll()
	for _, card := range cards {
		app.LeaderboardTab.CardsContainer.Add(card)
	}
	app.LeaderboardTab.CardsContainer.Refresh()
}

func (app *Config) RefreshLeaderboardWithoutFetch() {
	cards := app.makeTradersCards()
	app.LeaderboardTab.CardsContainer.RemoveAll()
	for _, card := range cards {
		app.LeaderboardTab.CardsContainer.Add(card)
	}
	app.LeaderboardTab.CardsContainer.Refresh()
}

func (app *Config) RefreshProfileSelector() {
	app.LeaderboardTab.ProfileSelector.ClearSelected()
	app.LeaderboardTab.ProfileSelector.Refresh()
}
