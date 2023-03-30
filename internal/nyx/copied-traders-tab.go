package nyx

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CopiedTradersTab struct {
	*container.TabItem
}

func (app *Config) copiedTradersTab() *fyne.Container {
	copiedTradersList := app.getCopiedTraders()

	max := container.NewMax(copiedTradersList)

	return max
}

func (app *Config) getCopiedTraders() *container.Split {
	data := make([]string, 1000)
	for i := range data {
		data[i] = "Trader " + strconv.Itoa(i)
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select A Copied Trader From The List")
	hbox := container.NewHBox(icon, label)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Template Object"), widget.NewToolbar())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			toolbar := item.(*fyne.Container).Objects[2].(*widget.Toolbar)
			if len(toolbar.Items) == 0 {
				toolbar.Append(widget.NewToolbarAction(theme.VisibilityIcon(), func() {
					//app.traderDialog()
				}))
				toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Stop copying?", "Confirming will delete all current positions copied from this trader.", func(deleted bool) {

					}, app.MainWindow)
				}))
			}
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id] + "\ncopied positions: 0")
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(fmt.Sprintf("Copied positions from %s go here", data[id]))
		icon.SetResource(theme.AccountIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select A Copied Trader From The List")
		icon.SetResource(nil)
	}
	for i := range data {
		list.SetItemHeight(i, 50)
	}

	return container.NewHSplit(list, container.NewCenter(hbox))
}

func (app *Config) getPositions() {
	positions, err := app.DB.AllPositions()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.Positions = positions
}

func (app *Config) getTraders() {
	traders, err := app.DB.AllTraders()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.Traders = traders
}
