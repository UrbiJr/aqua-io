package nyx

import (
	"fmt"

	"github.com/rivo/tview"
)

// ErrorView contains information about the view for the custom error messages
type ErrorView struct {
	Title    string
	PrevPage string // title of the page which showed the error view
	Modal    *tview.Modal
	View     *tview.Flex
}

// newModal returns an empty modal
func newErrorModal() *tview.Modal {

	modal := tview.NewModal().
		AddButtons([]string{"Ok"})

	return modal
}

// ShowError displays the custom error
func (nyx *Nyx) ShowError(currentPage string, err error) {
	nyx.ErrorView.Modal.
		// set function when OK button (any button) gets pressed
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// hide this page and go back to the previous page
			nyx.pages.SwitchToPage(currentPage)
		})
	nyx.ErrorView.Modal.SetText(fmt.Sprintf("Error: %s", err.Error()))
	nyx.pages.SwitchToPage(nyx.ErrorView.Title)
}

// ShowErrorAndExit displays the custom error and closes the app
func (nyx *Nyx) ShowErrorAndExit(err error) {
	nyx.ErrorView.Modal.
		// set function when OK button (any button) gets pressed
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// hide this page and go back to the previous page
			nyx.Quit()
		})
	nyx.ErrorView.Modal.SetText(fmt.Sprintf("Error: %s", err.Error()))
	nyx.pages.SwitchToPage(nyx.ErrorView.Title)
}

// NewErrorView returns a view for the custom error messages
func (nyx *Nyx) NewErrorView() *ErrorView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	modal := newErrorModal()

	flex.SetDirection(tview.FlexRow).
		AddItem(modal, 0, 4, true)

	return &ErrorView{Title: "Error", PrevPage: "Main Menu", Modal: modal, View: flex}
}
