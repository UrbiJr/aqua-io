package cactus

import (
	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StateSelectionView contains information about the "State selection" view
type StateSelectionView struct {
	Title string
	Form  *tview.Form
	View  *tview.Flex
}

// RefreshStateSelectionView fills the state selection form with the states of the specified country
func (cactus *Cactus) RefreshStateSelectionView(countryCode string, profile *user.Profile) {
	cactus.StateSelectionView.Form.Clear(true)

	states, err := utils.GetStates(countryCode)
	if err != nil {
		cactus.pages.SwitchToPage(cactus.NewProfileView.Title)
	}

	cactus.StateSelectionView.Form.AddDropDown("States", states, 0, func(state string, index int) {
		stateCode, err := utils.GetStateCode(countryCode, state)
		if err == nil {
			profile.State = stateCode
		}
	})

	cactus.StateSelectionView.Form.AddButton("Save", func() {
		cactus.pages.SwitchToPage(cactus.NewProfileView.Title)
	})
}

// NewStateSelectionView returns a view for the state selection by country
func (cactus *Cactus) NewStateSelectionView() *StateSelectionView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("State selection"), 0, 1, false).
		AddItem(form, 0, 4, true).
		SetBorder(true)

	return &StateSelectionView{Title: "State selection", Form: form, View: flex}
}
