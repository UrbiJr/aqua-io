package cactus

import (
	"errors"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI contains information about the user interface
type UI struct {
	tui   *tview.Application
	pages *tview.Pages
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

// indexToRune is a utility function to cast an int index to rune type in order to use it as menu shortcut
func indexToRune(index int) (rune, error) {
	s2 := strconv.Itoa(index)
	var b rune
	if len(s2) > 1 {
		return 0, errors.New("index cannot exceed 1 digit")
	} else {
		b = rune(s2[0])
	}
	return b, nil
}

// NewSitelist returns a view for the sitelist
func (cactus *Cactus) NewSitelist() *tview.Flex {

	list := tview.NewList()
	i := 1
	for _, entry := range cactus.SiteList {
		shortcut, _ := indexToRune(i)
		list.AddItem(entry.Name.String(), "", shortcut, nil)
		i += 1
	}

	shortcut, _ := indexToRune(i)
	list.AddItem("Go Back", "return to the previous page", shortcut, cactus.UI.OnGoBackSelected)

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Sitelist"), 0, 1, false).
		AddItem(list, 0, 4, true).SetBorder(true)

	return flex
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

// OnGoBackSelected should be called when a user choose to go to the previous view
func (ui *UI) OnGoBackSelected() {

	// switch selection based on current page
	//utils.Info().Println(ui.pages.GetTitle())
	switch currentPageTitle := ui.pages.GetTitle(); currentPageTitle {
	case "Sitelist":
		// if current page is Sitelist, go back to main menu
		ui.pages.SwitchToPage("Main Menu")
	}

}

// OnSitelistSelected should be called when a user choose Sitelist entry on main menu
func (ui *UI) OnSitelistSelected() {

	// switch current view to Sitelist
	ui.pages.SwitchToPage("Sitelist")
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
