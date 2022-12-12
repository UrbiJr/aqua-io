package cactus

import (
	"github.com/rivo/tview"
)

// UI contains information about the user interface
type UI struct {
	tui *tview.Application
}

// NewUI initializes a tview and assigns it to the returned UI object
func NewUI() *UI {
	ui := new(UI)
	ui.tui = tview.NewApplication()

	return ui
}

func (ui *UI) SetListeners() {
	/*
		ui.tui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 113 { // user presses 'q' key
				ui.tui.Stop()
			}
			return event
		})
	*/

}
