package nyx

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/resources"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

const (
	defaultPeriodType     = "WEEKLY"
	defaultStatisticsType = "ROI"
)

type LeaderboardTab struct {
	Traders []Trader
	*container.TabItem
	SelectedProfile user.Profile
	CardsContainer  *fyne.Container
}

func (app *Config) traderDialog(t Trader) dialog.Dialog {

	var slice [][]any
	slice = append(slice, []any{"Symbol", "Size", "Entry Price", "Mark Price", "PNL", "Action"})
	slice = append(slice, []any{"BTCUSDT Perpetual", "0.293", "26,996.00", "26,707.95", "84.40 (26.9625%)", t.EncryptedUid})

	traderCard := app.getTraderCard(t, true)
	positionsTable := widget.NewTable(
		func() (int, int) {
			return len(slice), len(slice[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			btn := widget.NewButton("Trade", nil)
			btn.Hide()
			return container.NewMax(lbl, btn)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			btn := container.Objects[1].(*widget.Button)

			if i.Row != 0 && i.Col == 5 {
				lbl.Hide()
				btn.Hidden = false
				btn.OnTapped = func() {
					app.Logger.Debug("trade button tapped")
				}
			} else {
				btn.Hide()
				lbl.Hidden = false
				// we're just putting in textual information
				lbl.SetText(
					slice[i.Row][i.Col].(string))
			}
		})
	colWidths := []float32{200, 100, 100, 100, 200, 60}
	for i, w := range colWidths {
		positionsTable.SetColumnWidth(i, w)
	}
	positionsCard := widget.NewCard("Positions", "", positionsTable)

	grid := container.NewGridWithRows(2, traderCard, positionsCard)
	scrollContent := container.NewVScroll(grid)

	traderDialog := dialog.NewCustom(
		"Trader Overview",
		"Close",
		scrollContent,
		app.MainWindow)

	traderDialog.Resize(fyne.NewSize(900, 600))
	traderDialog.Show()

	return traderDialog
}

func (app *Config) leaderboardTab() *fyne.Container {
	var profileGroups []string

	grid := container.NewAdaptiveGrid(3)
	app.LeaderboardTab.Traders, _ = app.fetchTraders(defaultStatisticsType, defaultPeriodType)
	cards := app.makeTradersCards()
	for _, card := range cards {
		grid.Add(card)
	}
	vScroll := container.NewVScroll(grid)

	for _, pfg := range app.User.ProfileManager.Groups {
		profileGroups = append(profileGroups, pfg.Name)
	}

	profileSelector := widget.NewSelect([]string{}, func(s string) {})
	profileSelector.Disable()
	groupSelector := widget.NewSelect(profileGroups, func(s string) {
		profileSelector.Options = []string{}
		profiles := app.User.ProfileManager.FilterByGroupName(s)
		for _, p := range profiles {
			profileSelector.Options = append(profileSelector.Options, p.Title)
		}
		if len(profileSelector.Options) > 0 {
			profileSelector.Enable()
		}
		profileSelector.Refresh()
	})

	searchEntry := widget.NewSelectEntry([]string{})
	searchEntry.SetPlaceHolder("Search by nickname...")
	searchResults := fyne.NewMenu("")
	onSearchSubmitted := func() {
		go func() {
			searchResults.Items = []*fyne.MenuItem{}
			filtered, err := app.searchByNickname(searchEntry.Text)
			if err != nil {
				searchResults.Items = append(searchResults.Items, fyne.NewMenuItem("API error", func() {}))
			} else {
				var options []string
				for _, trader := range filtered {
					if len(options) > 5 {
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
						app.traderDialog(Trader{
							NickName:      n,
							UserPhotoUrl:  photoUrl,
							EncryptedUid:  uid,
							FollowerCount: followers,
							TwitterUrl:    twitterUrl,
							Pnl:           pnl,
							Roi:           roi,
						})
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

	leftTopContainer := container.NewVBox(widget.NewLabel("Filter and sort"), container.NewHBox(widget.NewLabel("Time"), filterByPeriod, widget.NewLabel("Sort by"), sortByStatistics), searchEntry)
	rightTopContainer := container.NewVBox(widget.NewLabel("Select Profile"), groupSelector, profileSelector)
	topContainer := container.NewAdaptiveGrid(2, leftTopContainer, rightTopContainer)
	releasesContainer := container.NewWithoutLayout(vScroll, topContainer)

	topContainer.Move(fyne.NewPos(10, 10))
	topContainer.Resize(fyne.NewSize(1260, 50))
	vScroll.Move(fyne.NewPos(10, 130))
	vScroll.Resize(fyne.NewSize(1260, 540))

	app.LeaderboardTab.CardsContainer = grid

	return releasesContainer
}

func (app *Config) makeTradersCards() []*widget.Card {
	var cards []*widget.Card

	for _, trader := range app.LeaderboardTab.Traders {

		card := app.getTraderCard(trader, false)
		//card.SetImage(canvasImage)
		cards = append(cards, card)
	}

	return cards
}

func (app *Config) getTraderCard(trader Trader, showImage bool) *widget.Card {
	var twitterLink, binanceLink fyne.CanvasObject
	var canvasImage *canvas.Image
	copyButton := widget.NewButton("Copy", func() {})
	copyButton.Importance = widget.HighImportance

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

	card := widget.NewCard(
		trader.NickName,
		fmt.Sprintf("%d Followers", trader.FollowerCount),
		container.NewGridWithColumns(2,
			widget.NewLabel(fmt.Sprintf("ROI: %.2f%%", trader.Roi*100)), widget.NewLabel(fmt.Sprintf("PNL (USD): %.2f", trader.Pnl)),
			container.NewHBox(binanceLink, twitterLink), copyButton))

	if showImage {
		if utils.DoesFileExist(fmt.Sprintf("%s.jpg", trader.EncryptedUid)) {
			canvasImage = canvas.NewImageFromFile(fmt.Sprintf("%s.jpg", trader.EncryptedUid))
		} else {
			err := app.downloadFile(trader.UserPhotoUrl, trader.EncryptedUid)
			if err != nil {
				// return bundled error image
				canvasImage = canvas.NewImageFromResource(resources.ResourceNoImageAvailablePng)
			} else {
				canvasImage = canvas.NewImageFromFile(fmt.Sprintf("%s.jpg", trader.EncryptedUid))
			}
		}
		canvasImage.SetMinSize(fyne.NewSize(25, 25))
		canvasImage.FillMode = canvas.ImageFillContain
		card.SetImage(canvasImage)
	}

	return card
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
