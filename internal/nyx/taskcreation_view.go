package nyx

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TaskCreationView contains information about the "Task Creation" view
type TaskCreationView struct {
	Title string
	View  *tview.Flex
}

// NewTaskCreationView returns a view for the sitelist
func (nyx *Nyx) NewTaskCreationView() *TaskCreationView {

	list := tview.NewList()
	i := 1
	for _, entry := range nyx.SiteList {
		shortcut, _ := indexToRune(i)
		list.AddItem(entry.Name.String(), "", shortcut, nil)
		i += 1
	}

	list.AddItem("Go Back", "return to the previous page", 'b', nyx.UI.OnGoBackSelected)

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("🤖 Task Creation"), 0, 1, false).
		AddItem(list, 0, 4, true).
		SetBorder(true)

	return &TaskCreationView{Title: "Task Creation", View: flex}
}
