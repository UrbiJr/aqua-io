package cactus

import (
	"github.com/rivo/tview"
)

// ConfirmButtonFunc gets called if user presses "Yes" or "No" button
type ConfirmButtonFunc func()

// ConfirmView contains information about the view to confirm a user choice
type ConfirmView struct {
	Title    string
	PrevPage string // title of the page which showed the confirm view
	Modal    *tview.Modal
	View     *tview.Flex
}

// newConfirmModal returns an empty modal
func newConfirmModal() *tview.Modal {

	modal := tview.NewModal().
		AddButtons([]string{"Yes"}).
		AddButtons([]string{"No"})

	return modal
}

/*
ShowConfirm displays the confirm message with Yes/No buttons
yesFunc gets called if user presses Yes button
noFunc gets called if user presses No button
*/
func (cactus *Cactus) ShowConfirm(yesFunc ConfirmButtonFunc, noFunc ConfirmButtonFunc, confirmMessage string) {
	cactus.ConfirmView.Modal.
		// set function when OK button (any button) gets pressed
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// hide this page and go back to the previous page
			if buttonLabel == "Yes" {
				yesFunc()
			} else if buttonLabel == "No" {
				noFunc()
			}

		})
	cactus.ConfirmView.Modal.SetText(confirmMessage)
	cactus.pages.SwitchToPage(cactus.ConfirmView.Title)
}

// NewConfirmView returns a view to confirm a user choice
func (cactus *Cactus) NewConfirmView() *ConfirmView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	modal := newConfirmModal()

	flex.SetDirection(tview.FlexRow).
		AddItem(modal, 0, 4, true)

	return &ConfirmView{Title: "Confirm", PrevPage: "Main Menu", Modal: modal, View: flex}
}
