package copy_io

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type HomeTab struct {
	*container.TabItem
	CheckoutsChartContainer *fyne.Container
}

func (app *Config) homeTab(greetingMsg string) *fyne.Container {

	greetingMsgTitle := canvas.NewText(greetingMsg, nil)

	// define the homeTabContainer
	homeTabContainer := container.NewWithoutLayout(greetingMsgTitle)

	// resize and move the homeTab elements
	greetingMsgTitle.Move(fyne.NewPos(10, 10))
	greetingMsgTitle.Resize(fyne.NewSize(100, 50))

	return homeTabContainer
}
