package copy_io

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type HomeTab struct {
	*container.TabItem
	CheckoutsChartContainer *fyne.Container
}

func (app *Config) homeTab(greetingMsg string) *fyne.Container {

	greetingMsgTitle := widget.NewLabelWithStyle(greetingMsg, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	// define the homeTabContainer
	homeTabContainer := container.NewVBox(greetingMsgTitle)

	return homeTabContainer
}
