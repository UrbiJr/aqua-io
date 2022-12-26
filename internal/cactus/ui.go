package cactus

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI contains information about the user interface
type UI struct {
	tui *tview.Application
}

// MenuEntry contains information for each entry of a menu.
// name is the entry name to display
// description is a secondary text
// selected is the callback function which gets called when the entry is selected
type MenuEntry struct {
	name        string
	label       rune
	description string
	selected    func()
}

// Menu contains information about a Cactus-AIO ui menu
type Menu struct {
	name    string
	entries []MenuEntry
}

// NewMainMenu initializes the menu for the main view
func (ui *UI) NewMainMenu(welcomeMessage string, entries []MenuEntry) *tview.Flex {

	list := tview.NewList()
	for _, entry := range entries {
		list.AddItem(entry.name, entry.description, entry.label, entry.selected)
	}

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText(welcomeMessage), 0, 1, false).
		AddItem(list, 0, 4, true).SetBorder(true)

	return flex
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
