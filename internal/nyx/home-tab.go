package nyx

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	chart "github.com/wcharczuk/go-chart/v2"
)

func (nyx *Config) homeTab() *fyne.Container {
	// fill default checkout data
	data := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// get the checkouts chart
	checkoutsChart := nyx.getCheckoutsChart(data)
	// wrap checkoutschart in a VBox
	checkoutsChartContainer := container.NewVBox(checkoutsChart)
	// assign it to our config to refresh it easily
	nyx.CheckoutsChartContainer = checkoutsChartContainer

	// get the checkout feed
	checkoutFeed := nyx.getCheckoutFeed()

	// get the releases feed

	// define the homeTabContainer
	homeTabContainer := container.NewWithoutLayout(checkoutsChartContainer, checkoutFeed)

	// resize and move the homeTab elements
	checkoutsChart.Move(fyne.NewPos(10, 10))
	checkoutsChart.Resize(fyne.NewSize(800, 313))
	checkoutFeed.Move(fyne.NewPos(810, 10))
	checkoutFeed.Resize(fyne.NewSize(360, 520))

	return homeTabContainer
}

func (nyx *Config) getCheckoutFeed() *widget.List {
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

func (nyx *Config) getCheckoutsChart(data []float64) *canvas.Image {

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
		nyx.Logger.Error(err)
		log.Panic()
	}

	canvasImage := canvas.NewImageFromImage(img)
	canvasImage.SetMinSize(fyne.NewSize(800, 313))

	return canvasImage
}
