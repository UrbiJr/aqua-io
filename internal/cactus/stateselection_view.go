package cactus

import (
	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewStateSelectionView returns a view for the state selection by country
func (cactus *Cactus) NewStateSelectionView(countryCode string, profile *user.Profile) *tview.Flex {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	states, err := utils.GetStates(countryCode)
	if err != nil {
		cactus.pages.SwitchToPage("New Profile")
	}

	form.AddDropDown("States", states, 0, func(state string, index int) {
		stateCode, err := utils.GetStateCode(countryCode, state)
		if err == nil {
			profile.State = stateCode
		}
	})

	form.AddButton("Save", func() {
		cactus.pages.SwitchToPage("New Profile")
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("State selection"), 0, 1, false).
		AddItem(form, 0, 4, true).
		SetBorder(true)

	return flex
}
