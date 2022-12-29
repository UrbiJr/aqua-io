package cactus

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewSitelistView returns a view for the sitelist
func (cactus *Cactus) NewSitelistView() *tview.Flex {

	list := tview.NewList()
	i := 1
	for _, entry := range cactus.SiteList {
		shortcut, _ := indexToRune(i)
		list.AddItem(entry.Name.String(), "", shortcut, nil)
		i += 1
	}

	list.AddItem("Go Back", "return to the previous page", 'b', cactus.UI.OnGoBackSelected)

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Sitelist"), 0, 1, false).
		AddItem(list, 0, 4, true).
		SetBorder(true)

	return flex
}
