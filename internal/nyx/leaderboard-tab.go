package nyx

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/resources"
	"github.com/UrbiJr/nyx/internal/utils"
)

type LeaderboardTab struct {
}

func (app *Config) leaderboardTab() *fyne.Container {
	grid := container.NewAdaptiveGrid(3)
	cards := app.makeTradersCards()
	for _, card := range cards {
		grid.Add(card)
	}
	vScroll := container.NewVScroll(grid)
	releasesContainer := container.NewWithoutLayout(vScroll)

	vScroll.Move(fyne.NewPos(10, 10))
	vScroll.Resize(fyne.NewSize(1260, 700))

	return releasesContainer
}

func (app *Config) makeTradersCards() []*widget.Card {
	var cards []*widget.Card
	var canvasImage *canvas.Image

	traders, err := app.fetchTraders()
	if err == nil {
		for i, trader := range traders {
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
	}

	return cards
}
