package cactus

import (
	"errors"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

// MainMenuView contains information about the "Main Menu" view
type MainMenuView struct {
	Title string
	View  *tview.Flex
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

// NewMainMenuView initializes the menu for the main view
func (ui *UI) NewMainMenuView(entries []MenuEntry) *MainMenuView {

	list := tview.NewList()
	for _, entry := range entries {
		list.AddItem(entry.name, entry.description, entry.label, entry.selected)
	}

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Main Menu"), 0, 1, false).
		AddItem(list, 0, 4, true).
		SetBorder(true)

	return &MainMenuView{Title: "Main Menu", View: flex}
}
