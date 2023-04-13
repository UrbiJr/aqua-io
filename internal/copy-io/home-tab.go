package copy_io

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/UrbiJr/copy-io/internal/resources"
	"github.com/UrbiJr/copy-io/internal/utils"
	chart "github.com/wcharczuk/go-chart/v2"
)

type HomeTab struct {
	*container.TabItem
	CheckoutsChartContainer *fyne.Container
}

func (app *Config) homeTab(greetingMsg string) *fyne.Container {
	// fill default checkout data
	data := []float64{15, 5, 0, 6, 78, 2, 0, 10, 0, 20, 5, 14}
	// get the checkouts chart
	checkoutsChart := app.getCheckoutsChart(data)
	// wrap checkoutschart in a VBox
	checkoutsChartContainer := container.NewVBox(checkoutsChart)
	// assign it to our config to refresh it easily
	app.CheckoutsChartContainer = checkoutsChartContainer

	// get the checkout feed
	checkoutFeed := app.getCheckoutFeed()
	checkoutFeedTitle := canvas.NewText("Checkout Feed", nil)
	greetingMsgTitle := canvas.NewText(greetingMsg, nil)

	// get the releases feed
	releasesFeed := app.getReleasesFeed()

	// define the homeTabContainer
	homeTabContainer := container.NewWithoutLayout(greetingMsgTitle, checkoutFeedTitle, checkoutFeed, checkoutsChartContainer, releasesFeed)

	// resize and move the homeTab elements
	greetingMsgTitle.Move(fyne.NewPos(10, 10))
	greetingMsgTitle.Resize(fyne.NewSize(100, 50))

	checkoutFeedTitle.Move(fyne.NewPos(10, 60))
	checkoutFeedTitle.Resize(fyne.NewSize(100, 50))

	checkoutFeed.Move(fyne.NewPos(10, 120))
	checkoutFeed.Resize(fyne.NewSize(360, 500))

	checkoutsChartContainer.Move(fyne.NewPos(380, 10))
	checkoutsChartContainer.Resize(fyne.NewSize(850, 313))

	releasesFeed.Move(fyne.NewPos(400, 333))
	releasesFeed.Resize(fyne.NewSize(800, 310))

	return homeTabContainer
}

func (app *Config) makeReleaseList() []*widget.Card {
	var cards []*widget.Card
	var canvasImage *canvas.Image

	releases := app.fetchReleases()
	for i, release := range releases {
		card := widget.NewCard(
			release.Title,
			release.Date,
			widget.NewHyperlink("Check StockX", release.StockXLink),
		)
		if utils.DoesFileExist(fmt.Sprintf("downloads/%d-release.jpg", i)) {
			canvasImage = canvas.NewImageFromFile(fmt.Sprintf("downloads/%d-release.jpg", i))
		} else {
			err := app.downloadFile(release.ImageURL.String(), fmt.Sprintf("%d-release", i))
			if err != nil {
				// return bundled error image
				canvasImage = canvas.NewImageFromResource(resources.ResourceNoImageAvailablePng)
			} else {
				canvasImage = canvas.NewImageFromFile(fmt.Sprintf("downloads/%d-release.jpg", i))
			}
		}
		canvasImage.SetMinSize(fyne.NewSize(100, 100))
		canvasImage.FillMode = canvas.ImageFillContain
		card.SetImage(canvasImage)
		cards = append(cards, card)
	}

	return cards
}

func (app *Config) getReleasesFeed() *fyne.Container {

	hList := app.makeReleaseList()
	cardsWrapper := container.NewHBox()
	for _, x := range hList {
		cardsWrapper.Add(x)
	}
	hScroll := container.NewHScroll(cardsWrapper)

	releasesContainer := container.NewVBox(canvas.NewText("Releases", nil), hScroll)

	return releasesContainer
}

func (app *Config) getCheckoutFeed() *widget.List {
	// Generate random elements
	var listItems []string
	for i := 0; i < 30; i++ {
		listItems = append(listItems, fmt.Sprintf("Checkout %d", i+1))
	}

	// Create the list widget
	list := widget.NewList(
		func() int {
			return len(listItems)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if label, ok := item.(*widget.Label); ok {
				label.SetText(listItems[id])
			}
		},
	)

	return list
}

func (app *Config) getCheckoutsChart(data []float64) *canvas.Image {

	// the x ticks are generated automatically,
	// the x and y ranges are established automatically,
	// the canvas "box" is adjusted to fit the space the axes occupy so as not to clip.

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Ticks: []chart.Tick{
				{Value: 1.0, Label: "Jan"},
				{Value: 2.0, Label: "Feb"},
				{Value: 3.0, Label: "Mar"},
				{Value: 4.0, Label: "Apr"},
				{Value: 5.0, Label: "May"},
				{Value: 6.0, Label: "Jun"},
				{Value: 7.0, Label: "Jul"},
				{Value: 8.0, Label: "Aug"},
				{Value: 9.0, Label: "Sep"},
				{Value: 10.0, Label: "Oct"},
				{Value: 11.0, Label: "Nov"},
				{Value: 12.0, Label: "Dec"},
			},
		},
		YAxis: chart.YAxis{
			Name: "Checkouts",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0},
				YValues: data,
			},
		},
	}

	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	img, err := collector.Image()
	if err != nil {
		app.Logger.Error(err)
		log.Panic()
	}

	canvasImage := canvas.NewImageFromImage(img)
	canvasImage.SetMinSize(fyne.NewSize(850, 313))
	canvasImage.FillMode = canvas.ImageFillOriginal

	return canvasImage
}
