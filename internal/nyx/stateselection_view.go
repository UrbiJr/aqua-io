package nyx

import (
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
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
func (nyx *Config) RefreshStateSelectionView(currentPage string, countryCode string, profile *user.Profile) {
	nyx.StateSelectionView.Form.Clear(true)
	var stateCode string

	states, err := utils.GetStates(countryCode)
	if err != nil {
		nyx.pages.SwitchToPage(currentPage)
	}

	nyx.StateSelectionView.Form.AddDropDown("States", states, 0, func(state string, index int) {
		stateCode, err = utils.GetStateCode(countryCode, state)
	})

	nyx.StateSelectionView.Form.AddButton("Save", func() {
		profile.State = stateCode
		nyx.pages.SwitchToPage(currentPage)
	})

	nyx.StateSelectionView.Form.AddButton("Cancel", func() {
		nyx.pages.SwitchToPage(currentPage)
	})
}

// NewStateSelectionView returns a view for the state selection by country
func (nyx *Config) NewStateSelectionView() *StateSelectionView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	form := tview.NewForm()

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("State selection"), 0, 1, false).
		AddItem(form, 0, 4, true).
		SetBorder(true)

	return &StateSelectionView{Title: "State selection", Form: form, View: flex}
}
