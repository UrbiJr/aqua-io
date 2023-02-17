package nyx

import (
	"github.com/rivo/tview"
)

// TaskCreationView contains information about the "Task Creation" view
type TaskCreationView struct {
	Title string
	View  *tview.Flex
}

// NewTaskCreationView returns a view for the sitelist
func (nyx *Config) NewTaskCreationView() *TaskCreationView {

	return nil
}
