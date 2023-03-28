package nyx

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/resources"
	"github.com/UrbiJr/nyx/internal/utils"
)

const (
	defaultPeriodType     = "WEEKLY"
	defaultStatisticsType = "ROI"
)

type LeaderboardTab struct {
	Traders []Trader
	*container.TabItem
	CardsContainer *fyne.Container
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
					//uid := trader.EncryptedUid
					searchResults.Items = append(searchResults.Items, fyne.NewMenuItem(n, func() {
						app.Logger.Debug("Selected trader " + n)
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
	var canvasImage *canvas.Image

	for i, trader := range app.LeaderboardTab.Traders {
		var twitterLink, binanceLink fyne.CanvasObject
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

		copyButton := widget.NewButton("Copy", func() {})
		copyButton.Importance = widget.HighImportance

		card := widget.NewCard(
			trader.NickName,
			fmt.Sprintf("%d Followers", trader.FolloweCount),
			container.NewGridWithColumns(2,
				widget.NewLabel(fmt.Sprintf("ROI: %.2f%%", trader.Roi*100)), widget.NewLabel(fmt.Sprintf("PNL (USD): %.2f", trader.Pnl)),
				container.NewHBox(binanceLink, twitterLink), copyButton))

		if utils.DoesFileExist(fmt.Sprintf("%d-trader.jpg", i)) {
			canvasImage = canvas.NewImageFromFile(fmt.Sprintf("%d-trader.jpg", i))
		} else {
			err := app.downloadFile(trader.UserPhotoUrl, fmt.Sprintf("%d-trader", i))
			if err != nil {
				// return bundled error image
				canvasImage = canvas.NewImageFromResource(resources.ResourceNoImageAvailablePng)
			} else {
				canvasImage = canvas.NewImageFromFile(fmt.Sprintf("%d-trader.jpg", i))
			}
		}
		canvasImage.SetMinSize(fyne.NewSize(25, 25))
		canvasImage.FillMode = canvas.ImageFillContain
		//card.SetImage(canvasImage)
		cards = append(cards, card)
	}

	return cards
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
