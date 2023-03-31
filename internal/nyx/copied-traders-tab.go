package nyx

import (
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
	app.getTraders()
	copiedTradersList := app.getCopiedTraders()

	max := container.NewMax(copiedTradersList)

	return max
}

func (app *Config) getCopiedTraders() *container.Split {
	var data []string
	for _, t := range app.User.CopiedTradersManager.Traders {
		data = append(data, t.NickName)
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
			t := app.User.CopiedTradersManager.Traders[id]
			toolbar := item.(*fyne.Container).Objects[2].(*widget.Toolbar)
			if len(toolbar.Items) == 0 {
				toolbar.Append(widget.NewToolbarAction(theme.VisibilityIcon(), func() {
					app.traderDialog(t)
				}))
				toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Stop copying?", "Confirming will delete all current positions copied from this trader.", func(deleted bool) {

					}, app.MainWindow)
				}))
			}
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id] + "\n" + t.EncryptedUid)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText("Loading Positions...")
		t := app.User.CopiedTradersManager.Traders[id]
		icon.SetResource(theme.AccountIcon())
		go func() {
			positions, err := app.fetchTraderPositions(t.EncryptedUid)
			if err != nil {
				app.Logger.Error(err)
				return
			}
			positionsText := ""
			for _, p := range positions {
				positionsText = positionsText + p.Symbol + "\n"
			}
			label.SetText(positionsText)
		}()
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

func (app *Config) getTraders() {
	traders, err := app.DB.AllTraders()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.Traders = traders
}
