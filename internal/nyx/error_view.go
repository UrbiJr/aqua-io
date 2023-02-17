package nyx

import (
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
func (nyx *Config) ShowError(currentPage string, err error) {
}

// ShowErrorAndExit displays the custom error and closes the app
func (nyx *Config) ShowErrorAndExit(err error) {
}

// NewErrorView returns a view for the custom error messages
func (nyx *Config) NewErrorView() *ErrorView {

	return nil
}
