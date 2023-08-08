package aqua_io

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HomeTab struct {
	*container.TabItem
	CheckoutsChartContainer *fyne.Container
}

func (app *Config) homeTab(greetingMsg string) *fyne.Container {

	closeGreetingMsgButton := widget.NewButtonWithIcon("", theme.CancelIcon(), nil)
	greetingMsgTitle := widget.NewLabelWithStyle(greetingMsg, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	greetingMsgCard := widget.NewCard("", "", container.NewWithoutLayout(greetingMsgTitle, closeGreetingMsgButton))

	greetingMsgTitle.Move(fyne.NewPos(0, 0))
	closeGreetingMsgButton.Move(fyne.Position{X: greetingMsgTitle.MinSize().Width + 15, Y: 0})
	closeGreetingMsgButton.Resize(fyne.NewSize(15, 15))

	closeGreetingMsgButton.OnTapped = func() {
		greetingMsgCard.Hide()
	}

	// define the homeTabContainer
	homeTabContainer := container.NewWithoutLayout(greetingMsgCard)
	greetingMsgCard.Move(fyne.NewPos(10, 10))
	greetingMsgCard.Resize(fyne.Size{Width: greetingMsgTitle.MinSize().Width + 60, Height: greetingMsgTitle.MinSize().Height + 10})

	return homeTabContainer
}
