package nyx

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/user"
)

type CopiedTradersTab struct {
	*container.TabItem
	CopiedPositionsSlice [][]any
}

func (app *Config) copiedTradersTab() *fyne.Container {
	app.getTraders()
	copiedTradersList := app.getCopiedTraders()

	max := container.NewMax(copiedTradersList)

	return max
}

func (app *Config) getCopiedTraders() *container.Split {
	var data []string
	var slice [][]any

	for _, t := range app.User.CopiedTradersManager.Traders {
		data = append(data, t.NickName)
	}

	label := widget.NewLabel("Loading Your Positions...")
	slice = append(slice, []any{"Symbol", "Leverage", "Mark Price", "Unrealised PNL", "Last Update"})
	app.CopiedPositionsSlice = slice
	positionsTable := widget.NewTable(
		func() (int, int) {
			return len(app.CopiedPositionsSlice), len(app.CopiedPositionsSlice[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			canvasText := canvas.NewText("", nil)
			canvasText.Hide()
			return container.NewMax(lbl, canvasText)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			canvasText := container.Objects[1].(*canvas.Text)

			if i.Col == 3 && i.Row != 0 {
				lbl.Hide()
				canvasText.Hidden = false
				pnl := app.CopiedPositionsSlice[i.Row][i.Col].(float64)
				if pnl > 0 {
					canvasText.Color = color.RGBA{R: 14, G: 203, B: 129, A: 255}
				} else {
					canvasText.Color = color.RGBA{R: 246, G: 70, B: 93, A: 255}
				}
				canvasText.Text = fmt.Sprintf("%.2f", pnl)
			} else {
				canvasText.Hide()
				lbl.Hidden = false
				lbl.SetText(
					app.CopiedPositionsSlice[i.Row][i.Col].(string))
			}
		})
	colWidths := []float32{200, 100, 100, 200, 100}
	for i, w := range colWidths {
		positionsTable.SetColumnWidth(i, w)
	}
	go func() {
		app.CopiedPositionsSlice = app.getPositionsSlice()
		label.SetText("Your Positions")
		positionsTable.Refresh()
	}()

	border := container.NewBorder(
		label,
		nil,
		nil,
		nil,
		container.NewVScroll(positionsTable),
	)

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
	}
	list.OnUnselected = func(id widget.ListItemID) {
	}
	for i := range data {
		list.SetItemHeight(i, 50)
	}

	return container.NewHSplit(list, border)
}

func (app *Config) getTraders() {
	traders, err := app.DB.AllTraders()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.Traders = traders
}

func (app *Config) getPositionsSlice() [][]any {
	var slice [][]any
	var allPositions []user.Position

	slice = append(slice, []any{"Symbol", "Leverage", "Mark Price", "Unrealised PNL", "Last Update"})
	for _, p := range app.User.Profiles {
		//TODO: better check if bybit or binance api key, add binance api key code
		if p.BybitApiKey != "" {
			bybitPositions := app.getBybitPositions(p.BybitApiKey, p.BybitApiSecret, p.TestMode)
			if bybitPositions != nil {
				allPositions = append(allPositions, bybitPositions...)
			}
		}
	}

	for _, x := range allPositions {
		var currentRow []any

		currentRow = append(currentRow, fmt.Sprintf("%s Perpetual", x.Symbol))

		currentRow = append(currentRow, fmt.Sprintf("%d", x.Leverage))

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.MarkPrice))

		currentRow = append(currentRow, x.Pnl)

		currentRow = append(currentRow, fmt.Sprintf("%d", x.UpdateTimestamp))

		slice = append(slice, currentRow)
	}

	return slice
}
