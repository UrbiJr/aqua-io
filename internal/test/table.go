package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var data = [][]string{[]string{"label", "actions"},
	[]string{"hello", "toolbar goes here"}}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Table Widget")

	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			toolbar := widget.NewToolbar()
			toolbar.Hide()
			return container.NewMax(lbl, toolbar)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			toolbar := container.Objects[1].(*widget.Toolbar)

			if i.Row != 0 && i.Col == 1 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {}))
				}

			} else {
				toolbar.Hide()
				lbl.Hidden = false
				lbl.SetText(data[i.Row][i.Col])
			}
		})

	// hide first column (ID)
	colWidths := []float32{100, 100}
	for i, w := range colWidths {
		list.SetColumnWidth(i, w)
	}

	myWindow.SetContent(list)
	myWindow.ShowAndRun()
}
